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

## Additional documentation

- [Contributing](docs/CONTRIBUTING.md)
