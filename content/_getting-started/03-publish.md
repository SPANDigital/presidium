---
title: Publishing your Site
---

# Github Pages
To publish using Github Pages, commit and push your site to a Github repository and run the following:
```sh
$ npm run gh-pages
```
This will push your generated site to a gh-pages branch in your repository. You will need to 
[enable gh-pages](https://help.github.com/articles/configuring-a-publishing-source-for-github-pages/) 
in your repository settings to make your site public.

# Static Site
The generated static site can be found in `dist/site` if you choose to host your site elsewhere.