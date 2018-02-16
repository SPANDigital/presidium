---
title: Link Validation
---
The **link** **validation** tool enables you to validate all the links in your site.
The tool indicates which links are valid, broken, and/or external, and provides a warning for potentially broken links.

# Run Link Validator

```
$ npm run validate
```

# Output

The link validation tool produces the following output for each unique link found in your site:

## Valid Links

```
VALID:          /recipes/structure/
VALID:          /recipes/structure/#nested-articles
```

## External Links

```
EXTERNAL:       http://bootswatch.com/
```

## Broken Links

```
BROKEN:         /broken-link
BROKEN:         /recipes/structure/#incorrect-anchor
```

## Potentially Broken Links

```
WARNING:        /recipes/structure is missing a trailing '/'
```

# Known Issues

* The link validation tool currently marks any assets from `/media` as BROKEN.
