---
title: Configure
---

# Update Variables In _config.yml

Presidum supports the use of simple variables in `_config.yml`. In order to use versioning you must make sure that any helper variables that include the base url of your site are defined using a `${...}` variable. For example, this is incorrect:

```yaml
baseurl: /myRepositoryName
code_examples: /myRepositoryName/media/code_examples
```

This is correct:

```yaml
baseurl: /myRepositoryName
code_examples: ${baseurl}/media/code_examples
```

Be careful to not introduce circular dependencies with the use of variables.

# Set The Base Url & Turn On Versioning

In order to use Presidium versioning on `gh-pages`, you are required to use the base url of the repository. You must ensure that this is set in `_config.yml` in the root of your project:

```yaml
baseurl: /myRepositoryName
...

# Optional Support For Versioning
versioned: true
```

Note the lack of a trailing slash.