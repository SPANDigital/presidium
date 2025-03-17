---
title: Produce a Simple Documentation Site
weight: 5
---

#### Prerequisites
To embark on producing a simple documentation site with Presidium, ensure that you have a functional development
environment. This includes having Git installed for version control, as well as Hugo, which Presidium relies on for
site generation. Familiarity with Markdown is advantageous as it is the primary language for content creation.
Additionally, access to a terminal and a text editor will facilitate the development process.

#### Install Presidium

1. **Homebrew Installation (for macOS):**
    - Install Homebrew if not already installed:
      ```bash
      /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
      ```

2. **Tap SPAN Digital's Homebrew Repository:**
    - Use the following command to access Presidium:
      ```bash
      brew tap SPANDigital/homebrew-tap https://github.com/SPANDigital/homebrew-tap.git
      ```

3. **Install Presidium and Hugo:**
    - Execute these commands to install the necessary components:
      ```bash
      brew install presidium
      brew install hugo
      ```
    - Golang installation is also suggested but not mandatory unless building from source.

#### Run Presidium
1. **Initialize the Project:**
    - First, navigate to your chosen directory and initialize your Presidium project:
      ```bash
      presidium init
      ```
    - Go through the setup wizard:
        - Define the Project Name and Title.
        - Select appropriate Template and Theme according to your project’s purpose.
2. **Start Local Development Server:**
    - After initialization, navigate into your project directory:
      ```bash
      cd [your-project-name]
      ```
    - Serve your site locally to preview changes using:
      ```bash
      hugo serve
      ```
    - This launches the site at `http://localhost:1313`.

#### Add Content
1. **Edit the `config.yml` file:**
    - Update the `config.yml` to structure your site's top-level sections.
    - Example:
      ```yaml
      - name: Overview
        URL: /
      - name: Getting Started
        URL: /getting-started/
      ```

2. **Create Content Using Markdown:**
    - Add articles in the `content` directory. Use clear filenames and directories that reflect the content structure:
      ```bash
      hugo new content/overview/introduction.md
      ```
    - Edit articles using your text editor, maintaining consistency in format with project guidelines.

#### Run Presidium Locally
- Re-run `hugo serve` whenever you make changes. This refreshes your local server, allowing you to view updates in
- real time. It ensures documentation is accurate before publishing.
