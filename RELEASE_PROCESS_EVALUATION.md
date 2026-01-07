# Release Process Evaluation & Recommendations

## Executive Summary

Your current release process has several issues that deviate from industry best practices. This document outlines critical problems, recommended improvements, and alternatives to consider.

---

## Current State Analysis

### What You Have

**✅ Good:**
- GoReleaser configured for multi-platform builds (Linux, macOS, Windows)
- Multi-architecture support (amd64, arm64, arm)
- Docker image publishing to GHCR
- Homebrew tap for easy installation
- Security scanning (CodeQL, OSSF, Dependency Review, Codacy)
- Checksums for binary verification

**❌ Critical Issues:**

1. **No Branch Protection on Releases**
   - `.github/workflows/release.yml:4-6` triggers on ANY tag `v*` from ANY branch
   - Tags created on feature branches will trigger production releases
   - No validation that releases come from `main` branch only

2. **No CI/CD Pipeline for PRs**
   - No automated testing before merge
   - No build validation on pull requests
   - Security scans run but no integration tests

3. **Manual Version Management**
   - `Makefile:73-77` requires manual tag creation
   - Error-prone process (typos, wrong version numbers)
   - No semantic versioning automation
   - No changelog generation automation

4. **Incomplete Release Artifacts**
   - No SBOM (Software Bill of Materials)
   - No release notes automation
   - No signing of binaries or container images
   - No artifact attestation

5. **Cache Clearing is Wasteful**
   - `.github/workflows/release.yml:29-32` clears both module and build caches
   - Significantly slows down release builds
   - Should use GitHub Actions cache instead

6. **GoReleaser Configuration Issues**
   - `.goreleaser.yaml:13` has `CGO_ENABLED=0` but `Makefile:50,55` sets `CGO_ENABLED=1`
   - Inconsistent build configuration between local and release
   - Missing Windows-specific configurations (icons, manifests)
   - No APT/RPM package generation for Linux

7. **Docker Image Issues**
   - Only builds for `linux/amd64` (line 45-46 in `.goreleaser.yaml`)
   - Should support multi-arch (arm64) given your binary builds do
   - Uses `whoami` in Dockerfile which won't work properly
   - No vulnerability scanning for images

---

## Industry Best Practices You're Missing

### 1. Automated Semantic Versioning

**Top Projects Use:**
- Conventional Commits (e.g., `feat:`, `fix:`, `BREAKING CHANGE:`)
- Automated version bumping based on commit messages
- Tools: semantic-release, release-please, or changesets

