# Presidium

This repository tracks Presidium's [user documentation](http://presidium.spandigital.net/overview/) and provides a sample of how presidium is used.

# What is Presidium?

[Presidium](http://presidium.spandigital.net) is a software documentation management system for agile teams and their users made from the stuff software engineers love. Presidium uses familiar tools already in use by many software development teams to easily compose and publish and manage quality documentation.

Presidium integrates [Hugo](https://github.com/gohugoio/hugo) as a Go library to provide static site generation capabilities compatible with our Presidium Themes.

Read Presidium's [documentation](http://presidium.spandigital.net/overview/) for details on all features and for help on setting up a new project using the [sample template](https://github.com/SPANDigital/presidium-template).

# How do I use Presidium?

For detailed instructions on getting started, please see [Getting Started](http://presidium.spandigital.net/getting-started/).

# Features and Issues

All upcoming features and issues are tracked [here](https://github.com/SPANDigital/presidium/issues).

# Development

Presidium integrates Hugo as a Go library, ensuring that the user is building their docs with a Hugo version that is compatible with our Presidium Themes.

Common Hugo CLI commands are available directly as Presidium commands (e.g., `presidium server`), and all Hugo commands are available through `presidium hugo`.

## Theme Embedding

Presidium themes are embedded in the binary as a compressed zip archive to ensure offline functionality and avoid Go module embedding restrictions:

1. **Themes are Git Submodules**: The three Presidium themes (styling-base, layouts-base, layouts-blog) are tracked as git submodules for easy updates
2. **Build-time Compression**: During `make build`, themes are zipped into `themes.zip` and embedded into the binary using Go's `//go:embed` directive
3. **Runtime Extraction**: When presidium runs, the zip is extracted to a temporary directory and Hugo is configured to use the local themes via `HUGO_MODULE_REPLACEMENTS`
4. **No Network Required**: Once built, the presidium binary requires no network access to use the embedded themes
5. **Fast Extraction**: Zip extraction is typically 10-50ms, adding negligible overhead

### Why Zip Instead of Direct Embedding?

Go's `//go:embed` directive cannot embed directories that contain `go.mod` files (other Go modules). The themes are separate Go modules with their own `go.mod` files, so we:

- Package them as a zip at build time
- Embed the single zip file (no module restrictions)
- Extract at runtime (fast and simple)

## Getting started

Run the build command:

```
make build
```

Move the binary into a valid documentation site directory, then run the hugo build cmd:

```
./presidium hugo 
```

Or if you want to serve the site after the hugo build, then run:

```
./presidium hugo server
```

## Testing

### Unit Tests

Run the full test suite:

```
make test
```

### Offline Build Test

Verify that the binary works without network access using Docker:

```
make test-offline
```

This test:

- Prepares the themes.zip bundle from git submodules
- Builds Presidium inside a Docker container using multi-stage build (Go 1.25 + extended Hugo with LibSass)
- Runs `presidium hugo` in an isolated container with `--network none`
- Verifies all expected output files are generated (index.html, presidium.js, links.js, assets, images)
- Confirms no network access was attempted

The multi-stage Docker build ensures proper Linux binary compilation with all required C++ dependencies for SCSS support.

**Requirements**: Docker must be installed and running.
