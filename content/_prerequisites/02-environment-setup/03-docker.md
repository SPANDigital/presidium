---
title: "Docker"
---

Using VS Code's extension [Remote - Containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) to create Presidium content is a simple and consistent way to set up a development environment without needing to update any configuration on your local machine.

## Required Software
**macOS:** 
- [Docker Desktop](https://www.docker.com/products/docker-desktop) 2.0+
- [VS Code](https://code.visualstudio.com/) or [VS Code Insiders](https://code.visualstudio.com/insiders/)

**Linux:** 
- [Docker CE/EE](https://docs.docker.com/install/#supported-platforms) 18.06+ 
- [Docker Compose](https://docs.docker.com/compose/install) 1.21+. (The Ubuntu snap package is not supported.)
- [VS Code](https://code.visualstudio.com/) or [VS Code Insiders](https://code.visualstudio.com/insiders/)

**Windows:** 
- [Docker Desktop](https://www.docker.com/products/docker-desktop) 2.2+
- [WSL2 back-end](https://aka.ms/vscode-remote/containers/docker-wsl2). (Docker Toolbox is not supported.)
- [VS Code](https://code.visualstudio.com/) or [VS Code Insiders](https://code.visualstudio.com/insiders/)

## Installation

1. Install any required software (above).

2. Ensure that git user config is configured on your **local** environment: 
    ```bash
    git config --global user.name "Your Name"
    git config --global user.email "your.email@address"
    ```
    **Note:** If a Git Credential Manager (e.g. Github Desktop) is configured on the local environment, this step can be skipped.


4. Inside VS Code install the extension: [Remote - Containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers).

## Setting up your project
1. Copy the link to the Github repository to be cloned.

1. In [VS Code](https://code.visualstudio.com/) open the command pallet (F1 or &#8984;&#8679;`P`) and type 
`>Remote-Containers: Clone Repository in Container Volume...`
    
    **Note**: This step replaces the first step of the [Getting Started](/getting-started/#download-the-template) instructions.

1. Paste the URL for the repository which should be cloned.

1. A popup asking which branch to clone may appear. Select the one you want to checkout initally, once the container is created you are free to checkout any branch in the repository.

1. When prompted for a [Volume](https://docs.docker.com/storage/volumes/) to use there are a few options:
    - (Recommended) Create a single volume for all project. This is the most space efficient, and unless there is a reason two project should not share a filesystem, there is no harm to this for most development environments.
    - Create a Persistent volume for every/related projects. Share a filesystem when needed and isolate them otherwise.
    - Create a unique volume per project. This is the least space efficient and means that no two projects can (easily) access each others files.
    
    The name you choose for your volume is entirely up to you and should have no impact on the project itself.

1. A prompt asking for a project name will pop-up, the pre-generated name is pulled from the github repository name, and is usually a good choice.

1. The container will start building, this could take a few minutes the first time as it is downloading the required images.

1. A popup asking to "Add Development Configuration Files" may appear, select Jekyll (Community). You may need to click on the "Show All Definitions..." to find it. This popup will only appear if the repository does not already contain a `'.devcontainer'` folder containing the environment configuration details.

1. The development container will continue building, the very bottom left corner of VS Code should say: `Dev Container: Jekyll (Community)` once the environment is ready.
