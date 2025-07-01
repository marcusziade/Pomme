# GitHub Actions Configuration

## Required Secrets

For the release workflow to work properly, you need to configure the following secrets in your repository settings:

### HOMEBREW_TAP_TOKEN

This token is required for GoReleaser to push the Homebrew formula to your homebrew-tap repository.

1. Go to https://github.com/settings/tokens/new
2. Create a new token with the following permissions:
   - `repo` (Full control of private repositories)
3. Copy the token
4. Go to your Pomme repository settings → Secrets and variables → Actions
5. Create a new repository secret named `HOMEBREW_TAP_TOKEN` with the token value

Without this token, the release will succeed but the Homebrew formula won't be updated automatically.