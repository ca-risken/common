# Sensitive content scanning

Use `Scanner.Scan(ctx, ScanRequest{FileName: name, Reader: reader})` to classify one file with the Core sensitive-content API. The caller opens and closes the reader. The shared `ScanResponse` holds either a file classification or PDF page classifications.

- The scanner converts UTF-8, UTF-16 with a BOM, and Windows-31J CSV, TSV, and TXT to UTF-8. Text samples are limited to 32 KiB.
- JPEG and PNG inputs must meet the Core limits of 2 MiB and a 2,048 px maximum edge.
- For PDF files up to 5 MiB, PDFium WASM extracts text or renders JPEG images from the first 20 pages.

Use `NewFindingScan` to convert classifications into the shared Finding summary fields.
