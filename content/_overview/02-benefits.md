---
title:  Benefits
author:  virtualtraveler
---
# Easy Installation

Presidium uses Github pages to publish content. If you have access to Github and [Github Pages](https://pages.github.com/) is turned on for your site, you have everything you need. All you need to do it write your content and commit changes, and your content is published automatically. It's really that easy.

# Micro-article Based

Presidium is built on sound theoretical strategies for developing learning content and managing knowledge assets. The core of this approach is the development of specific focused micro-articles that explain individual concepts. This splits the problem of content development into two parts:

1. Defining a site outline based on the concepts that must be documented and how they should be sequenced and presented.
1. Developing articles that each explain each concept, using a standard template to ensure consistency.

Micro-articles break the documentation development and maintenance activities into manageable chunks that can be prioritized based on need. In many cases, it's possible to write and review an article in under 20 minutes.

This Presidium documentation site not only explains the Presidium software but also includes best practice guidelines for writing your documentation.

# Easy Content Maintenance

Presidium is designed to make it easy to write *and* maintain your documentation. Instead of a monolithic document, Presidium manages a directory structure of micro-articles. The directory structure maps to the sections of your site, while the articles are joined together to build the content of each section. You can create a directory to add new sections and sub-sections. Each section contains a sample article template you can use as a starting point when creating new content. This site explains other [best practices]({{site.base_url}}/best-practices/) for writing articles.

If you know how to  maintain a codebase in a Github repository, you already know how to maintain Presidium documentation. You can use your team's current Github workflow to manage the approval and publication process and use Github’s features to handle merge conflicts and other issues, just as you do for your source code.

# Configurable Features

Presidium is shipped pre-configured for software documentation. You can easily modify the suggested site structure and behavior. Most teams only need to remove a few sections that are not relevant to their project. Presidium has several configurable features that can be enabled as needed.

# Serverless

Presidium uses Github for both content management and as a publication server. If you have access to Github and Github pages is turned on for your repository, then you have all you need to start publishing. You don't need to deploy a server or get permission from some dark corner of the enterprise. All you need to do is commit the branch that contains your documentation and Github pages using Jekyll, and Presidium will do the rest.

Despite having no server, there are enough features in the Presidium.js component to meet the needs of most teams. But, if you really need more features, Presidium is open-source and we welcome contributions from the community.
