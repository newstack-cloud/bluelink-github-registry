# Celerity GitHub Registry

[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=two-hundred_celerity-github-registry&metric=coverage)](https://sonarcloud.io/summary/new_code?id=two-hundred_celerity-github-registry)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=two-hundred_celerity-github-registry&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=two-hundred_celerity-github-registry)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=two-hundred_celerity-github-registry&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=two-hundred_celerity-github-registry)

A [Celerity Registry Protocol](https://www.celerityframework.io/plugin-framework/docs/registry-protocols-formats/registry-protocol) implementation for plugins sourced from private GitHub repositories.

## Usage

The Celerity GitHub Registry is distributed as a Docker image that is published to the GitHub Container Registry.
To run the registry, you can use the following command:

```bash
docker run -d \
  --name celerity-github-registry \
  -p 8085:8085 \
  ghcr.io/two-hundred/celerity-github-registry:latest
```

Docker compose example:

```yaml
services:
    celerity-github-registry:
        image: ghcr.io/two-hundred/celerity-github-registry:latest
        container_name: celerity-github-registry
        ports:
        - "8085:8085"
```

You can deploy the registry to any environment that supports Docker, including Kubernetes, AWS ECS, and Azure Container Instances.

### Releases and Docker tags

The Celerity GitHub Registry is versioned using [semantic versioning](https://semver.org/).
The Docker image is tagged with the version number, and the `latest` tag always points to the latest stable release.
On every push or PR into the `main` branch, a development build is pushed to the `main` docker image tag.

## Configuration

Configuration for the registry is expected to be provided via environment variables.

The following environment variables are supported, where some are required and some are optional:

### Port 

`CELERITY_GITHUB_REGISTRY_PORT`

**_optional_**

The port on which the registry will listen for incoming requests.

**default value:** `8085`

### Access Log File

`CELERITY_GITHUB_REGISTRY_ACCESS_LOG_FILE`

**_optional_**

The file path to which the access log will be written.
If not specified, the access log will be written to `stdout`.


### Auth Token Header

`CELERITY_GITHUB_REGISTRY_AUTH_TOKEN_HEADER`

**_optional_**

The name of the header that will be used to pass the authentication token to the registry.
This will be used to authenticate requests and is where a GitHub personal access token should be provided.
The service discovery document for the registry will include this header in the `auth.v1` section as the value of the `apiKeyHeader` field.

**default value:** `celerity-gh-registry-token`

### Registry Base URL 

`CELERITY_GITHUB_REGISTRY_BASE_URL`

**_required_**

The base URL of the registry. This is used to construct the endpoints in the service discovery document
for the registry.
This should be a fully qualified URL, including the scheme (http or https).
When deploying the registry this would be something like `https://registry.example.io`,
or when running locally it could be `http://registry.example.local`.

The fully-qualified domain name in the URL must be the same as the prefix for plugins in the registry.
For example, if the base URL is `https://registry.example.io`, then the plugin prefix must be `registry.example.io` in order for the Celerity CLI to be able to find the correct service discovery document.

_Even when running locally, it is best to use a host name alias instead of localhost with port numbers due to the way the CLI resolves the registry based on the domain prefix in the plugin ID._

**When deploying the registry, it is advised that you set up a TLS certificate for the registry to ensure that the registry is only accessible over HTTPS.**

## Additional documentation

- [Contributing](docs/CONTRIBUTING.md)
