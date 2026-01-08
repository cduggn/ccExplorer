# Release Process

This document describes the release process for ccExplorer.

## Overview

ccExplorer uses [GoReleaser](https://goreleaser.com/) to automate the release process. Releases are triggered when a version tag is pushed to the `main` branch.

## Prerequisites

- Write access to the repository
- All changes must be merged to `main` branch
- CI tests must pass on `main` branch

## Release Steps

### 1. Ensure Main Branch is Clean

```bash
git checkout main
git pull origin main
```

Verify that:
- All intended changes are merged to `main`
- CI workflows are passing
- No open critical issues

### 2. Determine Version Number

Follow [Semantic Versioning](https://semver.org/):

- **MAJOR** version (X.0.0): Incompatible API changes
- **MINOR** version (0.X.0): New functionality (backwards compatible)
- **PATCH** version (0.0.X): Bug fixes (backwards compatible)

Example: If current version is `v0.6.0`:
- Bug fix → `v0.6.1`
- New feature → `v0.7.0`
- Breaking change → `v1.0.0`

### 3. Create and Push Tag

**From Main Branch:**

```bash
# Create annotated tag
git tag -a v0.7.0 -m "Release v0.7.0"

# Push tag to trigger release
git push origin v0.7.0
```

> **⚠️ IMPORTANT**: Tags must be created from the `main` branch only. The release workflow will fail if the tag is not on `main`.

### 4. Monitor Release Workflow

1. Go to: https://github.com/cduggn/ccexplorer/actions
2. Find the "goreleaser" workflow for your tag
3. Monitor the workflow execution

The workflow will:
- ✅ Verify tag is on `main` branch
- ✅ Build binaries for all platforms (Linux, macOS, Windows)
- ✅ Build multi-architecture binaries (amd64, arm64, arm)
- ✅ Create Docker image and push to GHCR
- ✅ Update Homebrew tap
- ✅ Generate checksums
- ✅ Create GitHub release with artifacts

### 5. Verify Release

After workflow completes successfully:

1. **Check GitHub Release**
   - Go to: https://github.com/cduggn/ccexplorer/releases
   - Verify release notes are generated
   - Check all binary artifacts are attached

2. **Check Docker Image**
   ```bash
   docker pull ghcr.io/cduggn/ccexplorer:v0.7.0
   docker run --rm ghcr.io/cduggn/ccexplorer:v0.7.0 --help
   ```

3. **Check Homebrew** (may take a few minutes)
   ```bash
   brew update
   brew install cduggn/cduggn/ccexplorer
   ```

4. **Test Binary Download**
   - Download binary for your platform from GitHub releases
   - Verify checksum
   - Test execution

### 6. Update Documentation

After release:

1. Update README.md if needed (installation instructions, version references)
2. Announce release (if applicable):
   - Slack/Discord/Community channels
   - Social media
   - Blog post (for major releases)

## What Gets Released

Each release includes:

| Artifact | Description | Location |
|----------|-------------|----------|
| **Binaries** | Pre-compiled executables | GitHub Releases |
| - Linux (amd64, arm64, arm) | Native Linux binaries | - |
| - macOS (amd64, arm64) | Universal binary for macOS | - |
| - Windows (amd64) | Windows executable | - |
| **Docker Image** | Container image (linux/amd64) | ghcr.io/cduggn/ccexplorer |
| **Homebrew Formula** | macOS package | cduggn/homebrew-cduggn |
| **Checksums** | SHA256 checksums | checksums.txt in release |
| **Source Code** | Tarball and zip | GitHub automatic |

## Troubleshooting

### Release Workflow Fails: "Tag is not on main branch"

**Problem**: You created the tag from a feature branch.

**Solution**:
```bash
# Delete the tag
git tag -d v0.7.0
git push origin :refs/tags/v0.7.0

# Switch to main and recreate
git checkout main
git pull origin main
git tag -a v0.7.0 -m "Release v0.7.0"
git push origin v0.7.0
```

### Release Workflow Fails During Build

**Problem**: Build or test failures.

**Solution**:
1. Check the workflow logs for specific errors
2. Fix the issues on `main` branch
3. Delete the failed tag and recreate after fixes are merged

### Docker Image Push Fails

**Problem**: Authentication or permission issues with GHCR.

**Solution**:
1. Verify `GHCR_PUBLISHER` secret is set correctly
2. Check that the token has `write:packages` permission
3. Verify repository settings allow package publishing

### Homebrew Formula Not Updated

**Problem**: Homebrew tap repository not accessible.

**Solution**:
1. Check that `cduggn/homebrew-cduggn` repository exists
2. Verify GoReleaser has write access to the tap repository
3. Check `.goreleaser.yaml` configuration for tap settings

## Emergency Rollback

If a release has critical issues:

### Option 1: Delete Release (Before Users Download)

```bash
# Delete tag locally and remotely
git tag -d v0.7.0
git push origin :refs/tags/v0.7.0
```

Then delete the GitHub release from the web UI.

### Option 2: Hotfix Release (After Users Downloaded)

1. Fix the issue on `main`
2. Create a new patch release (e.g., `v0.7.1`)
3. Mark the problematic release as "pre-release" in GitHub UI
4. Add warning to release notes

### Option 3: Mark as Pre-release

In GitHub releases UI, edit the release and check "Set as a pre-release".

## CI/CD Pipeline

### On Pull Requests

When you open a PR to `main`, the CI workflow runs:

- ✅ Linting (golangci-lint)
- ✅ Tests (Go 1.22 and 1.23)
- ✅ Race condition detection
- ✅ Build verification (all platforms)
- ✅ Security scanning (Trivy)

All checks must pass before merging.

### On Main Branch Push

After merging to `main`:
- Same CI checks run
- CodeQL security analysis
- OSSF Scorecard analysis

### On Tag Push

When a `v*` tag is pushed to `main`:
- Release workflow triggers
- GoReleaser builds and publishes all artifacts

## Security Considerations

### Current Protections

- ✅ Releases only from `main` branch
- ✅ Automated branch verification
- ✅ Security scanning (CodeQL, Trivy, Codacy)
- ✅ Dependency scanning
- ✅ Docker image published to GHCR (private registry)

### Future Enhancements (Planned)

- ⏳ Artifact signing with Cosign
- ⏳ SBOM (Software Bill of Materials) generation
- ⏳ SLSA provenance attestation
- ⏳ Multi-arch Docker images
- ⏳ Automated changelog generation
- ⏳ Release notes automation

See [RELEASE_PROCESS_EVALUATION.md](../RELEASE_PROCESS_EVALUATION.md) for detailed roadmap.

## Configuration Files

### `.goreleaser.yaml`

Main GoReleaser configuration defining:
- Build targets and architectures
- Docker image configuration
- Homebrew tap settings
- Changelog generation rules

### `.github/workflows/release.yml`

GitHub Actions workflow that:
- Triggers on `v*` tags
- Verifies tag is on `main` branch
- Sets up Go environment with caching
- Executes GoReleaser
- Publishes artifacts

### Makefile

Development shortcuts:
- `make build` - Local build with CGO_ENABLED=0
- `make test` - Run tests
- `make lint` - Run linters
- `make release` - Local release testing (snapshot mode)

**Note**: Do NOT use `make tag` for production releases. Use the manual process above to ensure proper branch verification.

## Release Checklist

Before creating a release tag:

- [ ] All PRs for this release are merged to `main`
- [ ] CI is passing on `main` branch
- [ ] Version number determined (semantic versioning)
- [ ] CHANGELOG updated (if maintained manually)
- [ ] Documentation updated for new features
- [ ] Breaking changes documented
- [ ] Migration guide written (if needed for breaking changes)
- [ ] Local testing completed
- [ ] On `main` branch locally

Create and push tag:

- [ ] Tag created with `git tag -a vX.Y.Z -m "Release vX.Y.Z"`
- [ ] Tag pushed with `git push origin vX.Y.Z`
- [ ] Release workflow started successfully
- [ ] Release workflow completed successfully

Post-release verification:

- [ ] GitHub release created
- [ ] All binaries attached to release
- [ ] Docker image available on GHCR
- [ ] Homebrew formula updated
- [ ] Checksums verified
- [ ] Sample binary downloaded and tested
- [ ] Docker image pulled and tested
- [ ] Release announced (if applicable)

## Getting Help

If you encounter issues with the release process:

1. Check this documentation first
2. Review recent successful releases for reference
3. Check GoReleaser documentation: https://goreleaser.com/
4. Review GitHub Actions logs for errors
5. Open an issue in the repository

## References

- [GoReleaser Documentation](https://goreleaser.com/)
- [Semantic Versioning](https://semver.org/)
- [GitHub Releases](https://docs.github.com/en/repositories/releasing-projects-on-github)
- [Release Process Evaluation](../RELEASE_PROCESS_EVALUATION.md) - Detailed analysis and future improvements

---

*Last Updated: 2026-01-08*
*Process Owner: Project Maintainers*
