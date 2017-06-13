---
title: Single Repository
author: torabisu
roles: Developer
---

Presidium does not have to be hosted in a repository of its own.  It can exist within your software project in a 
sub folder, e.g. `/docs`

## Getting Presidium

The easiest way to incorpate Presidium into your project is to get the latest archived version from [https://github.com/SPANDigital/presidium-template/archive/master.zip](https://github.com/SPANDigital/presidium-template/archive/master.zip) and uncompressing it into your projects empty /docs folder.  If you wish to clone the project instead
just make sure you remove the .git folder from docs/ as we want your application repo to manage
the docs/ folder and not the presidium template repo.

The contents of your docs/ folder should look something like this:

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
```

From this point onwards you can follow the getting started guide, where your documentation root is simply in `/docs`
