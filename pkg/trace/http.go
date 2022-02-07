package trace

import (
	"net/http"
	"sync"

	"github.com/felixge/httpsnoop"
	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib"
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const tracerName = "github.com/ca-risken/common/pkg/trace"

type recordingResponseWriter struct {
	writer  http.ResponseWriter
	written bool
	status  int
}

var rrwPool = &sync.Pool{
	New: func() interface{} {
		return &recordingResponseWriter{}
	},
}

func getRRW(writer http.ResponseWriter) *recordingResponseWriter {
	rrw := rrwPool.Get().(*recordingResponseWriter)
	rrw.written = false
	rrw.status = http.StatusOK
	rrw.writer = httpsnoop.Wrap(writer, httpsnoop.Hooks{
		Write: func(next httpsnoop.WriteFunc) httpsnoop.WriteFunc {
			return func(b []byte) (int, error) {
				if !rrw.written {
					rrw.written = true
				}
				return next(b)
			}
		},
		WriteHeader: func(next httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
			return func(statusCode int) {
				if !rrw.written {
					rrw.written = true
					rrw.status = statusCode
				}
				next(statusCode)
			}
		},
	})
	return rrw
}

func putRRW(rrw *recordingResponseWriter) {
	rrw.writer = nil
	rrwPool.Put(rrw)
}

func OtelChiMiddleware(serverName, skipPath string) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// Skip tracing when the path is specified, ex. /healthz
			if skipPath == r.URL.Path {
				next.ServeHTTP(w, r)
				return
			}

			opts := []oteltrace.SpanStartOption{
				oteltrace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", r)...),
				oteltrace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(r)...),
				oteltrace.WithSpanKind(oteltrace.SpanKindServer),
			}
			tracer := otel.GetTracerProvider().Tracer(tracerName,
				oteltrace.WithInstrumentationVersion(contrib.SemVersion()))
			// We want to set the route pattern to span's name, but it is returned only after ServeHTTP.
			// So it is set after the execution.
			ctx, span := tracer.Start(r.Context(), "", opts...)
			defer span.End()
			rrw := getRRW(w)
			defer putRRW(rrw)
			r = r.WithContext(ctx)
			next.ServeHTTP(rrw.writer, r)

			routePattern := chi.RouteContext(r.Context()).RoutePattern()
			span.SetName(routePattern)
			span.SetAttributes(semconv.HTTPServerAttributesFromHTTPRequest(serverName, routePattern, r)...)
			spanStatus, spanMessage := semconv.SpanStatusFromHTTPStatusCode(rrw.status)
			span.SetStatus(spanStatus, spanMessage)
		}
		return http.HandlerFunc(fn)
	}
}
