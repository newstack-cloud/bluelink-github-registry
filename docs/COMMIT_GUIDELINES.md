# Commit Guidelines

## 1) Conventional Commits

For all contributions to this repo, you must use the conventional commits standard defined [here](https://www.conventionalcommits.org/en/v1.0.0/).

This is used to generate automated change logs, allow for tooling to decide semantic versions for all applications,
provide a rich and meaningful commit history along with providing
a base for more advanced tooling to allow for efficient searches for decisions and context related to commits and code.

### Commit types

**The following commit types are supported in the Bluelink Registry project:**

- `fix:` - Should be used for any bug fixes.
- `build:` - Should be used for functionality related to building the application.
- `release:` - Should be used for functionality related to releasing a version of the application.
- `revert:` - Should be used for any commits that revert changes.
- `wip:` - Should be used for commits that contain work in progress.
- `feat:` - Should be used for any new features added, regardless of the size of the feature.
- `chore:` - Should be used for tasks such as releases or patching dependencies.
- `ci:` - Should be used for any work on GitHub Action workflows or scripts used in CI.
- `docs:`- Should be used for adding or modifying documentation.
- `style:` - Should be used for code formatting commits and linting fixes.
- `refactor:` - Should be used for any type of refactoring work that is not a part of a feature or bug fix.
- `perf:` - Should be used for a commit that represents performance improvements.
- `test:` - Should be used for commits that are purely for automated tests.
- `instr:` - Should be used for commits that are for instrumentation purposes. (e.g. logs, trace spans and telemetry configuration)
- `local:` - Should be used for commits that are for local environment infrastructure/tools/scripts.

### Example commit

```bash
git commit -m 'feat: add endpoint to get plugin version package information

Adds the endpoint to get version package information for a requested
operating system and architecture.
'
```

## 2) You must use the imperative mood for commit headers.

https://cbea.ms/git-commit/#imperative

The imperative mood simply means naming the subject of the commit as if it is a unit of work that can be applied instead of reporting facts about work done.

If applied, this commit will **your subject line here**.

Read the article above to find more examples and tips for using the imperative mood.
