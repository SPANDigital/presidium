---
title:  The Big Picture
---


Presidium is an open-source, static website generator built on Hugo for creating and maintaining software documentation for Agile teams and their users. Being static simplifies publication - you don't need to provision a server and can host your documentation in a matter of seconds on Github Pages or Bitbucket. You can easily include documentation publication in your continuous integration pipeline.

Presidium gives Agile teams a way to evangelize, explain, and support solutions to their peers so they can be correctly used and successfully adopted. Wikis can be unstructured and lose their value, tribal knowledge and content stored on Slack channels or emails can go into a black hole. Engineers are experts at writing code, not documentation. *Software documentation that doesn't suck* is not just a tagline.

Presidium is based on sound theoretical strategies for developing learning content and managing knowledge assets. The core of this approach is the development of specific focused micro-articles that explain individual concepts. Predefined article templates and a menu structure get you up and running quickly.

In addition to the template for software documentation, Presidium includes templates for:
* On-boarding
* Design systems

Presidium supports:

* [Automatic Menus]({{site.baseurl}}/key-concepts/#menu-structure-and-behavior): The left navigation menu is dynamically created every time you publish your site.

* [Versioning]({{site.baseurl}}/recipes/versioning/): Support for multiple documentation versions. You can easily switch between versions without losing context.

* [Documentation Importers]({{site.baseurl}}/reference/importers/): Import documentation from source code and specs using the [Swagger]({{site.baseurl}}/reference/importers/#swagger), [Javadoc]({{site.baseurl}}/reference/importers/#javadoc) and [JSDoc]({{site.baseurl}}/reference/importers/#jsdoc) importers.

* [Theming]({{site.baseurl}}/reference/#themes): Easily select and modify [Bootswatch](https://bootswatch.com/) themes that meet your brand requirements.

* [Role Filtering]({{site.baseurl}}/reference/front-matter/#user-roles): Define user roles and filter site content by a specific selected role.

* [Article Status Tracking]({{site.baseurl}}/reference/#directories): Track the status of an article and manage simple authoring workflows.

* [Link Validation]({{site.baseurl}}/tools/): Make sure your links actually go somewhere!

Presidium was created by [SPAN Digital](http://www.spandigital.com) and is licensed under [Apache 2.0](/updates/#license)
