---
title: Host
---

The simplest way to publish a Github repository is to use Github Pages, but you can also host the generated site on any web server.

## On Github Pages

Commit and push your site to a Github repository and run the following:

```sh
$ npm run gh-pages
```

This will push your generated site to a `gh-pages` branch in your repository. You will need to 
[enable gh-pages](https://help.github.com/articles/configuring-a-publishing-source-for-github-pages/) 
in your repository.

## As a Static Site

The generated static site can be found in `dist/site`.