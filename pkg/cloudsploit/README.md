# CloudSploit

## Customize CloudSploit

You can customize several settings for CloudSploit by modifying the `cloudsploit.yaml` file.

```yaml
# defaultScore (1-10)
# If a plugin's score is not set, this default score will be applied.
defaultScore: 3

# ignorePlugin
# Specify plugins to be ignored here.
ignorePlugin:
  - EC2/ebsSnapshotPublic
  - Lambda/lambdaPublicAccess
  - SNS/topicPolicies
  - SQS/sqsPublicAccess

# specificPluginSetting
# You can set scores, tags, recommendations, etc. for each plugin.
specificPluginSetting:
  category/pluginName:
    # score (1-10):
    # Set the score for the plugin
    score: 8

    # skipResourceNamePattern:
    # Specify resource name patterns to ignore resources that match these patterns.
    skipResourceNamePattern:
      - "arn:aws:s3:::bucket-name"
      - "ignoreResourceName"

    # ignoreMessagePattern:
    # Specify message patterns to ignore messages that match these patterns.
    ignoreMessagePattern: "Domain: .+ expires in (?:2[5-9]|[3-9]\d|\d{3,}) days"

    # tags:
    # You can set tags for resources.
    # Tags can be used for search filters, etc.
    tags:
      - tag1
      - tag2

    # recommend:
    # You can set recommendations.
    recommend:
      risk: "..."
      remediation: "xxxxx"
```

This configuration allows you to customize CloudSploit's behavior, including setting default scores, ignoring specific plugins, and configuring plugin-specific settings such as scores, resource name patterns to skip, tags, and recommendations.

## Generate CloudSploit YAML file

You can generate the latest CloudSploit YAML file using the following command.

```bash
$ make generate-yaml
```

If you want to generate the YAML file with a specific commit hash, you can use the following command.

```bash
$ COMMIT_HASH=xxxxxxx go run generate-cloudsploit-yaml/main.go
```


## Update your CloudSploit YAML file

You can update the CloudSploit YAML file with the following command.

```bash
# AWS plugin
$ PLUGIN_FILE=path/to/your/cloudsploit.yaml \
  PLUGIN_DIR=plugins/aws \
  COMMIT_HASH=xxxxxxx \
  go run generate-cloudsploit-yaml/main.go

# GCP plugin
$ PLUGIN_FILE=path/to/your/cloudsploit.yaml \
  PLUGIN_DIR=plugins/gcp \
  COMMIT_HASH=xxxxxxx \
  go run generate-cloudsploit-yaml/main.go
```

