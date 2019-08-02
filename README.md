# Dotam

[English Docs](./README.en.md)

Dotam 读法 dot-am, 原意为dotfiles automation。

<!-- TOC -->

- [Dotam](#dotam)
    - [为什么创建另一个自动化工具](#为什么创建另一个自动化工具)
    - [特性](#特性)
    - [快览](#快览)
    - [安装](#安装)
        - [**homebrew**](#homebrew)
        - [**预编译二进制**](#预编译二进制)
        - [**通过go安装**](#通过go安装)
        - [**Docker**](#docker)
        - [**开发及源码安装**](#开发及源码安装)
    - [使用](#使用)
        - [运行](#运行)
        - [初始化](#初始化)
    - [文档](#文档)
        - [Template 模块](#template-模块)
        - [Docker模块](#docker模块)
        - [Plugin](#plugin)
        - [模板语法](#模板语法)
    - [注意及常见问题](#注意及常见问题)
        - [获取支持](#获取支持)
        - [语法冲突](#语法冲突)

<!-- /TOC -->

## 为什么创建另一个自动化工具

我自己在团队中的角色既是一个开发者（负责开发实现业务逻辑），同时我也是一个Devops（我需要自己写ci files,部署到kubernetes的相关yamls等等）。
我在做这些事情的时候通常会遇到在提交代码的时候同步一些状态到很多配置文件中，比如：我需要给Docker镜像的版本写到`Makefile`中，用于发布到kubernetes的`yaml`中还有`.drone.yml`
中，有的时候我总是会疏忽某一个文件，所以就需要再次修改，提交。。真的是很麻烦

我做了一些搜索，但是没有发现很符合我使用场景的自动化工具，所以我就写了这个。你可以往下看看它的主要特性是否符合你的使用场景来决定是否要使用和这个自动化工具

## 特性

* 编程语言无关
* 配置简单，几乎没什么学习成本
* 支持多配置语言json,yml,hcl
* 内置方便的插件git, docker等
* 简单的自定义插件化

PS: 构建效率目前还不是主要开发方向，在版本稳定后我会花更多的时间来优化构建效率，当然如果你愿意参与贡献，我是举双脚欢迎的

## 快览

通常你只需要在项目根目录下配置一个像下面的这样的示例配置文件，然后运行`dotam build`,所有的工作就完成了！
PS: 如果你懒得每次从别的项目拷贝一个配置文件过来那么你可以通过`dotam init`来初始一个模板文件。

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

## 安装

### **homebrew**

```bash
brew install deoops-net/tap/dotam

// update
brew uninstall dotam
brew install deoops-net/tap/dotam
```

### **预编译二进制**

到[发布](https://github.com/deoops-net/dotam/releases)页面下载对应平台的预编译版本。

### **通过go安装**

`go install github.com/deoops-net/dotam`

### **Docker**

comming soon!

### **开发及源码安装**

```bash
git clone https://github.com/deoops-net/dotam
cd dotam
make test
make install
```

## 使用

### 运行

`dotam build [Dotamfile.{yml,hcl,json}]`

默认情况下直接使用dotam build即可，dotam会在根目录下寻找Dotamfile的三种格式文件(.json, .hcl, .yml|.yaml)之一。你也可以指定一个具体的
文件让dotam 去跑

### 初始化

`dotam init [-t yml|yaml,json,hcl]`

此命令会为项目创建一个模板配置文件，可以通过`-t`指定创建的文件类型。


## 文档

### Template 模块

通常我们需要维护很多项目相关的静态文件，但是他们多数情况下需要一些动态属性，这就好比每次发布版本我们需要同步版本号到各种文件中去。
比如：
* 我们需要给CI配置文件指定要构建的Docker镜像的版本
* 或者我们需要给用于`kubectl`部署`deployment`的yaml文件指定镜像版本
* 我们的`main`函数也可能会读取某个文件来输出版本日志

等等这些，每次维护起来就比较麻烦，所以我们可以通过**一处定义，多处同步**的方式来使用Template模块。

Template模块的使用比较简单，这里有一个例子：

```hcl
temp "Makefile" {
  # 指定从哪里读取模板文件
  src = ".dotam/Makefile"
  # 指定渲染后的文件的存放目录，一般是根目录即可
  dest = "."
  # 最后将模板用到变量定义好就可以了
  var {
    # 比如我们想传递给Makefile模板一个version变量 
    version = "{{versions.production}}"
  }

}
// 定义我们要传递的变量
var {
  versions {
    prodution = "v1.0.0"
  }
}
```


### Docker模块

Docker模块可以帮助我们在本地构建好镜像，并push到指定仓库。一个Docker模块的示例：

```hcl
docker {
  repo = "deoops/dotam"
  tag = "_args.version"
  // 如果需要制定某一个Dockerfile可以使用这个选项
  dockerfile = "Dockerfile"
  // 一些情况需要用的--build-args参数
  buildArgs = [
    {
      key = "foo"
      value = "bar"
    } 
  ]
  // 如果是一个没有设置认证的内部registry可以用次参数
  // notPrivate = true

  // auth 参数有两个作用
  // push 镜像时的认证
  // 调度 caporal时的api认证
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

这里我们引入了一个新的内置变量`_args`，此变量用于获取从命令行传递的参数并传递给模板。由于一些敏感信息不适合在
文档中记录，所以我们可以把这些数据通过命令行参数隐式的传递给dotam，dotam将会通过内置对象_args向下传递。所以
上面的配置的命令行命令看起来应该是这样：

```bash
$ dotam build reg_user=tom reg_pass=password
```

[caporal](https://github.com/deoops-net/caporal)目前还是一个早期的版本，主要用来提供一个远程的api server来帮助我们运行和更新我们的容器服务，虽然
现在功能还比较简单，但是这个项目还在持续的开发中，相信越来越多的功能将会集成进来。


### Plugin

**Not Implemented**

### 模板语法

通常的变量只需要通过{{variable}}的方式使用即可, 你可以参考`example/.dotam`目录下的示例用法。
本项目的模板引擎目前依赖于[pongo2](https://github.com/flosch/pongo2)这个项目, 如果你对一些高阶的模板语法有需求可以到此项目下查看更多的文档。

## 注意及常见问题

### 获取支持

* 你可以随时发送邮件到techmesh@aliyun.com来获取支持，我会在检查邮件时尽快回复。
* 你也可以通过邮件来索要我的个人微信来获取即时支持。

### 语法冲突

如果模板中还有一些用于其他工具的模板标记比如我们的`.drone.yml`中用于slack插件的`{{#success}}`语法，这个会和项目自带的
pongo2模板语法有冲突，你可以直接在模板中用`safe`过滤器来处理:

```yml
  - name: notify
    image: plugins/slack
    settings:
      webhook: https://hooks.slack.com/services/TKR84LDNK/BL0A07VEG/xasdaww
      channel: deoops
      link_names: true
      icon_url: https://unsplash.it/256/256/?random
      image_url: http://auto.wegeek.fun/api/badges/techmesh/dkb-api/status.svg
      template: >
        {{"{{#success build.status}}"|safe}}
        {{"**API** build {{build.number}} succeeded. <@dayuoba> ready to be deployed. <@Vincent> [Doc]login update "|safe}}
        {{"{{else}}"|safe}}
        {{" **API** build {{build.number}} failed. Fix {{build.link}} please <@dayuoba>. "|safe}}
        {{"{{/success}} "|safe}}
```

EDIT:
通过 [issues](https://github.com/flosch/pongo2/issues/218) 了解到这个才是更好的办法







