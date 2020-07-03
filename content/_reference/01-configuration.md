---
title: Configuration
---

This is a sample Presidium configuration file:

```yml
#
# Site Metadata
#
# - name: Site name
# _ baseurl: Optional URL to use if documentation is hosted in a subdirectory `domain.com/{baseurl}`
# - footer: Footer copy
# - logo: Menu bar logo.
# - show: Hide or show article components. Defaults to true
# - external: Links to external sources
#    - authors-url: base url for article authors
# NOTE: If you have a baseurl be sure to append it to any static content path e.g.: ${baseurl}/media/images/logo.png
name: Presidium Template
#baseurl: /mysite
footer: Template Footer
logo: /media/images/logo.png

show:
  status: true
  author: true
external:
  authors-url: https://github.com/

#
# Jekyll collections.
#
# Collections represent the top level grouping of articles for a site.
# Collections directories require a leading underscore: "./content/_{collection-name}/"
#
collections:
    overview:
    key-concepts:
    prerequisites:
    getting-started:
    best-practices:
    reference:
    glossary:
    recipes:
    tools:
    updates:

#
# Site Sections
#
# Describes the top level structure of your site
# - tile: Title for top level menu title
# - path: Path to generated article collection
# - collection: The Jekyll collection to use for generating a sub menu of articles (optional).
# - collapsed: Whether articles should show in a tree or collapsed into a single item
# - always-expanded: Boolean indicating whether the section should be expanded by default and   never collapse
#
sections:
  - title: "Overview"
    url: "/"
    collection: overview
    collapsed: true

  - title: "Key Concepts"
    url: "/key-concepts/"
    collection: key-concepts

  - title: "Prerequisites"
    url: "/prerequisites/"
    collection: prerequisites

  - title: "Getting Started"
    url: "/getting-started/"
    collection: getting-started
    collapsed: true

  - title: "Best Practices"
    url: "/best-practices/"
    collection: best-practices

  - title: "Reference"
    url: "/reference/"
    collection: reference

  - title: "Glossary"
    url: "/glossary/"
    collection: glossary
    export-content: true

  - title: "Recipes"
    url: "/recipes/"
    collection: recipes

  - title: "Tools"
    url: "/tools/"
    collection: tools

  - title: "Updates"
    url: "/updates/"
    collection: updates

#
# Optional Role Filters
#

#roles:
#  label: "Show documentation for"
#  all: "All Roles"
#  options:
#    - "Business Analyst"
#    - "Developer"
#    - "Tester"

# Optional support for versioning.
versioned: false

sass:
    sass_dir: media/css/_sass
```