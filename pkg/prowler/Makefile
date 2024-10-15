.PHONY: generate-yaml
generate-yaml:
	go run generate-prowler-yaml/main.go

.PHONY: update-yaml
update-yaml:
	PLUGIN_FILE=prowler.yaml \
	PLUGIN_DIR=prowler/providers/azure/services \
	COMMIT_HASH=ea4bf5b484b3d39ebf3f6b91e9df61df1dc4643f \
	go run generate-prowler-yaml/main.go
