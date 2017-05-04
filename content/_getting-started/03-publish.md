---
title: Publishing your Site
---

# Github Pages
To publish using Github Pages, push your site to a Github repository and run the following:
```sh
$ npm run gh-pages
```
This will commit a generated site to a gh-pages branch. You will need to 
[enable gh-pages](https://help.github.com/articles/configuring-a-publishing-source-for-github-pages/) 
in your repository to make your site public.

# Static Site
The generated static site can be found in `dist/site` if you choose to host your site elsewhere.