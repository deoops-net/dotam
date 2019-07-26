# Dotam

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
* 内置方便的插件git, docker等
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
        password = "some key takes you home"
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
brew tap deoops-net/tap

brew install dotam
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







