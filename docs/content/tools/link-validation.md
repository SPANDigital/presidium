---
title: Link Validation
url: link-validation
slug: link-validation
---
This tool reports links within a Presidium site. You can use this tool see which links are broken, points to non standard protocols.

> The tool does not need a live Presidium site to work against, so the site does not need to be deployed.

## How to use the tool

1. Generate a local public site by running presidium hugo:

   ```shell
   presidium hugo
   ```

2. Next run the tool pointing to the `public` site:

   ```shell
   presidium report pagelinks ./public
   ```

## Inspect the report

For example:

```shell
presidium validate ./public
VALIDATION PATH: ./public

        total: 864
  valid links: 864
       broken: 0
     external: 0
     warnings: 0
```

Broken links, external and dynamic (and/or non standard) links will be reported in detail.