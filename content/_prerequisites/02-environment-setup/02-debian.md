---
title: Debian
---
The following script can be used to provision a Debian based environment:

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
