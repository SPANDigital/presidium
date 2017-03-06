# Presidium

Presiduim is a Jekyll based documentation template for writing technical documentation that may be published to Github pages or as a static site.
Jekyll is used to build a static site from markdown using a set of standard templates and layouts.

## Quick Start

- Requires `ruby`, `bundler` and `npm` to run locally. See the [prerequisites](#prerequisites) for further details and requirements.
- To install the required dependencies and start the documentation server locally, run:
```sh
$ npm install
$ npm start
```

This will install all required dependencies and fire up a Jekyll server on: http://localhost:4000/`baseurl`
Add and edit your content in `content` and `media`

`baseurl` is an optional path configuration set in `_config.yml`

## Setup and Configuration

### Directory Structure

```
presidium/

    content/
        [section folders]/ : Folders with markdown content

    media/ : 
        css/ : Sass styling
        images/ : Image resources
        [other]/ : Static files
    
    dist/ : Distribution directory with Jekyll sources and site
        _site/ : Generated static site
    
    .jekyll/ : Dependencies for Jekyll and Github pages
        
    node_modules/ : node module dependencies, including the presidium templates and structure
        
    _config.yml : Base Jekyll configuration file
    
    package.json : npm build scripts
```

### Content
Content and media resources are kept separate from the underlying Jekyll templates and layouts.
All content updates are done using the sections and markdown templates in `content`.
Static sources such as images, files and custom styles should be added to the `media`.

### Build Scripts
The following build scripts manage the main workflows:
- `npm install` : Installs all npm and Jekyll dependencies required to build, run and publish your site.
- `npm start` : Serves the Jekyll site from `dist/` and Watches for any changes to `content` and `media`.
- `npm run build` : Build your Jekyll site to `dist/`.
- `npm run gh-pages` : Publishes your site `dist/` to a `gh-pages` branch in your current repo.

### Configuration
Site configuration such as the site name and menu structure can be done in `_config.yml`. Changes made to config require a site rebuild to take effect.

### Styling and Theming
All styles are based on Bootstrap with Bootswatch themes and sass overrides.

To change a theme or to provide your own styling, make changes to:
 - `_custom.scss` : any custom styles that you would like to apply
 - `_variables.scss` : change the theme and override and bootstrap or bootswatch variables
 
### Publishing to Github Pages

To publish your documentation to Github pages, run:
```sh
$ npm run gh-pages
```

This pushes your site and to a gh-pages branch on your repo. To enable your site on Github, go to your repository settings and enable Github pages using your `gh-pages` branch. 

Further details may be found on [Github](https://help.github.com/articles/about-github-pages-and-jekyll/)

## <a name="prerequisites"></a>Prerequisites 
`Only OSX is currently supported`

The following tools are required to build and run Presidium locally:

- OS: Mac OSX 
    - [npm](http://www.npmjs.com) > 3.10
    - [ruby](https://www.ruby-lang.org/en/documentation/installation/) > 2.1
    - [bundler](http://bundler.io/) > 1.14.3

### NPM

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

### Debian Based Provisioning:

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