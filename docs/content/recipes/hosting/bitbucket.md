---
title: Bitbucket
weight: "2"
draft: true
status: UNKNOWN (very old content)
---

Bitbucket hosting is provided by setting up a separate repository to host your generated site. To serve your content on Bitbucket:

1. Create a hosting repository named `<your username>.bitbucket.io` or `<your team name>.bitbucket.io`.
2. Build the static site: `hugo`
3. Copy the generated site from `public` to the hosting repository and push to Bitbucket
4. The content will be served from: `<your username>.bitbucket.io` or `<your team name>.bitbucket.io` 

> Custom domain names are [currently not supported by Bitbucket](https://bitbucket.org/site/master/issues/3641/custom-domain-repo-url-without-user-name) by Bitbucket.
