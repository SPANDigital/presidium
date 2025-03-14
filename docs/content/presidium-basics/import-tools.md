---
title: Import Tools
weight: 60
---
The process of importing documentation involves parsing a reference source and generating articles that are included in the generated site. This can be useful if you want to add a client library, API specification, or similar reference content to a docset. If the reference content will continue to be maintained outside of Presidium, periodic imports can keep the two synchronized.

Presidium supports the following documentation sources:
OpenAPI3
JSON
HTML
RST
Embedding

OpenAPI3

Presidium includes a Golang tool (presidium-oapi3) for importing your OpenAPI 3 spec into Presidium documentation.
Add the presidium-oapi3 dependency to your site’s package.json or run npm install --save presidium-oapi3.
Add a script that invokes the tool.
Run npm run import-oapi whenever you need to update your API documentation.
{
   "scripts" : {
       "import-oapi" : "presidium-oapi3"
   },
}

Example:
$ npm run import-oapi convert -f <YOUR_API_SPEC> -o <THE_OUTPUT_DIRECTORY> -r <THE_PRESIDIUM_REFERENCE_URL>

The following options are available for presidium-oapi3:
Option
Description
-n, --apiName string
The name under which the generated docs will be grouped
-f, --file string
OpenAPI 3 spec file
-o, --outputDir string
The output directory
-r, --referenceURL string
The reference URL (default “reference”)
-h, --help
help for convert


HTML to Markdown converter

html2md

html2md is a tool that allows you to convert HTML files into Presidium markdown articles.
Installation
Usage
Advanced Usage
Limitations

Installation
Install from Homebrew
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh)"
Add SPAN’s Homebrew tap
brew tap SPANDigital/homebrew-tap https://github.com/SPANDigital/homebrew-tap.git
Install html2md
brew install html2md

Usage
Usage:
 html2md convert [source] [dest] [flags]

Flags:
 -d, --debug             enable debug logging
     --headers strings   article header tags (default [h1,h2])
     --select string     the part of the page to select and convert (default "body")
Params
source is the url of the website or the path to the local html file(s).
dest is the path to the directory where the converted markdown files will be saved.
Examples
Download and convert the Presidium website
html2md convert https://presidium.spandigital.net/ ./presidium --select="#presidium-content"
Convert local html files
html2md convert ./html-files ./presidium --select=.article --headers=h1
Note
This converter does not output all the files needed to build a complete Presidium website. You will still need to create a Presidium site first, or import the converted markdown files into an existing site. The converted markdown files will be saved in the content directory relative to the dest path along with assets (images, videos, etc) in the assets directory.

Advanced usage
You can define a config.yml file in your working directory or in$HOME/.html2md directory with additional options to control the conversion.
Remove HTML
The html.remove option allows you to selectively remove elements from your source document before converting it to Markdown. This can be useful when you want to clean up your HTML content or remove unwanted elements that have no relevance.
Example Usage:
html:
 remove: ['.nav-link', '#warning']
In this example, all element matching the CSS selectors .nav-link or #warning will be removed before conversion.
Replace HTML
The html.replace option allows you to transform specific HTML elements into custom Markdown syntax, giving you greater control over the appearance and structure of your converted document.
Example Usage:
Suppose you have HTML content that includes tooltips, and you want to convert them into Markdown-friendly syntax. Here’s how you can achieve this using the “Replace HTML” feature:

```
html:
 replace:
 - match: '.tooltips-term' # CSS selector for the element to be replaced
   select: ['?href', '.tooltips-text'] # ?href selects the href attribute of the matched element, .tooltips-text selects the content of a child element with the class "tooltips-text"
   replace: '{{</* tooltip "$1" text="$2" */>}}' # $1 and $2 are the selected elements
```

In this example, the CSS selector .tooltips-term is specified as the element to be replaced. The select field allows you to capture specific attributes and content relative to the matched element. Finally, the replace pattern converts the selected elements into a Markdown tooltip format.
Note
In addition to the standard CSS-selectors, select allows you to select attributes on the matched element using the ? prefix. You can also get the text of the matched element by passing text in the select list.
Replace Markdown
Similar to the HTML replace you can also replace markdown based on a Regex pattern.
Example Usage:
Let’s say you have Markdown content with links, and you want to prepend the base URL to all these links.
markdown:
 replace:
   - pattern: '\[([^]]+)\]\(([^\)]+)\)'
     with: "[$1]({{< baseurl >}}$2)"
Sample config
html:
 remove: ['.article-title .permalink', '#warning'] # CSS selectors for elements that should be removed before conversion.
 replace:
   - match: '.tooltips-term' # CSS selector for the element to be replaced.
     # Below are the arguments to select elements relative to the matched element.
     select: ['?href', '.tooltips-text']
     # Replacement pattern with the selected arguments.
     replace: '{{</* tooltip "$1" text="$2" */>}}'
markdown:
 replace:
   - pattern: '\[([^]]+)\]\(([^\)]+)\)' # Regex pattern used for selecting and capturing specific content.
     # The captured content is then utilized in the replacement pattern below.
     with: "[$1](/$2)" # This is the replacement pattern for converting Markdown links.
whitelist: ['https://spandigital.net/assets/'] # URLs that should be whitelisted for conversion
assetDir: 'assets' # The directory where assets should be saved
contentDir: 'content' # The directory where the converted markdown files should be saved


RST2MD
RST2MD is a tool that allows you to convert RST files into Presidium MD. It’s a custom wrapper around Pandoc, that slightly modifies the output to be more compatible with Presidium. (Link to RST2MD binary still needs to be added)

Installation

Installation is dependent on how the tool will be hosted.

Usage

The tool accepts an input(RST source files) and output parameter(location where generated MD files should be saved).

./rst2md -v -input /path/to/rst/source/ -output /path/to/output/dir

JSON Schema
presidium-json-schema is a CLI tool for importing your JSON Schema spec into Presidium documentation.
Install
brew tap SPANDigital/homebrew-tap https://github.com/SPANDigital/homebrew-tap.git brew install presidium-json-schema
Usage
Usage:
 presidium-json-schema convert [path] [flags]

Flags:
 -d, --destination string   the output directory (default ".")
 -e, --extension string     the schema extension (default "*.schema.json")
 -o, --ordered              preserve the schema order (defaults to alphabetical)
 -w, --walk                 walk through sub-directories
To convert a file you simply:
presidium-json-schema convert <PATH_TO_SCHEMA_DIR> -d <THE_DESTINATION_DIR>


Embed
A fallback approach to importing generated documentation is to embed documentation in an iframe. This approach is not recommended because items are not indexed or available on the main menu. However, it will work for certain cases when an importer is not yet available.
When possible, use a simple template when embedding documentation in an iframe.
To include documentation in an iframe:
Generate the static site documentation for your component.
Put the documentation in the /static folder so that it’s statically served. The Presidium convention is to place it under /static/import/{my-reference}.
Add a reference article to the Reference section:

---
title: My Reference
---

# foo.bar

<div>
   <iframe>
           src='/static/import/{my-reference}/foo/bar/package-summary.html'
   </iframe>
</div>
You can create multiple Markdown files for different components as required.

