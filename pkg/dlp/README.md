# Sensitive content scanning

Use `Scanner.Scan(ctx, ScanRequest{FileName: name, Reader: reader})` to classify one file with the Core sensitive-content API. The caller opens and closes the reader. The shared `ScanResponse` holds either a file classification or PDF page classifications.

- The scanner converts UTF-8, UTF-16 with a BOM, and Windows-31J CSV, TSV, and TXT to UTF-8. Text samples are limited to 32 KiB.
- JPEG and PNG inputs must meet the Core limits of 2 MiB and a 2,048 px maximum edge.
- For PDF files up to 5 MiB, PDFium WASM extracts text or renders JPEG images from the first 20 pages.

Use `NewFindingFile` to convert a scan to file-level Finding metadata with human-readable classification names. This combines the page decisions of a PDF into one file result. Use `NewFindingScan` to generate the shared Finding summary fields.
Use the file's sensitivity decision to select detected files for Finding metadata. The serialized `test_data` field is Boolean and is omitted when unknown.