**Example:** [Kubernetes](https://github.com/kubernetes/kubernetes), [Terraform](https://github.com/hashicorp/terraform)

### 2. Comprehensive CI/CD Pipeline

**Missing Workflows:**

```
┌─────────────┐
│   PR        │
│   Workflow  │──► Lint, Test, Build, Security Scan
└─────────────┘

┌─────────────┐
│   Main      │
│   Workflow  │──► Integration Tests, Build Artifacts
└─────────────┘

┌─────────────┐
│   Release   │
│   Workflow  │──► Version, Changelog, Sign, Publish
└─────────────┘
```

**Top Projects:** [k6](https://github.com/grafana/k6), [Traefik](https://github.com/traefik/traefik)

### 3. Supply Chain Security

**SLSA Level 3 Requirements:**
- Build provenance
- Signed artifacts
- SBOM generation
- Reproducible builds

**Tools:**
- [sigstore/cosign](https://github.com/sigstore/cosign) - Signing
- [anchore/syft](https://github.com/anchore/syft) - SBOM generation
- [slsa-framework/slsa-github-generator](https://github.com/slsa-framework/slsa-github-generator)

**Example:** [Docker Buildx](https://github.com/docker/buildx), [Helm](https://github.com/helm/helm)

### 4. Automated Changelog

**Best Practice:**
- Auto-generated from conventional commits
- Grouped by type (features, fixes, breaking changes)
- Links to PRs and issues
- Contributors list

**Tools:**
- [git-chglog](https://github.com/git-chglog/git-chglog)
- [github-changelog-generator](https://github.com/github-changelog-generator/github-changelog-generator)
- Built into semantic-release/release-please

### 5. Pre-release Testing

**What's Missing:**
- Canary releases
- Beta/RC releases
- Testing against real AWS accounts in CI
- Performance benchmarking

### 6. Comprehensive Release Artifacts

**Modern Projects Include:**
- Binaries for all platforms
- Container images (multi-arch)
- Package manager formats (Homebrew ✅, APT, RPM, Chocolatey, Scoop)
- Checksums ✅
- **SBOM** (JSON/SPDX format)
- **Signatures** (GPG or Sigstore)
- **Attestations** (SLSA provenance)

---

## GoReleaser vs Alternatives

### GoReleaser (Current)

**Pros:**
- Purpose-built for Go projects
- Excellent multi-platform support
- Rich ecosystem (Homebrew, Docker, etc.)
- Well-documented
- Used by major projects (Hugo, Terraform providers, K6)

**Cons:**
- Go-specific (not polyglot)
- Less flexibility for complex pipelines
- Configuration can be verbose
- Limited to release-time (not full CI/CD)

**Verdict:** ✅ **Keep GoReleaser** - It's the industry standard for Go CLIs

### Dagger

**Pros:**
- Polyglot (useful if you add other languages)
- Programmable CI/CD (write in Go, Python, TypeScript)
- Local reproducibility
- Containerized builds
- Composable pipelines

**Cons:**
- Learning curve (new paradigm)
- Younger ecosystem than GoReleaser
- Requires more custom code
- Overkill for single-language projects
- Less "batteries included" for releases

**Use Dagger If:**
- You need complex multi-language builds
- You want to write CI logic in Go
- You need the same pipeline locally and in CI
- You're building a platform/framework

**Verdict:** ❌ **Not Recommended** for your use case - stick with GoReleaser for releases, but consider Dagger for CI testing

### Other Alternatives Considered

**goreleaser/release-action** ✅ (Already using - good choice)

**Earthly** - Similar to Dagger, same assessment

**Mage** - Build tool, not release tool

**GitHub Releases API directly** - Too low-level, reinventing the wheel

---

## Recommended Solution Architecture

### Option A: Enhanced GoReleaser (Recommended)

**Keep GoReleaser** for releases, add:

1. **CI Pipeline** (test on every PR)
2. **Automated versioning** (release-please or semantic-release)
3. **Supply chain security** (SBOM, signing)
4. **Branch protection** (releases from main only)

**Implementation Complexity:** Low-Medium
**Timeline:** 1-2 days
**Maintenance:** Low

### Option B: Hybrid Dagger + GoReleaser

**Use Dagger** for CI/testing, **GoReleaser** for releases

**Pros:**
- Reproducible test environment
- Flexibility for future needs

**Cons:**
- Higher complexity
- More maintenance
- Longer learning curve

**Implementation Complexity:** Medium-High
**Timeline:** 3-5 days
**Maintenance:** Medium

---

## Detailed Recommendations

### 1. Fix Critical Security Issue ⚠️ **URGENT**

**Problem:** Releases can be triggered from any branch

**Solution:**
```yaml
# .github/workflows/release.yml
on:
  push:
    tags:
      - "v*"
    branches:
      - main  # Only allow releases from main
```

Better yet, use a separate release trigger:
```yaml
on:
  release:
    types: [created]
```

Then use release-please to create releases automatically.

### 2. Implement Automated Versioning

**Option A: release-please (Recommended)**

Creates PRs that:
- Bump version based on conventional commits
- Generate changelog
- Create GitHub release when PR merged
- Used by Google, many GCP projects

**Setup:**
```yaml
# .github/workflows/release-please.yml
name: release-please
on:
  push:
    branches:
      - main

jobs:
  release-please:
    runs-on: ubuntu-latest
    steps:
      - uses: googleapis/release-please-action@v4
        with:
          release-type: go
```

**Option B: semantic-release**

More configurable, JavaScript-based, used by npm ecosystem.

### 3. Add Comprehensive CI Pipeline

**Create `.github/workflows/ci.yml`:**
```yaml
name: CI

on:
  pull_request:
  push:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
          cache: true
      - run: make lint
      - run: make test
      - run: make test-race
      - run: make build

  integration-test:
    runs-on: ubuntu-latest
    steps:
      # Add AWS integration tests with mock credentials
      - run: make test-integration

  security:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: aquasecurity/trivy-action@master
        with:
          scan-type: 'fs'
          scan-ref: '.'
```

### 4. Add Supply Chain Security

**Update `.goreleaser.yaml`:**
```yaml
# Add SBOM generation
sboms:
  - artifacts: archive

# Add signing with Cosign
signs:
  - cmd: cosign
    args:
      - "sign-blob"
      - "--output-signature=${signature}"
      - "--output-certificate=${certificate}"
      - "${artifact}"
      - "--yes"
    artifacts: all

# Add Docker image signing
docker_signs:
  - cmd: cosign
    args:
      - "sign"
      - "${artifact}@${digest}"
      - "--yes"
    artifacts: all
```

**Add to workflow:**
```yaml
- uses: sigstore/cosign-installer@v3
- uses: anchore/sbom-action/download-syft@v0
```

### 5. Improve Docker Configuration

**Multi-arch images:**
```yaml
# .goreleaser.yaml
dockers:
  - image_templates:
      - 'ghcr.io/cduggn/ccexplorer:{{ .Tag }}-amd64'
      - 'ghcr.io/cduggn/ccexplorer:{{ .Tag }}-arm64'
    goos: linux
    goarch: amd64
    use: buildx
    build_flag_templates:
      - "--platform=linux/amd64"
  - image_templates:
      - 'ghcr.io/cduggn/ccexplorer:{{ .Tag }}-arm64'
    goos: linux
    goarch: arm64
    use: buildx
    build_flag_templates:
      - "--platform=linux/arm64"

docker_manifests:
  - name_template: 'ghcr.io/cduggn/ccexplorer:{{ .Tag }}'
    image_templates:
      - 'ghcr.io/cduggn/ccexplorer:{{ .Tag }}-amd64'
      - 'ghcr.io/cduggn/ccexplorer:{{ .Tag }}-arm64'
  - name_template: 'ghcr.io/cduggn/ccexplorer:latest'
    image_templates:
      - 'ghcr.io/cduggn/ccexplorer:{{ .Tag }}-amd64'
      - 'ghcr.io/cduggn/ccexplorer:{{ .Tag }}-arm64'
```

### 6. Add More Package Formats

**Linux packages:**
```yaml
# .goreleaser.yaml
nfpms:
  - id: packages
    package_name: ccexplorer
    homepage: https://github.com/cduggn/ccexplorer
    maintainer: Colin Duggan <duggan.colin@gmail.com>
    description: AWS Cost Explorer CLI tool
    license: MIT
    formats:
      - deb
      - rpm
      - apk
```

**Windows packages:**
```yaml
scoops:
  - repository:
      owner: cduggn
      name: scoop-bucket
    homepage: https://github.com/cduggn/ccexplorer
    description: AWS Cost Explorer CLI tool
    license: MIT
```

### 7. Fix CGO Inconsistency

**Decision needed:**
- If you need CGO (for SQLite, C libraries): Set `CGO_ENABLED=1` everywhere
- If you don't need CGO: Set `CGO_ENABLED=0` everywhere (recommended for pure Go)

**Recommendation:** `CGO_ENABLED=0` for better portability

### 8. Remove Cache Clearing

**Replace in `.github/workflows/release.yml`:**
```yaml
# Remove these:
# - name: Clear module cache
#   run: go clean -modcache
# - name: Clear build cache
#   run: go clean -cache

# Add this instead:
- uses: actions/cache@v4
  with:
    path: |
      ~/.cache/go-build
      ~/go/pkg/mod
    key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
```

### 9. Add Release Documentation

**Create `.github/RELEASE.md`:**
```markdown
# Release Process

## Automated Releases

1. Merge PRs to main using conventional commits:
   - `feat: add new feature` → minor version bump
   - `fix: resolve bug` → patch version bump
   - `feat!: breaking change` or `BREAKING CHANGE:` in body → major bump

2. release-please bot creates a release PR automatically

3. Review the release PR:
   - Check generated changelog
   - Verify version number
   - Review included commits

4. Merge release PR → automated release triggers

5. GoReleaser publishes:
   - GitHub release with notes
   - Binaries for all platforms
   - Docker images
   - Homebrew formula update
   - Checksums and signatures

## Manual Release (Emergency)

Only if automation fails:

1. Create tag: `git tag -a v1.2.3 -m "Release v1.2.3"`
2. Push tag: `git push origin v1.2.3`
3. Monitor: https://github.com/cduggn/ccexplorer/actions
```

### 10. Add Missing Workflows

**Branch protection:**
```yaml
# Set in GitHub UI or using Terraform/code:
# - Require PR reviews
# - Require CI to pass
# - No direct pushes to main
# - Require signed commits (optional)
```

**Dependabot:**
```yaml
# .github/dependabot.yml
version: 2
updates:
  - package-ecosystem: gomod
    directory: /
    schedule:
      interval: weekly
    groups:
      aws:
        patterns:
          - "github.com/aws/aws-sdk-go-v2*"

  - package-ecosystem: github-actions
    directory: /
    schedule:
      interval: weekly
```

---

## Migration Plan

### Phase 1: Critical Fixes (Week 1)

**Priority: URGENT**

1. ✅ Fix release workflow branch restriction
2. ✅ Add CI pipeline for PRs
3. ✅ Fix CGO_ENABLED inconsistency
4. ✅ Add caching to workflows
5. ✅ Document current release process

### Phase 2: Automation (Week 2)

**Priority: HIGH**

1. ✅ Implement release-please
2. ✅ Remove manual tagging from Makefile
3. ✅ Add conventional commit validation
4. ✅ Set up branch protection rules
5. ✅ Add Dependabot

### Phase 3: Security (Week 3)

**Priority: MEDIUM**

1. ✅ Add SBOM generation
2. ✅ Implement artifact signing (Cosign)
3. ✅ Add Docker image scanning
4. ✅ Generate SLSA provenance
5. ✅ Add security policy (SECURITY.md)

### Phase 4: Distribution (Week 4)

**Priority: LOW**

1. ✅ Add multi-arch Docker images
2. ✅ Add APT/RPM packages
3. ✅ Add Windows package managers (Scoop/Chocolatey)
4. ✅ Add installation verification tests

---

## Comparison: Top Go CLI Projects

| Feature | ccExplorer | kubectl | k6 | terraform | hugo | Your Target |
|---------|------------|---------|----|-----------| -----|-------------|
| GoReleaser | ✅ | ❌ | ✅ | ❌ | ✅ | ✅ |
| Auto Versioning | ❌ | ✅ | ✅ | ✅ | ✅ | ✅ |
| PR CI | ❌ | ✅ | ✅ | ✅ | ✅ | ✅ |
| SBOM | ❌ | ✅ | ✅ | ✅ | ❌ | ✅ |
| Signing | ❌ | ✅ | ✅ | ✅ | ❌ | ✅ |
| Multi-arch Docker | ❌ | ✅ | ✅ | ✅ | ✅ | ✅ |
| Package Managers | 1 | 5+ | 4+ | 5+ | 5+ | 4+ |
| Release Notes | Manual | Auto | Auto | Auto | Auto | Auto |
| Changelog | ❌ | ✅ | ✅ | ✅ | ✅ | ✅ |

---

## Cost-Benefit Analysis

### Staying with GoReleaser + Improvements

**Effort:** 2-3 days
**Benefits:**
- Fix critical security issues
- Automated versioning saves hours per release
- Better supply chain security
- More professional project image
- Easier contributor onboarding

**Cost:**
- Initial setup time
- Learning new tools (release-please, cosign)

**ROI:** Very High ⭐⭐⭐⭐⭐

### Switching to Dagger

**Effort:** 5-7 days
**Benefits:**
- Reproducible local testing
- More control over pipeline
- Future-proof for polyglot

**Cost:**
- Significant learning curve
- Ongoing maintenance complexity
- Lose GoReleaser's release features
- Would still need release tooling

**ROI:** Low ⭐⭐ (for your current needs)

---

## Final Recommendation

### Keep GoReleaser, Enhance the Pipeline

**Immediate Actions (This Week):**
1. ✅ Fix release workflow to require main branch
2. ✅ Add CI pipeline for PRs
3. ✅ Fix CGO_ENABLED to be consistent
4. ✅ Add workflow caching

**Next Sprint:**
1. ✅ Implement release-please for automated versioning
2. ✅ Add SBOM generation
3. ✅ Add artifact signing with Cosign
4. ✅ Improve Docker multi-arch support

**Later:**
1. Add more package manager support (APT, RPM, Scoop)
2. Add integration test suite
3. Consider Dagger only if you need complex multi-language builds

---

## Questions to Answer

Before proceeding, decide:

1. **CGO:** Do you need CGO_ENABLED=1? (Check if you use any C dependencies)
2. **Signing:** GPG keys or Keyless signing (Cosign/Sigstore)?
3. **Versioning:** Comfortable with conventional commits?
4. **Testing:** Can you add AWS integration tests with mock credentials?
5. **Packages:** Which package managers are priority? (APT/RPM for Linux users?)

---

## References

**Best Practices:**
- [SLSA Framework](https://slsa.dev/)
- [OpenSSF Best Practices](https://bestpractices.coreinfrastructure.org/)
- [Conventional Commits](https://www.conventionalcommits.org/)

**Tools:**
- [GoReleaser](https://goreleaser.com/)
- [release-please](https://github.com/googleapis/release-please)
- [Cosign](https://github.com/sigstore/cosign)
- [Syft](https://github.com/anchore/syft)

**Example Projects:**
- [k6 releases](https://github.com/grafana/k6/blob/master/.github/workflows/release.yml)
- [Traefik releases](https://github.com/traefik/traefik/blob/master/.github/workflows/release.yaml)
- [Helm releases](https://github.com/helm/helm/blob/main/.github/workflows/release.yml)

---

*Generated: 2026-01-07*
*Evaluated by: Claude Code*
