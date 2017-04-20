# Presidium

Presiduim is a Jekyll based documentation framework that enforces industry best practices for documenting software systems. 
Presidium provides software teams with familiar tools to acquire, revise, categorize, utilize and evaluate small 
document fragments that get aggregated into well structured website.

## Quick Start

- Download the [template project](https://github.com/SPANDigital/presidium)
- Requires `ruby`, `bundler` and `npm`. See [setting up your environment](#setup) for further details and requirements.
- To install the required dependencies and start the documentation server locally, run:
```sh
$ npm install
$ npm start
```

This will install all required dependencies and fire up a Jekyll server on: http://localhost:4000/

Add and edit your content using examples in `content` and `media`
Configure your site using options in `_config.yml`

## Setup and Configuration

### Directory Structure
Presidium uses the following directories:
```
presidium/

    content/
        [section folders]/ : Folders with markdown content

    media/ : 
        css/ : Sass styling
        images/ : Image resources
        [other]/ : Static files
    
    dist/ : Distribution directory 
        src/: Presidium and Jekyll sources
        site/ : Generated static site
    
    .jekyll/ : Jekyll dependencies
        
    node_modules/ : node module dependencies, including the Presidium templates and structure
        
    _config.yml : Jekyll configuration file
    
    package.json : npm build scripts
```

### Content
All content updates are done using the sections and markdown templates in `content` using markdown. See the [template project](https://github.com/SPANDigital/presidium) for example usage.
* [Jekyll](http://jekyllrb.com/) a simple, blog-aware, static site generator that integrates with GitHub
* [GitHub Flavored Markdown](https://help.github.com/articles/github-flavored-markdown/) GitHub Flavored Markdown to document your content.

Static sources such as images, files and custom styles should be added to the `media`.

Content and media resources are kept separate from the underlying Jekyll layouts and templates.

### Build Scripts
The following build scripts manage the main workflows:
- `npm install` : Installs all npm and Jekyll dependencies required to build, run and publish your site.
- `npm start` : Serves the Jekyll site from `dist/site` and Watches for any changes to `content` and `media`.
- `npm run build` : Build your Jekyll site to `dist/site`.
- `npm run gh-pages` : Publishes `dist/site` to a `gh-pages` branch in your current repo.

### Configuration
Site configuration such as the site name and menu structure can be done in `_config.yml`. Changes made to config require a site rebuild to take effect.

### Styling and Theming
All styles are based on Bootstrap with Bootswatch themes and sass overrides. Jekyll themes are not supported.

* [Sass](http://sass-lang.com/) a CSS extension language
* [Bootswatch](https://bootswatch.com/) Boostrap-based themes

To change a theme or to provide your own styling or overrides, make changes to:
 - `_custom.scss` : any custom styles that you would like to apply
 - `_variables.scss` : change the theme and override and bootstrap or bootswatch variables
 
### Publishing to Github Pages
To publish your documentation to Github pages, run:
```sh
$ npm run gh-pages
```

This pushes your site and to a gh-pages branch on your repo. To enable your site on Github, go to your repository settings and enable Github pages using your `gh-pages` branch. 

Further details may be found on [Github](https://help.github.com/articles/about-github-pages-and-jekyll/)

## <a name="setup"></a>Setting up your Environment
`Instructions are currently only available for OSX`

The following tools are required to build and run Presidium locally:

- OS: Mac OSX 
    - [npm](http://www.npmjs.com) > 3.10 : Node package manager to build, serve and publish
    - [ruby](https://www.ruby-lang.org/en/documentation/installation/) > 2.1 : To build Jekyll site
    - [bundler](http://bundler.io/) > 1.14.3 : To manage Jekyll dependencies

### NPM

- node v6.10.1
- npm v3.10.10

For [homebrew](http://brew.sh/) users, install and upgrade node:

```sh
$ brew install node
$ brew update
$ brew upgrade node
```

### Ruby

Requires ruby >= 2.1

Newer versions are available via Homebrew:
```sh
$ brew install ruby
```

### Bundler
Bundler is recommended to install Jekyll and its dependencies locally:
```sh
$ gem install bundler
```

For further details on setting up Jekyll for Github pages, see: [Setting up github pages locally](https://help.github.com/articles/setting-up-your-github-pages-site-locally-with-jekyll/)

## Build Server

### Provisioning a Debian Based build server:

```sh
#!/bin/bash

sudo apt-get update
sudo apt-get install -y curl git

# Install Ruby
# http://rvm.io/rvm/install
gpg --keyserver hkp://keys.gnupg.net --recv-keys 409B6B1796C275462A1703113804BB82D39DC0E3
\curl -sSL https://get.rvm.io | bash -s stable --ruby=2.4.0
source ~/.rvm/scripts/rvm
rvm install 2.4.0 --quiet-curl
rvm use 2.4.0 --default
ruby --version

# Install bundler
gem install bundler

# Install NPM
\curl -sL https://deb.nodesource.com/setup_6.x | sudo -E bash -
sudo apt-get install -y nodejs
```