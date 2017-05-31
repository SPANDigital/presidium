---
title: Hosting Presidium in project /docs folder
author: torabisu
roles: Developer
---

Presidium does not have to be hosted in a repository of its own.  It can exist within your software project in a sub folder, e.g. `/docs`

This recipe will show you how to accomplish this.

# Getting Presidium

The easiest way to incorpate Presidium into your project is to get the latest
archived version from [https://github.com/SPANDigital/presidium-template/archive/master.zip](https://github.com/SPANDigital/presidium-template/archive/master.zip) and
uncompressing it into your projects empty /docs folder.  If you wish to clone the project instead
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

Now you can follow the [Getting Started Guide]

# Splitting Away From Master

Firstly you need to ignore the published folder from within your projects master branch:

```bash
$ echo "docs/node_modules/" >> .gitignore
$ echo "docs/.jekyll/" >> .gitignore
$ echo "docs/dist/site" >> .gitignore
```

Next we can set up the gh-pages branch:

```bash
$ git checkout --orphan gh-pages
$ git reset --hard
$ git commit --allow-empty -m "Initializing gh-pages branch"
$ git push origin gh-pages
$ git checkout master
```

This creates a new branch called gh-pages, and it will be managed totally separate
from your existing tree.

We can how book out the gh-pages branch into the dist/sites folder:

```bash
$ rm -rf docs/dist/site
$ git worktree add -B gh-pages docs/dist/site origin/gh-pages
```

# Reference

[https://gohugo.io/tutorials/github-pages-blog/](https://gohugo.io/tutorials/github-pages-blog/)