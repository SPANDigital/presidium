---
title: The Big Picture
weight: 1
status: NEEDS WORK
---

> **OUTDATED** : What is referred to as `Presidium` below is called 'Presidium Writing Tools' in Separation of Concerns, which helps a Presidium user set up an initial doc site.
>   
> `Presidium` is an open-source, static website generator built on Hugo for creating and maintaining software 
> documentation for Agile teams and their users. Being static simplifies publication - you don’t need to provision 
> a server and can host your documentation in a matter of seconds on Github Pages.
> You can easily include documentation publication in your continuous integration pipeline.

Presidium gives Agile teams a way to evangelize, explain, and support solutions to their peers so they can be 
correctly used and successfully adopted. Wikis can be unstructured and lose their value, tribal knowledge and 
content stored on Slack channels or emails can go into a black hole. Engineers are experts at writing code, not 
documentation. *Software documentation that doesn’t suck* is not just a tagline.

Presidium is based on sound theoretical strategies for developing learning content and managing knowledge assets. 
The core of this approach is the development of specific focused micro-articles that explain individual concepts. 
Predefined article templates and a menu structure get you up and running quickly.

In addition to the content template for software documentation, the Presidium writing tools includes content templates for:
- Services (default)
- On-boarding
- Design systems
- Blog

Presidium supports:
- [Automatic Menus]({{< ref "key-concepts/#menu-structure" >}}): The left navigation menu is dynamically created every time you publish your site.
- [Theming]({{< ref "reference/#themes" >}}): Easily select supported themes
- [Role Filtering]({{< ref "reference/front-matter/#user-roles" >}}): Define user roles and filter site content by a specific selected role.
- [Article Status Tracking]({{< ref "reference/front-matter/#statuses" >}}): Track the status of an article and manage simple authoring workflows.
- [Link Validation]({{< ref "tools/#link-validation" >}}): Make sure your links actually go somewhere!
<!-- - [Versioning: Support] for multiple documentation versions. You can easily switch between versions without losing context. // TODO insert correct link to versionin article -->

Presidium was created by [SPAN Digital](http://www.spandigital.com/) and is licensed under [Apache 2.0](https://github.com/SPANDigital/presidium/blob/develop/LICENSE)
