---
title: config.yml
weight: 30
---
The `config.yml` file
* Sets specifications for the top-level sections, including their title and order
* Sets module-wide specifications, such as the module title, language, copyright text, and others

This file can use many parameters, as described in [the Hugo documentation](https://gohugo.io/), but in Presidium it typically contains the following sections and fields:

#### `markup` Section
Configures the markup parser. The default is Goldmark. For use of other parsers see [the Hugo documentation](https://gohugo.io/getting-started/configuration-markup/).

#### `params` Section
These are some of the more common fields used in this section:

* `favicon`: Path to the favicon image. Favicon and logo are stored in the `static/images` folder.
* `logo`: Path to the logo image (displayed at the top of the left navigation bar)
* `sortByFilePath`: If set to `true`, navigation is sorted alphabetically instead of by `weight`
<!--* `quality_category`: Enterprise only -->

#### `menu` Section
Typically contains only a subsection called `main` that defines the top-level sections of the module's left-side navigation menu.
The fields are:
* `identifier`: Name of the directory associated with the menu item
* `name`: Section title as it appears when rendered
* `url`: String appended to the module URL to identify this section
* `weight`: Order of this section relative to other top-level sections

#### `outputformats` Section
Defines the formats to be generated when building the module. All the necessary formats are already set up by SPAN. For custom output formats, refer to https://gohugo.io/methods/page/outputformats/#article

#### `outputs` Section
Defines what parts of the site use specific output formats. As with `outputformats`, the necessary setup is already done by SPAN.

#### `module` Section
Typically contains only a subsection called `imports`, which brings in styling and layouts and, optionally, other externally-stored resources. For example: 
```  imports:
  - path: github.com/spandigital/presidium-styling-base
  - path: github.com/spandigital/presidium-layouts-base
```

#### `frontmatter` Section
Defines what frontmatter fields are available in articles, as well as the required data type for each field and whether the field is required or not.

> **Note:** If the field is marked as required but not present in an article, the site throws an error when attempting to build. If the field is not required and not present, the site throws a warning, but still builds successfully.

#### Other Items
These are not contained in a section.

* `baseUrl`: The URL used for the \{\{% baseurl %}} shortcode
* `assetDir`: Directory where the site assets are stored
* `pluralizelisttitles`: True/False. Automatically pluralize menu titles. Defaults to `false`
* `enableGitInfo`: Enables the ability to retrieve Git information of the last commit made to articles


