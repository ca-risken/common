# Sensitive content scanning

Use `Scanner.Scan(ctx, ScanRequest{FileName: name, Reader: reader})` to classify one file with the Core sensitive-content API. The caller opens and closes the reader. The shared `ScanResponse` holds file, PDF page, or ZIP entry classifications.

- The scanner converts UTF-8, UTF-16 with a BOM, and Windows-31J CSV, TSV, and TXT to UTF-8. Text samples are limited to 32 KiB.
- JPEG and PNG inputs must meet the Core limits of 2 MiB and a 2,048 px maximum edge.
- For PDF files up to 5 MiB, PDFium WASM extracts text or renders JPEG images from the first 20 pages.
- For `.zip` files, the scanner reads supported members one at a time without extracting files to disk. It reviews up to 100 archive entries and 25 MiB of expanded data per ZIP. The configured single-file limit applies to the ZIP and each member. Encrypted ZIPs and nested archives are outside the supported scope.
- For `.log.gz` files, the scanner decompresses the log and uses the text classification path. The configured single-file limit applies to both compressed and expanded data.

Use `NewFindingFile` to convert a scan to file-level Finding metadata with human-readable classification names. A PDF or ZIP produces one file result. ZIP results include detected member paths in `sensitive_entries`. Use `NewFindingScan` to generate the shared Finding summary fields.
Use the file's sensitivity decision to select detected files for Finding metadata. The serialized `test_data` field is Boolean and is omitted when unknown.
