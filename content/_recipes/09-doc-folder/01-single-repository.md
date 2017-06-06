---
title: Single Repository
author: torabisu
roles: Developer
---

Presidium does not have to be hosted in a repository of its own.  It can exist within your software project in a 
sub folder, e.g. `/docs`

## Getting Presidium

The easiest way to incorporate Presidium into your project is to get the [latest archived version](https://github.com/SPANDigital/presidium-template/archive/master.zip) and uncompress it into your project's empty `/docs` folder.  If you wish to clone the project instead,
be sure to remove the `.git` folder from `docs/` after cloning as your project repo will manage
the `docs/` folder and **not** the presidium template repo.

The contents of your `docs/` folder should look something like this:

```bash
LICENSE
NOTICE
README.md
_config.yml
content/
dist/
media/
package.json
```

Add the following to your project's `.gitignore` file:

```bash
$ echo "docs/node_modules/" >> .gitignore
$ echo "docs/.jekyll/" >> .gitignore
$ echo "docs/dist/" >> .gitignore
$ echo "docs/dist/.versions" >> .gitignore
```

From this point onwards you can follow the getting started guide, where your documentation root is simply in `/docs`.
