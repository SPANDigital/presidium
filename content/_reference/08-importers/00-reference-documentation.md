---
title: Importing Documentation
---

The Reference section of Presidium should be reserved as a low-level reference for users.
Typical components documented in the Reference section may include a client library or API specification.
The process of importing documentation involves parsing a reference source and generating articles that are included in the generated site.

>Where possible, reference documentation should be generated to ensure that your documentation is in sync with the
system being documented.

Presidium supports the following documentation sources:
- [javadoc comments](#javadoc)
- [jsdoc comments](#jsdoc)
- [swagger api](#swagger)

For other sources that do not yet have an importer, documentation can be [embedded](#embed) into references.
