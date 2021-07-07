---
title: Github Pages
---

[Github Pages](https://pages.github.com/) provides a quick and convenient means of hosting and serving your documentation from a Github repository.
1. Checkout the branch that you would like to publish from.
2. Commit and push all your changes to your Github repository.
3. Run `npm run gh-pages` to build and push a generated site to a `gh-pages` branch in your repository.
4. You will need to [enable gh-pages](https://help.github.com/articles/configuring-a-publishing-source-for-github-pages/) 
in your repository settings to make your documentation available online.
5. You may host your documentation using a [custom domain](https://help.github.com/articles/using-a-custom-domain-with-github-pages/). Setting the `cname` property in your `_config.yml` exports your custom domain with your generated site.

> Make sure the **baseurl** variable in `_config.yml` is set to match your repository's name (and path if applicable). For example, `#baseurl: /presidium-site`.  
