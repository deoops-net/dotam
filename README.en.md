# Dotam

Dotam pronounced as dot-am, means dotfiles automation。

<!-- TOC -->

- [Dotam](#dotam)
    - [Why Another Automation Tool](#why-another-automation-tool)
    - [Features](#features)
    - [Quick Start](#quick-start)
    - [Installation](#installation)
        - [**homebrew**](#homebrew)
        - [**Pre-build binaries**](#pre-build-binaries)
        - [**goland**](#goland)
        - [**Docker**](#docker)
        - [**debugging or install from source**](#debugging-or-install-from-source)
    - [Usage](#usage)
        - [run](#run)
        - [initial a template](#initial-a-template)
    - [References](#references)
        - [Template module](#template-module)
        - [Docker module](#docker-module)
        - [Plugin](#plugin)
        - [Tempalte Syntax](#tempalte-syntax)
    - [FAQ](#faq)
        - [Get Support](#get-support)

<!-- /TOC -->

## Why Another Automation Tool

As myself are not only a developer but also a devops. I use drone for CI/CD automation, but for 
some senses drone may not meet my needs, like: when a finish some code, I should change the release version
for many files(deploy.yml for kubernetes, .drone.yml for CI/CD, Makefile, and some other script files)
If I forgot to change one i need to re-commit them.

So, why there isn't a tool can do something like `ansible` which can help me maintain these files
using template file. Both it can also run some `pipeline` tasks like drone. Because sometimes drone 
runs task slowly in a remote environment and to use it's cache feature is hard(we need config git plugin to cache it, but
at the same time we also need to implement a system lock for resource racing).

I did some research, but no one meets this, so I try to build this one.

## Features

* language independent
* simple, low study curve
* support most popular DSL(json, yml, hcl)
* built-in plugins(template, docker, git)
* build custom plugins easily 

PS: optimizing performance is a long term job, so as early version it will not be a high priority
if you are good at this, pls PR! 

## Quick Start

Follow the [Installation](#Installation) section to install it, it's quit easy.
Then what you need is a `Dotamfile.{json,yml,yaml,hcl}`  tells `dotam` what to do,
below is a demo, normally you can generate it with `dotam init`

Dotamfile.hcl:

```hcl
temp "RELEASE" {
    src = ".dotam/RELEASE"
    dest = "."
    var {
        version = "{{versions.release}}"
    }
}

docker {
    repo = "deoops/dotam"
    tag = "{{versions.release}}"
    
    auth {
        username = "tom"
        password = "_args.reg_pass"
    }
}

var "versions" {
    prod = "v0.1.1"
    release = "v0.1.3-beta"
}

```

After editing the config file, you can run `dotam build`. That's it!

## Installation

### **homebrew**

```bash
brew tap deoops-net/tap

brew install dotam
```

### **Pre-build binaries**

Go to [Release](https://github.com/deoops-net/dotam/releases) page to download a prebuild version.

### **goland**

`go install github.com/deoops-net/dotam`

### **Docker**

comming soon!

### **debugging or install from source**

```bash
git clone https://github.com/deoops-net/dotam
cd dotam
make test
make install
```

## Usage

### run

`dotam build [Dotamfile.{yml,hcl,json}]`

by default, `dotam` will try to find a `Dotamfile` with json, hcl, yml|yaml file extension.
but you can also specify one to let him run it: `dotam build Dotamfile.json`

### initial a template

`dotam init [-t yml|yaml,json,hcl]`

this will generate a demo file into current folder, you can use `-t` to specify a file extension


## References

### Template module

Normally we need to maintain multi static files, but usually the files needs some dynamic properties
such as:

* we need to specify a Docker tag for CI files.
* we need to specify some version or tags for deployment files of kubernetes.
* we may also write something into a `RELEASE` file to let our program read and log it.

Change these manually is boring, so we need a method that can define once then use any where.

Template module is quite simple, here is a example:

```hcl
temp "Makefile" {
  # specify the source file for rendering
  src = ".dotam/Makefile"
  # the destination dir
  dest = "."
  # specify the variables used for rendering
  var {
    # such as we want to trans a production version for Makefle
    version = "{{versions.production}}"
  }

}
// defines the variables for rendering
var {
  versions {
    prodution = "v1.0.0"
  }
}
```

That in a Makefile we can use the variable wrote above:

```Makefile
BUILD_VERSION={{version}}

.PHONEY: build
build:
    npm run build ${BUILD_VERSION}

```

### Docker module

Docker module can help us build and push images to a specify registry.
here is a demo:

```hcl
docker {
  repo = "deoops/dotam"
  tag = "_args.version"
  auth {
    username = "_args.reg_user"
    password = "_args.reg_pass"
  }
  caporal {
    host = "your deployed caporal host"
    name = "container name your want to start or update"
    opts {
      // the -p flag of docker run 
      publish = ["8080:8080"]
      // the network flag of docker run 
      network = ""
    }
  }
}
```

As you can see we use particular variable `_args`, this will retrieve data from the command line
parameters. As the example above to can type the command like this: `dotam build version=v1.0.12 reg_user=tom reg_pass=pass`

That is quite simple, why you should use the Docker module not use a build tool like Makefile?
In a further version, i will integrate a remote api-server for scheduling containers with this
module.For teams already built there k8s cluster, there is no need, but for some small teams
a strong and simple docker scheduler is necessary.

[caporal](https://github.com/deoops-net/caporal) is an early version of container scheduling api server, basiclly it's now only support simple jobs
like run or update a container remotely.but it will be under development for more features.

### Plugin

**Not Implemented**

### Tempalte Syntax

For some static variables you can just use {{variable}} mentioned above.
For some advanced template features, you can follow the [pongo2](https://github.com/flosch/pongo2) project

## FAQ

### Get Support

* if you have questions when using this tool, feel free to send me an email: `techmesh@aliyun.com.`

