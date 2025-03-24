---
title: HTML to Markdown
weight: 30
---

html2md is a tool that allows you to convert HTML files into Presidium markdown articles.

### Installation
Install from Homebrew:
```
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh)"
```
Add SPAN’s Homebrew tap:
```
brew tap SPANDigital/homebrew-tap https://github.com/SPANDigital/homebrew-tap.git
```
Install html2md:
```
brew install html2md
```
### Usage

`html2md convert [source] [dest] [flags]`


Params:
* `source`: The URL of the website or the path to the local HTML file(s).
* `dest`: The path to the directory where the converted markdown files are saved.
* `flags`:
    * ` -d, --debug`: Enable debug logging
    * `--headers <strings>`: Article header tags (default is `h1,h2`])
    * `--select <string>:` The part of the page to select and convert (default is `body`)


### Examples
Download and convert the [Presidium](https://presidium.spandigital.net/) website
```
html2md convert https://presidium.spandigital.net/ ./presidium --select="#presidium-content"
```
Convert local html files
```
html2md convert ./html-files ./presidium --select=.article --headers=h1
```

> **Note**
>
> This converter does not output all the files needed to build a complete Presidium website.
> You will still need to [create]({{<ref "getting-started">}}) a Presidium site first, or import the converted markdown files into an existing site.
> The converted markdown files will be saved in the `content` directory relative to the `dest` path along with assets (images, videos, etc) in the `assets` directory.

### Advanced usage
You can define a `config.yml` file in your working directory or in `$HOME/.html2md` directory with additional options to control the conversion.

#### Remove HTML
The `html.remove` option allows you to selectively remove elements from your source document before converting it to Markdown. This can be useful when you want to clean up your HTML content or remove unwanted elements that have no relevance.

**Example Usage:**
```
html:
 remove: ['.nav-link', '#warning']
```
In this example, all element matching the CSS selectors .nav-link or #warning will be removed before conversion.

#### Replace HTML
The `html.replace` option allows you to transform specific HTML elements into custom Markdown syntax, giving you greater control over the appearance and structure of your converted document.

**Example Usage:**

Suppose you have HTML content that includes [tooltips]({{<ref "reference/markdown/tooltips.md">}}), and you want to convert them into Markdown-friendly syntax. Here’s how you can achieve this using the “Replace HTML” feature:

```
html:
 replace:
 - match: '.tooltips-term' # CSS selector for the element to be replaced
   select: ['?href', '.tooltips-text'] # ?href selects the href attribute of the matched element, .tooltips-text selects the content of a child element with the class "tooltips-text"
   replace: '{{</* tooltip "$1" text="$2" */>}}' # $1 and $2 are the selected elements
```

In this example, the CSS selector `.tooltips-term` is specified as the element to be replaced. The `select` field allows you to capture specific attributes and content relative to the matched element. Finally, the `replace` pattern converts the selected elements into a Markdown tooltip format.

>**Note**
>
>In addition to the standard CSS-selectors, select allows you to select attributes on the matched element using the ? prefix. You can also get the text of the matched element by passing text in the select list.

#### Replace Markdown
Similar to the HTML replace you can also replace markdown based on a Regex pattern.

**Example Usage:**

Let’s say you have Markdown content with links, and you want to prepend the base URL to all these links.
```
markdown:
 replace:
   - pattern: '\[([^]]+)\]\(([^\)]+)\)'
     with: "[$1]({{< baseurl >}}$2)"
```

#### Sample config
```
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
```