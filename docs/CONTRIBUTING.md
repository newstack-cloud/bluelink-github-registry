# Contributing to the Celerity GitHub Registry

## Setup

Ensure that git uses the custom directory for git hooks

```bash
git config core.hooksPath .githooks
```

### Prerequisites

- Docker >= 20.10.6 (Comes with Docker Desktop)
- [Go](https://golang.org/dl/) >=1.23
- [Node.js](https://nodejs.org/en/download/) >= 20.14.0 (For tooling)
- [Yarn](https://yarnpkg.com/getting-started/install) >= 1.22.19 (For tooling)

### NPM dependencies for tooling

There are npm dependencies that provide tools that are used in git hooks to ensure commits follow the [commit guidelines](./COMMIT_GUIDELINES.md).

Yarn should be used as the package manager for these dependencies.

Install dependencies from the root directory by simply running:
```bash
yarn
```

### Setting hosts for the local registry

To run the registry locally, you need to set up a local DNS entry for the registry.
This can be done by adding the following line to your `/etc/hosts` (or equivalent) file:

```bash
127.0.0.1 gh-registry.celerity.local
```

This will allow you to access the registry at `http://gh-registry.celerity.local`
when you spin up the registry locally using the provided docker compose stack.

**If you are not using the default port (8085), you will need to update `nginx.local.conf` to use your alternate port.**

### Go dependencies

Dependencies are managed with Go modules (go.mod) and will be installed automatically when you first
run tests.

If you want to install dependencies manually you can run:

```bash
go mod download
```

### Preparing Environment Variables

You will need to prepare environment variables for the local development environment,
both for tests and for running the service locally.

1) Copy the `.env.example` file to `.env`:

```bash
cp .env.example .env
```

2) Prepare an ASCII-armored public key file (extracted from a GPG key) named `public.key` in the root directory.

3) Run the tool to inject the public key into the signing keys environment variable:

```bash
go run ./tools/signing-keys/main.go -insert=.env public.key
```

4) Repeat steps 1-3 for the `.env.test.example` file to create and populate the `.env.test` file.

## Running tests

```bash
bash ./scripts/run-tests.sh
```

## Local Development

### Running the service locally

To run the service locally, you can use the following command:

```bash
docker compose up
# Or to run in the background:
docker compose up -d
```

## Releasing

To release a new version of the deploy engine, you need to create a new tag and push it to the repository.
The release workflow will automatically build and publish the new version to the GitHub Container Registry.

The format must be `vX.Y.Z` where `X.Y.Z` is the semantic version number.

See [here](https://go.dev/wiki/Modules#publishing-a-release).

1. add a change log entry to the `CHANGELOG.md` file following the template below:

```markdown
## [0.2.0] - 2024-06-05

### Fixed:

- Corrects bug in extracting assets from plugin repository release.

### Added

- Adds graceful error handling for fetching plugin repository releases.
```

2. Create and push the new tag:

```bash
git tag -a v0.2.0 -m "chore: Release v0.2.0"
git push --tags
```

Be sure to add a release for the tag with notes following this template:

Title: `v0.2.0`

```markdown
## Fixed:

- Corrects bug in extracting assets from plugin repository release.

## Added

- Adds graceful error handling for fetching plugin repository releases.
```
