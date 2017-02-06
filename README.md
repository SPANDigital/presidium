# Presidium

Presidium is a technical documentation template based on Jekyll

## Prerequisites

- OS: Mac OSX
    - [Homebrew](http://brew.sh/)
    - [npm](http://www.npmjs.com)
 
### Install npm

```sh
$ brew install node
$ brew update
$ brew upgrade node
```
 
### Install Ruby

On OS X El Capitan, Yosemite, Mavericks, and macOS Sierra, Ruby 2.0 is included. OS X Mountain Lion, Lion, and Snow Leopard ship with Ruby 1.8.7. Newer versions are available via Homebrew, [rbenv](https://github.com/rbenv/rbenv#readme), or [RVM](http://rvm.io/).

### Install Jekyll

```
$ gem install jekyll
```

## Initialize

```sh
$ npm install
```

## Run

Run a local jekyll server and watches for changes http://localhost:4000/
```sh
$ npm run serve
```

## Publish
To publish to gh-pages using the docs directory:
```sh
$ npm run gh-pages
```

## Directory Structure

```
presidium/

    content/
        [section folders]/ : Content folders for sections

    media/ : 
        css/ : Sass styling overrides
        images/ : Image resources
        [other]/ : Static files
    
    dist/ : TODO review: Distribution directory with packaged Jekyll site
    docs/ : TODO review: gh-pages
    
    .build/ : Presidium dependencies
        presidium-core/ : submodule dependency to presidium core jekyll template
        presidium-js/ : submodule dependency to presidium javascript components
    
    _config.yml : Jekyll configuration
    
    package.json : npm build, package and serve scripts        
```