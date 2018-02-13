---
title: Link Validation
---
A **link** **validation** tool is provided so that you can validate all the links in your site.
The tool will indicate which links are valid, broken and/or external and will provide a warning for potentially broken links.

# Run Link Validator

```
$ npm run validate
```

# Output

The validator tool will produce the following output for each unique link found in your site:

## Valid Links:

```
VALID:          /recipes/structure/
VALID:          /recipes/structure/#nested-articles
```

## External Links:

```
EXTERNAL:       http://bootswatch.com/
```

## Broken Links:

```
BROKEN:         /broken-link
BROKEN:         /recipes/structure/#incorrect-anchor
```

## Potentially Broken Links:

```
WARNING:        /recipes/structure is missing a trailing '/'
```

# Known Issues:

1. The link validator will currently mark any assets served out of `/media` as BROKEN.