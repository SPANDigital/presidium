---
title: Link Validation
---
The **link** **validation** tool enables you to validate all the links in your site.
The tool indicates which links are valid, broken, and/or external, and provides a warning for potentially broken links.

# Run Link Validator

Basic usage:

```
$ npm run validate
```

## Output

The link validation tool produces the following output for each unique link found in your site:

### Valid Links

```
VALID:          /recipes/structure/
VALID:          /recipes/structure/#nested-articles
```

### External Links

```
EXTERNAL:       http://bootswatch.com/
```

### Broken Links

```
BROKEN:         /broken-link
BROKEN:         /recipes/structure/#incorrect-anchor
```

### Potentially Broken Links

```
WARNING:        /recipes/structure is missing a trailing '/'
```

# Configuration 

The link validator supports multiple arguments: 

| Argument         | Options                                        | Default value | Description                                                                                          | 
| ---------------- |----------------------------------------------- | ------------- | -----------------------------------------------------------------------------------------------------| 
| `fail_on_errors` | `true`, `false`                                | `true`        | Specify if the process should produce an error (and stop) when invalid links are detected.           |
| `log`            | `all`, `valid`, `broken`, `warning`, `external`| `all`         | Limits the logging of console messages based on the parameter specified.                             |
| `check`          | `author`                                       |  N/A          | Specifies which front-matter properties should be validated - only _author_ is currently supported.   |


Full example: 

```sh
$ npm run validate -- --fail_on_errors=false --log=broken --check=author 
```

# Known Issues

* The link validation tool currently marks any assets from `/media` as BROKEN.
