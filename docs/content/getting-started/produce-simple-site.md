---
title: Produce a Simple Documentation Site
weight: 5
---

#### Prerequisites
To produce a simple documentation site with Presidium, ensure that you have a functional development
environment: 
* Git, for version control
* Hugo version 0.131.0 or greater, for site generation
* Go, for support of modules

Also helpful:
* Familiarity with Markdown and YAML
* Access to a terminal and a text editor

#### Install Presidium

1. Install Homebrew if you don't have it already. For macOS:
      ```bash
      /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
      ```

2. Tap SPAN Digital's Homebrew Repository:
      ```bash
      brew tap SPANDigital/homebrew-tap https://github.com/SPANDigital/homebrew-tap.git
      ```

3. Install Presidium and Hugo:
      ```bash
      brew install presidium
      brew install hugo
      ```
4. Optionally, [install Golang](https://go.dev/doc/install). This is required only if you are building from source.

#### Run Presidium
1. Initialize the project:
   1. Navigate to your chosen directory and initialize your Presidium project:
      ```bash
      presidium init
      ```
   2. Go through the setup wizard:
        - Define the project name and title.
        - Select appropriate Template and Theme according to your project’s purpose.
2. Start a local development server:
   1. After initialization, navigate into your project directory:
      ```bash
      cd [your-project-name]
      ```
   2. Serve your site locally:
      ```bash
      hugo serve
      ```
3. Navigate to `http://localhost:1313` to see the Hello Presidium article.

#### Add Content
1. Update the `config.yml` file to structure your site's top-level sections. For example:
      ```yaml
      - identifier: overview
        name: Overview
        URL: /
        weight: 10
      - identifier: getting-started
        name: Getting Started
        URL: /getting-started/
        weight: 20
      ```

    For the use of the `weight` variable, see [Ordering Items]({{<ref "/presidium-basics/structuring-content/#ordering-items">}}).

2. Create content Using Markdown:
    - Add articles and directories in the `content` directory. Use clear filenames and directories that reflect the content structure:
      ```bash
      hugo new content/overview/introduction.md
      ```
    - Edit articles using your text editor, maintaining consistency in keeping with project guidelines.
    - Add any desired image files to the directory containing the article they are used in.

At any point, you can check the appearance of articles by running `hugo serve` and going to `http://localhost:1313`.

#### Run Presidium Locally
If you are running your site locally, Presidium generally updates whenever you make changes.
Some more complex changes may require you to quit (control-C), then relaunch with `hugo serve`.
