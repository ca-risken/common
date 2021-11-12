PKG_TARGETS = database logging portscan sqs xray rpc
LINT_TARGETS = $(PKG_TARGETS:=.lint)

.PHONY: lint
lint: $(LINT_TARGETS)
%.lint: FAKE
	sh hack/golinter.sh pkg/$(*)

FAKE:
