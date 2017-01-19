# presidium-reference

A reference documentation site using presidium.

## Installation Requirements

- OS: Mac OSX
- [Homebrew](http://brew.sh/)
 
### Install latest npm

```sh
$ bew install node
$ brew update
$ brew upgrade node
```
 
### Install Ruby

On OS X El Capitan, Yosemite, Mavericks, and macOS Sierra, Ruby 2.0 is included. OS X Mountain Lion, Lion, and Snow Leopard ship with Ruby 1.8.7. Newer versions are available via Homebrew, [rbenv](https://github.com/rbenv/rbenv#readme), or [RVM](http://rvm.io/).

### Install Jekyll

```
$ gem install jekyll
```

## Bells & Whistles

- Install presidium as an npm module.
- Allow users to specify version in packge.json.
- Remove underscores from content directory, and prepend them during the build/copy.
- Let npm handle installation and set up of Jekyll gem.
- Content becomes a reference repo.
- Users allowed to switch sections off in sections.yml. sections.yml is versioned by presidium.
- Users allowed to change _config.yml. _config.yml versioned by presidium.
- npm init
    - Pull content reference repo.
    - Pull npm presidium module.
    - Copy latest config from module to root, for user override.
    - Copy latest sections from module to root, for user override.
- npm build
    - Clean and copy files into dist.
    - Serve/compile using jekyll.
