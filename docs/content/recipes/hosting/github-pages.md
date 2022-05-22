---
title: Github Pages
weight: 1
---

[GitHub Pages](https://pages.github.com/) provides a quick and convenient means of hosting and serving your 
documentation from a Github repository. The recommended way to host your Presidium site in Github Pages is to use
[GitHub Actions](https://github.com/features/actions)

If you are using GitHub Actions for the first time you can simply create a `push.yaml` file under `.github/workflows`
and define your GitHub Pages job, if you already using GitHub Actions, you can simply add a job to deploy your
documentation and configure your `baseURL` accordingly:

```yaml
name: Presidium Github Pages
on:
  - push

jobs:
  gh-pages:
    runs-on: ubuntu-20.04
    concurrency:
      group: ${{ github.workflow }}-${{ github.ref }}
    steps:
      - uses: actions/checkout@v2
        with:
          submodules: true
          fetch-depth: 0
          persist-credentials: false
      - name: Setup Hugo
        uses: peaceiris/actions-hugo@v2
        with:
          hugo-version: '0.87.0'
          extended: true
      - name: Build
        run: hugo --minify --baseURL `https://<ORGANIZATION>.github.io/<REPOSITORY_NAME>`
      - name: Deploy
        uses: peaceiris/actions-gh-pages@v3
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          publish_dir: ./public
```

If your documentation site uses private repository modules, you need to change the git config in order to be able to 
fetch from private repositories:

```yaml
    ...
    steps:
      - uses: actions/checkout@v2
        with:
          submodules: true
          fetch-depth: 0
          persist-credentials: false
      - name: Inject Golang Git Config
        run: git config --global url."https://${PRIVATE_GITHUB_TOKEN}:@github.com/".insteadOf "https://github.com/"
    ...
    env:
      PRIVATE_GITHUB_TOKEN: ${{ secrets.PRIVATE_GITHUB_TOKEN }}
      GOPRIVATE: github.com/spandigital
```

Note that we are using the environment variable `PRIVATE_GITHUB_TOKEN` which you need to create in your repository
`Settings` > `Secrets`, containing a personal access token capable of downloading your private resources.
