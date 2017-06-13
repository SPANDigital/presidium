---
title: Semantic Versioning
author: dfollettsmith
---

Presidium supports semantic versioning of your documentation when publishing to Github pages.

# Branch Setup

Presidium versioning relies on the use of the gh-pages branch to publish current and versioned documentation.
In order to use Presidium versioning, you will need to ensure you have a gh-pages branch set up in your repository. 
You may skip these step if you already have a branch set up.

Change the directory to the root of your project:

```sh
cd path/to/your/project
```

If the branch already exists (be sure you know what you are doing!), you can clean gh-pages and start again:

```
$ git push origin --delete gh-pages
$ git branch -D gh-pages
```

Create a new gh-pages branch with nothing on it:

```sh
$ git checkout --orphan gh-pages && git reset --hard
$ git commit --allow-empty --m "Initialize gh-pages"
$ git push -u origin gh-pages
$ git checkout master
```

# Update .gitignore

Presidium uses a hidden folder as a staging area for publishing. To make sure you don't accidentally commit this run the following command:

```sh
$ echo ".versions" >> .gitignore
```

# Set The Base Url & Turn On Versioning

In order to use Presidium versioning on gh-pages, you are required to use the base url of the repository. You must ensure that this is set in `_config.yml` in the root of your project:

```yaml
baseurl: /myRepositoryName
...

# Optional Support For Versioning
versioned: true
```

Note the lack of a trailing slash.

# Update Variables In _config.yml

Presidum supports the use of simple variables in _config.yml. In order to use versioning you must make sure that any helper variables that include the base url of your site are defined using a `${...}` variable. For example, this is incorrect:

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

# Running The Tool

To publish a version:

```sh
$ npm run gh-pages -- -v 1.2
```

If you want to publish the latest version:

```sh
$ npm run gh-pages
```

# Versioning Conventions

1. Presidium will display only the latest five semantically versioned releases of your documentation.
1. Presidium supports only numeric semantic versions.
1. Presidium names the latest version as 'latest', this means if you have 1.0, 1.1, 1.2, 1.3, 1.5, versions, then you will see: latest, 1.5, 1.4, 1.3 and 1.2.

# Known Issues:

1. Versioning UI component is not yet responsive.
1. A user cannot view versions by serving the documentation locally.
1. Turning off versioning does not prevent access to the version if the url is known.
1. Currently, only publishing to Github pages is supported.
1. If the use of a CNAME is disabled, it is not removed from `.versions` on a subsequent republish.
1. When running Presidium locally the following error is reported: `../versions.json' not found.` on the console.

