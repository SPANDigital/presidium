---
title: Set Up
---

Presidium supports semantic versioning of your documentation when publishing to Github pages.

# Github Pages

Presidium versioning relies on the use of the `gh-pages` branch to publish current and versioned documentation.

In order to use Presidium versioning, you will need to ensure you have a gh-pages branch set up in your repository. 
You may skip these step if you already have a branch set up.

Change the directory to the root of your project:

```sh
cd path/to/your/project
```

If the branch already exists (be sure you know what you are doing), you can clean `gh-pages` and start again:

```
$ git push origin --delete gh-pages
$ git branch -D gh-pages
```

Create a new `gh-pages` branch with nothing on it:

```sh
$ git checkout --orphan gh-pages && git reset --hard
$ git commit --allow-empty --m "Initialize gh-pages"
$ git push -u origin gh-pages
$ git checkout master
```