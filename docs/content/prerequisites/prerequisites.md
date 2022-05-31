---
title: Prerequisites

---

The following tools are required to use Presidium:

1. Homebrew installed

    ```
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
    ```

1. Tell git to use your ssh credentials for SPANDigital's private repos

    ```
    git config --global url.git@github.com:.insteadOf https://github.com/
    export GOPRIVATE=github.com/SPANDigital
    ```

1. Github personal access token for private repos

   To access private repos you will need to generate a [personal acess token](https://github.com/settings/tokens) with *repo* access. 
   Then set this as an environment variable:
   
   ```
   export HOMEBREW_GITHUB_API_TOKEN=xxxxxx
   ```
