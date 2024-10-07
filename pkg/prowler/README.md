# Prowler

## Customize Prowler

You can customize several settings for Prowler by modifying the `prowler.yaml` file.

```yaml
# ignorePlugin
# Specify plugins to be ignored here.
ignorePlugin:
  - entra_global_admin_in_less_than_five_users

# specificPluginSetting
# You can set scores, tags, recommendations, etc. for each plugin.
specificPluginSetting:
  category/pluginName:
    # score (0.1-1.0):
    # Set the score for the plugin
    # If no score is set, the score will be according to Severity
    score: 0.8

    # skipResourceNamePattern:
    # Specify resource name patterns to ignore resources that match these patterns.
    skipResourceNamePattern:
      - "your_resource_uid"

    # ignoreMessagePattern:
    # Specify message patterns to ignore messages that match these patterns.
    ignoreMessagePattern: "Non-privileged user [a-zA-Z0-9-_]+ does not have MFA."

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

This configuration allows you to customize Prowler's behavior, including setting ignoring specific plugins, and configuring plugin-specific settings such as scores, resource name patterns to skip, tags, and recommendations.

## Generate Prowler YAML file

You can generate the Prowler YAML file using the following command.

```bash
$ make generate-yaml
```

If you want to generate a YAML file with a specific commit hash, specify COMMIT_HASH.

```bash
$ COMMIT_HASH=xxxxxxx go run generate-prowler-yaml/main.go
```


## Update your Prowler YAML file

If updating your Prowler YAML, specify the current YAML in PLUGIN_FILE

```bash
# Azure plugin
$ PLUGIN_FILE=path/to/your/prowler.yaml \
  PLUGIN_DIR=providers/azure/services \
  COMMIT_HASH=xxxxxxx \
  go run generate-prowler-yaml/main.go
```

