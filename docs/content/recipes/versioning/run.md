---
title: Run
weight: "3"
---

To publish a specific version:

```sh
$ npm run gh-pages -- -v 1.2
```

To publish the latest version:

```sh
$ npm run gh-pages
```

# Versioning Conventions

* Presidium will display only the latest five semantically versioned releases of your documentation.
* Presidium only supports numeric semantic versions.
* Presidium names the latest version as 'latest'. This means that if you have versions 1.0, 1.1, 1.2, 1.3, 1.5, you will see: latest, 1.5, 1.4, 1.3 and 1.2.

# Known Issues

* The versioning UI component is not yet responsive.
* A user cannot view versions by serving the documentation locally.
* Turning off versioning does not prevent access to the version if the url is known.
* Currently, only publishing to Github pages is supported.
* If the use of a CNAME is disabled, it is not removed from `.versions`when republished.
