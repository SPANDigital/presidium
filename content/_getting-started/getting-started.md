---
title: Setting up your First Site
---

# Create a New Site
1. Download and extract the [sample template](https://github.com/SPANDigital/presidium-template)
2. Run the following from your project directory to install required dependencies:

```sh
$ npm install
```

# Use an Existing Presidium Site
1. Clone the git repository.
2. Run the following from your project directory to install required dependencies:

```sh
$ npm install
```

If you get an error, check that you have the following [prerequisites](/prerequisites/) installed on your environment:
- npm v3.10+
- ruby v2.1+
- bundler v1.14.3+

# Run Presidium
To start your site, run the following from your project folder:
```sh
$ npm start
```

This will build your site to `dist/site` and serve it locally on: [http://localhost:4000/](http://localhost:4000/)

# Edit Content

Once your site is up and running, you can start making changes to the following files and folders:

- `content/` Content folder for all your articles
- `media/` Media folder for static assets such as images, attachments or css
- `_config.yml` Site configuration properties

Updates to content, media or css can be made without having to restart the local server. 
Structural or changes to configuration currently require a server restart.

# Publish your Site
The simplest way to publish a github repository is to use Github Pages, but you can also host the generated site on any Web server.
## To Github Pages
To publish using Github Pages, commit and push your site to a Github repository and run the following:
```sh
$ npm run gh-pages
```
This will push your generated site to a gh-pages branch in your repository. You will need to 
[enable gh-pages](https://help.github.com/articles/configuring-a-publishing-source-for-github-pages/) 
in your repository.

## As a Static Site
The generated static site can be found in `dist/site`