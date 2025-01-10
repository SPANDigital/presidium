---
title: Setup
weight: "1"
status: NEEDS MORE INFO
---

Presidium supports semantic versioning of your documentation when publishing to Github pages.

# Github Pages

Presidium versioning relies on the use of the `gh-pages` branch to publish current and versioned documentation.

To use Presidium versioning, you will need to make you sure have a  `gh-pages` branch set up in your repository.
Skip these steps if you already have a branch set up.

1. Change the directory to the root of your project:

  ```sh
cd path/to/your/project
```

  If the branch already exists (be sure you know what you are doing), you can clean `gh-pages` and start again:

  ```
$ git push origin --delete gh-pages
$ git branch -D gh-pages
```

1. Create a new `gh-pages` branch with nothing in it:

  ```sh
$ git checkout --orphan gh-pages && git reset --hard
$ git commit --allow-empty --m "Initialize gh-pages"
$ git push -u origin gh-pages
$ git checkout master
```

# Update .gitignore

Presidium uses a hidden folder as a staging area for publishing. To make sure you don't accidentally commit this folder, run the following command:

```sh
$ echo ".versions" >> .gitignore
```
