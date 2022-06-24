---
title: Configure
weight: "2"
---

# Update Variables in _config.yml

Presidum supports the use of simple variables in `_config.yml`. To use versioning, you must make sure that any helper variables that include the base url of your site are defined using a `${...}` variable.

For example, this is correct:

```yaml
baseurl: /myRepositoryName
code_examples: ${baseurl}/media/code_examples
```
This is incorrect:

```yaml
baseurl: /myRepositoryName
code_examples: /myRepositoryName/media/code_examples
```

Make sure you don't introduce circular dependencies when using variables.

# Set the Base URL & Turn on Versioning

To use Presidium versioning on `gh-pages`, you must use the base url of the repository. Make sure that this is set in `_config.yml` in the root of your project. (Note the lack of a trailing slash.)

```yaml
baseurl: /myRepositoryName
...

# Optional Support For Versioning
versioned: true
```
