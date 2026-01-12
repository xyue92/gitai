# Release Guide

This guide explains how to create a new release for GitAI using the automated build system.

## Overview

The release process is fully automated using GitHub Actions. When you push a new tag, the system will:
1. Build binaries for all platforms (Linux, macOS, Windows)
2. Generate SHA256 checksums
3. Create a GitHub Release
4. Upload all binaries and checksums

## Prerequisites

Before your first release, you need to:

1. **Update repository URLs** in the following files:
   - Replace `xyue92/gitai` with your actual GitHub username/organization
   - Files to update:
     - `.github/workflows/release.yml`
     - `scripts/install.sh` (line 5)
     - `scripts/install.ps1` (line 4)
     - `homebrew/gitai.rb`
     - `README.md` (all download URLs)

2. **Test the build locally** (optional but recommended):
   ```bash
   # Build all platforms
   ./scripts/build.sh v0.1.0

   # Check output
   ls -lh dist/
   ```

## Creating a Release

### Step 1: Prepare for Release

1. Make sure all changes are committed and pushed to `main` branch
2. Update version numbers if needed in documentation
3. Test the application locally

### Step 2: Create and Push a Tag

```bash
# Create a new tag (use semantic versioning)
git tag -a v1.0.0 -m "Release v1.0.0: Initial release with multi-platform support"

# Push the tag to GitHub
git push origin v1.0.0
```

**Tag naming convention:**
- Use semantic versioning: `v{MAJOR}.{MINOR}.{PATCH}`
- Examples: `v1.0.0`, `v1.2.3`, `v2.0.0-beta.1`
- Always prefix with `v`

### Step 3: Monitor the Build

1. Go to your repository on GitHub
2. Click on "Actions" tab
3. You should see a "Release" workflow running
4. Wait for all builds to complete (usually 3-5 minutes)

### Step 4: Verify the Release

Once the workflow completes:

1. Go to "Releases" section in your GitHub repository
2. You should see a new release with your tag
3. Verify that all binaries are attached:
   - `gitai-linux-amd64` + `.sha256`
   - `gitai-linux-arm64` + `.sha256`
   - `gitai-darwin-amd64` + `.sha256`
   - `gitai-darwin-arm64` + `.sha256`
   - `gitai-windows-amd64.exe` + `.sha256`

### Step 5: Test the Installation Scripts

Test that users can download and install:

```bash
# Test Unix install script
curl -sSL https://raw.githubusercontent.com/xyue92/gitai/main/scripts/install.sh | bash

# Verify
gitai --version
```

### Step 6: Update Homebrew Formula (if using Homebrew)

See [homebrew/README.md](homebrew/README.md) for detailed instructions.

Quick steps:
1. Download each binary and calculate SHA256
2. Update `Formula/gitai.rb` in your homebrew-tap repository
3. Test with `brew install gitai`

## Release Checklist

Use this checklist for each release:

- [ ] All changes committed and pushed to main
- [ ] Version numbers updated in docs (if applicable)
- [ ] Local testing completed
- [ ] Tag created with proper version number
- [ ] Tag pushed to GitHub
- [ ] GitHub Actions workflow completed successfully
- [ ] All binaries present in GitHub Release
- [ ] Checksums file (checksums.txt) uploaded to release
- [ ] Install script tested
- [ ] Self-update command tested (`gitai update --check`)
- [ ] Homebrew formula updated (if applicable)
- [ ] Release announcement written (optional)

## Troubleshooting

### Build Failed

Check the GitHub Actions logs:
1. Go to Actions tab
2. Click on the failed workflow
3. Expand the failed step to see error details

Common issues:
- **Missing permissions**: Make sure the workflow has `contents: write` permission
- **Build errors**: Check if code compiles locally with `go build`
- **Invalid tag**: Make sure tag follows `v*` pattern

### Binaries Not Uploaded

Check that the `release` job completed:
- The `build` job creates artifacts
- The `release` job downloads artifacts and creates the release
- Both jobs must succeed

### Checksum Verification Fails

The checksum file format should be:
```
SHA256_HASH  filename
```

If format is wrong, check the checksum generation step in the workflow.

## Manual Release (Fallback)

If GitHub Actions is not working, you can create a release manually:

```bash
# Build all platforms
./scripts/build.sh v1.0.0

# Create checksums
cd dist/
sha256sum * > checksums.txt

# Create a GitHub Release manually and upload files
```

## Version Numbering Guide

Follow semantic versioning:

- **MAJOR** version (1.0.0 → 2.0.0): Breaking changes, incompatible API changes
- **MINOR** version (1.0.0 → 1.1.0): New features, backwards compatible
- **PATCH** version (1.0.0 → 1.0.1): Bug fixes, backwards compatible

Examples:
- First release: `v1.0.0`
- Add new feature: `v1.1.0`
- Fix a bug: `v1.0.1`
- Breaking change: `v2.0.0`

## Tips

1. **Test locally first**: Always run `./scripts/build.sh` before creating a tag
2. **Use annotated tags**: `git tag -a` creates annotated tags with metadata
3. **Write good release notes**: Explain what changed in the release
4. **Delete failed tags**: If a release fails, delete the tag:
   ```bash
   git tag -d v1.0.0
   git push origin :refs/tags/v1.0.0
   ```
5. **Pre-releases**: Use tags like `v1.0.0-beta.1` for testing

## Automation Ideas

Future improvements you can add:

1. **Automated Homebrew updates**: Create a workflow in homebrew-tap repo that updates the formula automatically
2. **Release notes generation**: Use conventional commits to auto-generate release notes
3. **Binary size optimization**: Add UPX compression to reduce binary size
4. **Checksums in release body**: Auto-include checksums in the release description
5. **Download statistics**: Track how many times each binary is downloaded
