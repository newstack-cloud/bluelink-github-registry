# Configuring Bluelink CLI to use the registry

This document provides guidance on how to configure the Bluelink CLI to use an instance of the Bluelink GitHub Registry for private plugins.

## Configuration

There is only one piece of configuration required to use the Bluelink GitHub Registry with the Bluelink CLI, and that is to configure the authentication token.

The Bluelink CLI will automatically resolve the discovery document based on the plugin ID prefix (e.g. `registry.example.io` for a plugin with an ID of `registry.example.io/org/provider`), so there is no need to configure registry URLs for the CLI.

Authentication needs to be configured to let the CLI know how to authenticate with the registry when calling the registry protocol endpoints.

### Configure authentication

In your `$HOME/.bluelink/auth.json`, you will need to the following entry:

```json
{
    "{registry_domain}": {
        "apiKey": "{githubAccessToken}"
    }
}
```

Where `{registry_domain}` is the domain of your private registry (e.g. `registry.example.io`) and `{githubAccessToken}` is a fine-grained GitHub personal access token with permissions to read from the private repositories that host the plugins and access their release artifacts.
You will need to make sure that the fine-grained token is for the resource owner of the plugin repositories (e.g. the GitHub user or organisation that owns the plugin repositories).
