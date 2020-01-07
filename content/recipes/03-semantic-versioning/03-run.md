---
title: Run
---

To publish a version:

```sh
$ npm run gh-pages -- -v 1.2
```

If you want to publish the latest version:

```sh
$ npm run gh-pages
```

# Versioning Conventions

1. Presidium will display only the latest five semantically versioned releases of your documentation.
1. Presidium supports only numeric semantic versions.
1. Presidium names the latest version as 'latest', this means if you have 1.0, 1.1, 1.2, 1.3, 1.5, versions, then you will see: latest, 1.5, 1.4, 1.3 and 1.2.

# Known Issues

1. Versioning UI component is not yet responsive.
1. A user cannot view versions by serving the documentation locally.
1. Turning off versioning does not prevent access to the version if the url is known.
1. Currently, only publishing to Github pages is supported.
1. If the use of a CNAME is disabled, it is not removed from `.versions` on a subsequent republish.
