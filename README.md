# presidium-hugo

## Quickstart

###$ installation of hugo

### Prerequisites
Install [homebrew](https://brew.sh/) if you don't have it already.

### Issue appropriate command

If you are installing hugo from scratch:

```brew install hugo```

If you have a previous version of hugo.:

```brew upgrade hugo```

Confirming your version:

```hugo version```

Your version needs to be at or greater then this version

```Hugo Static Site Generator v0.75.0```

### Tell git to use your ssh credentials for SPANDigital's private repos

```
git config --global url.git@github.com:.insteadOf https://github.com/
export GOPRIVATE=github.com/SPANDigital
```

### Install / upgrade binary

To install first time
```
go install github.com/SPANDigital/presidium-hugo
```
To upgrade to latest version
```
go get -u github.com/SPANDigital/presidium-hugo
```

### Converting a Jekyll content directory

```
mkdir ~/Documents/mycontentdirectory
cd ~/Documents/mycontentdirectory
export PATH=$PATH:$HOME/go/bin
presidium-hugo convert -s ~/spandigital/span-handbook-docs -p name-of-initial-section
```

### Converting a Jekyll content directory with a locally changed presidium-hugo

```
go build main.go && ./main  convert -s ~/Documents/projects/site-to-convert/  --destDir  ~/Documents/mycontentdirectory -p name-of-initial-section
```

### Initializing Hugo's module system

````
cd ~/mycontendirectory
hugo mod init github.com/spandigital/mycontendirectory
hugo mod get -u
````
### Serving (non-enterprise)
````
cd ~/mycontendirectory
hugo serve
````

### Vendoring
````
cd ~/mycontendirectory
hugo mod vendor
````

### Zipping for enterprise (requires vendoring)
````
zip -r documentation.zip content static config.yml _vendor
````

### Usage

```CLI tools for managing Presidium Hugo content

   Usage:
     presidium-hugo [command]

   Available Commands:
     convert     Convert Jekyll to Hugo content
     help        Help about any command

   Flags:
         --config string   config file (default is $HOME/.presidium-hugo.yaml)
     -h, --help            help for presidium-hugo

   Use "presidium-hugo [command] --help" for more information about a command.```
```

#### Convert usage

```Convert Jekyll to Hugo content

   Usage:
     presidium-hugo convert [flags]

   Flags:
     -m, --commonmarkAttributes         Convert to commonmark attribute format
     -y, --convertConfigYml             Convert jekyll _config.yml to hugo config.yml (default true)
     -C, --copyMediaToStatic            Copy Jekyll media to Hugo static folder (default true)
     -d, --destDir string               Destination directory (default "/Users/richardwooding/spandigital/presidium-hugo")
     -c, --enableColor                  Enable colorful output (default true)
     -e, --eraseMarkdownWithNoContent   Erase markdown files with no content (default true)
     -i, --fixImages                    Fix images in same path (default true)
     -a, --fixImagesWithAttributes      Replace images with attributes with shortcodes (default true)
     -h, --help                         help for convert
     -R, --removeRawTags                Remove {% raw %} tags (default true)
     -t, --removeTargetBlank            Remove target="blank" variants (default true)
     -b, --replaceBaseUrl               Replace {{site.baseurl}} with {{ site.BaseURL }} (default true)
     -j -- replaceBaseUrlWithSpaces     Replace {{ site.BaseURL }} with {{ site.baseurl }} (default true)
     -o, --replaceCallOuts              Replace callout HTML with callout shortcodes (default true)
     -V, --replaceIfVariables           Replace {% if site.variable =} with with-param shortcodes (default true)
     -p, --replaceRoot string           Replace this path with root
     -T, --replaceTooltips              Replace tooltip HTML with callout shortcodes (default true)
     -g, --slugBasedOnFileName          Base front matter slug on filename (default true)
     -s, --sourceRepoDir string         Source directory
     -u, --urlBasedOnFilename           Base front matter url on filename (default true)
     -w, --weightBasedOnFilename        Base front matter weight on filename (default true)

   Global Flags:
         --config string   config file (default is $HOME/.presidium-hugo.yaml)
  ```

### How to make changes to the templates (styling and html)

After vendoring the site
```
hugo mod vendor
```

The converted site will contain a ```_vendor``` folder

To make changes locally to ***styling*** locate
```
_vendor/github.com/spandigital/presidium-theme-website/assets/presidium.scss
```
Or
```
_vendor/github.com/spandigital/presidium-theme-website/assets/sass
```
To make changes locally to ***html*** locate
```
_vendor/github.com/spandigital/presidium-theme-website/layouts
```

Changes should reflect instantly if
```
hugo serve
```
is running.

To commit any styling or html changes, for now, you will need to commit to three repos (in the same folder you made the locall changes)
1. [presidium-theme-website](https://github.com/spandigital/presidium-theme-website)
2. [presidium-theme-pdf](https://github.com/SPANDigital/presidium-theme-pdf)
3. [presidium-hugo](https://github.com/SPANDigital/presidium-hugo/tree/master/themes/presidium) in the themes folder  ```themes/presidium```