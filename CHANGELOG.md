# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Added Hugo theme module dependencies: `presidium-layouts-base`, `presidium-layouts-blog`, and `presidium-styling-base`
- Added automated pre-build cleanup for `presidium hugo` command (removes `public`, `resources`, `go.sum`, `.hugo_build.lock`, `themes`)
- Added automated module management before Hugo builds (`hugo mod clean`, `go mod tidy`, `hugo mod tidy`, `hugo mod get`)
- Added default Hugo flags: `--templateMetrics --ignoreCache --logLevel info`

### Changed
- Updated embedded Hugo from v0.87.0 to v0.148.0+extended
- Updated Go version to 1.25
- Migrated from deprecated `packr2` to `go:embed` for template embedding
- Added `--tags extended` to build target to enable SCSS support via LibSass
- Enhanced `presidium hugo` command to execute pre-build cleanup and module management automatically

### Fixed
- Fixed compilation errors caused by Hugo API changes (replaced deprecated `paths.IsAbsURL`)
- Fixed build errors by removing `packr2` dependency

### Removed
- Removed `toolchain` directive from `go.mod`
- Removed `packr2` dependency and related files
- Removed `packrd/packed-packr.go` and `pkg/domain/service/template/template-packr.go`

### Technical Notes
- Template `go.mod` files renamed to `go.mod.tpl` to work around `go:embed` restrictions
- `ProcessTemplate` updated to automatically rename `go.mod.tpl` back to `go.mod` during generation

### Known Issues
- `presidium-layouts-base@v0.10.0` has compatibility issues with Hugo v0.148+
- Navigation templates assume `.File` is always available, causing nil pointer errors for certain page types
- Workaround: Override `layouts/partials/navigation/nav-item.html` in your project or wait for `presidium-layouts-base` update
- See `.tmp/HUGO_COMPATIBILITY_ISSUE.md` for detailed fix
