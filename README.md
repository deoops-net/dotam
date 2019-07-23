# Dotam

Dotam 读法 dot-am, 原意为dotfiles automation。

## 为什么创建另一个自动化工具

我自己在团队中的角色既是一个开发者（负责开发实现业务逻辑），同时我也是一个Devops（我需要自己写ci files,部署到kubernetes的相关yamls等等）。
我在做这些事情的时候通常会遇到在提交代码的时候同步一些状态到很多配置文件中，比如：我需要给Docker镜像的版本写到`Makefile`中，用于发布到kubernetes的`yaml`中还有`.drone.yml`
中，有的时候我总是会疏忽某一个文件，所以就需要再次修改，提交。。真的是很麻烦

我做了一些搜索，但是没有发现很符合我使用场景的自动化工具，所以我就写了这个。你可以往下看看它的主要特性是否符合你的使用场景来决定是否要使用和这个自动化工具

## 特性

* 编程语言无关
* 配置简单
* 支持多配置语法`json,yml,hcl`
* 内置方便的插件git, docker等
* 简单的自定义插件化

PS: 构建效率目前还不是主要开发方向，在版本稳定后我会花更多的时间来优化构建效率，当然如果你愿意参与贡献，我是举双脚欢迎的

## 快览

通常你只需要在项目跟目录下配置一个下面的这样的示例配置文件，然后运行`dotam build`,所有的工作就完成了！
PS: 如果你懒得每次从别的项目拷贝一个配置文件过来那么你可以通过`dotam init`来初始一个模板文件。

```hcl
temp "Makefile" {
    src = "conf"
    dest = "./"
    var {
        version = "{{ versions.prod }}"
        tag = "0.1.2"
    }
}

plugin "docker" {
    command = "docker"
    args = ["build", "-t", "{{docker.repo}}", "{{version.prod}}", "."]
    settings {
        version = "{{ versions.prod }}"
        passed = "{{ status.build_pass }}"
    }
}

var "versions" {
    prod = "v1.0.0"
    stage = "v1.0.3"
}

var "docker" {
    repo = "deoops/dotam"
}

```

## 安装

目前只支持go环境安装，后面我会发布更多的预编译版本到各个平台

`go install github.com/deoops-net/dotam`
