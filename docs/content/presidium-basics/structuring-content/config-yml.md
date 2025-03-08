---
title: config.yml
weight: 30
---
The `config.yml` file
* Sets specifications for the top-level sections, including their title and order
* Sets docset-wide specifications, such as the docset title, language, 

### `markup` Section
Used to configure the markup parser. Will most likely be the same for all Presidium docs.
Source: https://gohugo.io/getting-started/configuration-markup/

### `params` Section
* frontmatter: https://docs.spandigital.net/docs/presidium-docs-internal/official-features/paved-path/hugo/#frontmatter-validation
  * key
  * type
  * strict

  <span style="color:purple">needs more explanation. don't understand how the linked page relates.</span>
* enterprise_key: Key that is used in the URL in Enterprise instances; for example, `docs.spandigital.net/docs/{key}/overview`
* favicon: Path to the favicon image
* logo: Path to the logo image (displayed at the top of the left navigation bar)
* quality_category: 
* sortByFilePath: If set to `true`, navigation is sorted alphabetically instead of by `weight`

### `menu` Section
Typically contains only a subsection called `main` that defines the top-level sections of the docset.
The fields are:
* identifier: Name of the directory associated with the menu item
* name: Section title as it appears when rendered
* url: String appended to the docset URL to identify this section
* weight: order of this section relative to other top-level section

### `outputformats` Section
Defines the formats that should be generated when building docs. All the necessary formats are already set up by SPAN. For custom output formats, refer to https://gohugo.io/methods/page/outputformats/#article

### `outputs` section
Defines what parts of the site use specific output formats. As with `outputformats`, the necessary setup is already done by SPAN.

### `module` Section
Typically contains only a subsection called `imports`, which brings in styling and layouts and, optionally, other externally-stored resources. For example, 
```  imports:
  - path: github.com/spandigital/presidium-styling-base
  - path: github.com/spandigital/presidium-layouts-base
```

### `frontmatter` Section
Defines what frontmatter fields are available in articles, as well as the required data type for each field and whether the field is required or not. If the field is marked as required but not present on an article, the site will throw an error when attempting to build. If the field is not required and not present, the site will throw a warning, but still successfully build.

### Other Items
These are not contained in a section.

* baseUrl: The URL that will be used for the {{% baseurl %}} shortcode
* assetDir: Directory where the site assets are expected to be
* pluralizelisttitles: True/False. Automatically pluralize menu titles. Defaults to false
* enableGitInfo: Enables the ability to retrieve git information of the last commit made to articles


