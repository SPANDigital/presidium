---
title: Github Pages
---

[Github Pages](https://pages.github.com/) provides a quick and convenient way of hosting and serving your documentation from a Github repository.
1. Checkout the branch that you would like to publish from.
1. Commit and push all your changes to your Github repository.
1. Run `npm run gh-pages` to build and push a generated site to a `gh-pages` branch in your repository.
1. You will need to [enable gh-pages](https://help.github.com/articles/configuring-a-publishing-source-for-github-pages/) 
in your repository settings to make your documentation available online.
1. You may set the `cname` property in your `_config.yml` to export a custom domain name with your generated site. This will require [further configuration](https://help.github.com/articles/using-a-custom-domain-with-github-pages/).
