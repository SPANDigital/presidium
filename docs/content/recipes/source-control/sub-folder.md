---
title: Sub Folder
weight: "2"
---

Presidium can exist within a sub-folder of an existing source code repository, for example, `/docs`.

## Getting Presidium

The easiest way to incorporate Presidium into your project is to get the [latest archived version](https://github.com/SPANDigital/presidium-template/archive/master.zip) and uncompress it into your project's empty `/docs` folder.  If you would rather clone the project, make sure you remove the `.git` folder from `docs/` after cloning because your project repo will manage
the `docs/` folder and **not** the Presidium template repo.

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

```
docs/node_modules/
docs/.jekyll/
docs/dist/
docs/dist/.versions
```

From this point on, you can follow instructions in the Getting Started section. The only difference is that your documentation root is `/docs`.