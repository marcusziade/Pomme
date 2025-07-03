# AUR Submission Guide for pomme

This guide documents how pomme is maintained on the AUR (Arch User Repository).

## Package Information

- **Stable package**: `pomme` - Tracks stable releases
- **Development package**: `pomme-git` - Tracks the latest git commits

## AUR URLs

- https://aur.archlinux.org/packages/pomme
- https://aur.archlinux.org/packages/pomme-git

## Maintaining the Packages

### Updating for New Releases

When a new version is released:

```bash
# Clone the AUR repository
git clone ssh://aur@aur.archlinux.org/pomme.git
cd pomme

# Update the pkgver in PKGBUILD
vim PKGBUILD  # Update pkgver=X.X.X

# Regenerate .SRCINFO
makepkg --printsrcinfo > .SRCINFO

# Test the build locally
makepkg -si

# Commit and push
git add PKGBUILD .SRCINFO
git commit -m "Update to version X.X.X"
git push
```

### Important Notes

1. **Directory Name**: The GitHub tarball extracts to `Pomme-X.X.X` (capital P)
2. **Tag Format**: Tags are without the 'v' prefix (e.g., `2.0.0` not `v2.0.0`)
3. **Tests**: Currently disabled in check() due to build issues that need to be fixed upstream

### Package Structure

The PKGBUILD:
- Downloads from `https://github.com/marcusziade/pomme/archive/${pkgver}.tar.gz`
- Extracts to `Pomme-${pkgver}` directory
- Builds the binary with version information
- Installs binary, license, and documentation

### Testing Installation

Users can install with:
```bash
# Stable version
yay -S pomme

# Development version
yay -S pomme-git
```

### Troubleshooting

If users report issues:
1. Check that the GitHub release exists
2. Verify the tarball URL is correct
3. Test the build locally with `makepkg -si`
4. Check for any new dependencies

### Orphaning

If you can no longer maintain the package, orphan it on the AUR web interface so someone else can adopt it.