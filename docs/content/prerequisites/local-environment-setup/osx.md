---
title: OSX
weight: "1"
---
# Xcode

Install Xcode command line tools
```sh
$ xcode-select --install
```


# NPM

Install using the [node package](https://nodejs.org/en/) or [nvm](http://nvm.sh)
- node v6.10 (LTS)
- npm v3.10


# Ruby

Requires ruby >= 2.1 available using [homebrew](https://brew.sh):
```sh
$ brew install ruby
```

# Bundler
```sh
$ gem install bundler
```

# M1 Mac IMPORTANT ADDITIONAL STEPS
   
If you are using a Mac (2020 or newer) with the M1 chipset, Ruby will need to be version 2.7.3 (recommend using RVM to install this). Also, there may be errors encountered while building the ffi library. The following commands can be run to rebuild it manually for the M1 architecture:
```sh 
git clone https://github.com/libffi/libffi
cd libffi/
./autogen.sh
./configure
make
make install
Delete libffi directory if you dont need it
```