# Presidium

Presiduim is a Jekyll based documentation framework to write technical documentation that may be hosted on Github pages.

## Prerequisites
`Currently only OSX is supported`

The following is required to build and run Presidium locally:

- OS: Mac OSX 
    - [npm](http://www.npmjs.com) > 4.1.2
    - [ruby](https://www.ruby-lang.org/en/documentation/installation/) > 2.4
    - [bundler](http://bundler.io/) > 1.14.3

For further details on setting up Jekyll for Github pages, see: https://help.github.com/articles/setting-up-your-github-pages-site-locally-with-jekyll/

### NPM

For [homebrew](http://brew.sh/) users, to upgrade node:

```sh
$ brew install node
$ brew update
$ brew upgrade node
```

### Ruby

Ruby 2.0 is already included on OS X (El Capitan, Yosemite, Mavericks and Sierra)
Newer versions are available via Homebrew if required: [rbenv](https://github.com/rbenv/rbenv#readme), or [RVM](http://rvm.io/).

### Bundler
Bundler is recommended to install Jekyll:
```sh
$ gem install bundler
```

## Getting Started

Fork the reference repo: https://github.com/SPANDigital/presidium.git

On the target environment:
```
$ git clone {forked repo}
$ npm install
$ npm start
```

This will install all required dependencies and fire up a Jekyll server on your machine on http://localhost:4000/. 
Changes to content and styling will be automatically watched. Changes to config will only be applied on startup

### Directory Structure

```
presidium/

    content/
        [section folders]/ : Folders with markdown content

    media/ : 
        css/ : Sass Styling
        images/ : Image Resources
        [other]/ : Static Files
    
    dist/ : Distribution Directory with Jekyll Sources and Site
        
    _config.yml : Jekyll configuration
    
    package.json : npm build, serve and package scripts        
```

### Content
All content updates are done using the template structure in the `content` directory.

### Configuration
Base site configuration such as copy and menu structure can be done in `_config.yml`

### Styling and Theming
All styles for Presidium are based on Bootstrap with Bootswatch overrides using sass.

To change the theme or to provide your own styling you can provide overrides using:
 - Update `_custom.scss` : set any custom css
 - Update `_variables.scss` : set a different theme and any custom bootstrap variables

## Publish to Github Pages

To publish to Github pages, you must configure your repository to use a gh-pages branch. Further details may be found on [Github](https://help.github.com/articles/about-github-pages-and-jekyll/)

Once you are ready to publish your site, run:
```sh
$ npm run gh-pages
```