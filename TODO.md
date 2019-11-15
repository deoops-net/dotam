# TODO

叹号为优先级

## V 1.0.0

**目标：**

* 模板插件完整支持
* cli基本命令支持`build, init`
* plugin完整支持
* 完善文档

**todos:**

- [x] !! Temp插件
    - [ ] \_args参数优先
- [x] !! 多配置语言支持
    - [x] json
    - [x] yaml
    - [x] hcl
- [x] !! init命令
- [ ] !! 文档
    - [x] template syntax
        - [x] en
        - [x] zh
- [ ] ! 自举
- [ ] ! git 插件
    - [x] git hook
- [ ] !! docker 插件
    - [x] !! build
    - [x] !! push
    - [x] ! pass auth from command line
    - [ ] !! support multi tags for pushing @ refactor code
    - [ ] !! health probe
    - [ ] !! ready probe
    - [x] !! support --build-args
    - [x] !! support specifying Dockerfile
    
    
    
- [ ] ! refactor code
    - [ ] !! better file&package structure
- [ ] !! caporal
    - [x] 创建容器
    - [x] 更新容器
    - [x] 指定网络
    - [x] 指定volume
    - [x] 鉴权分离
    - [ ] more opts support
    
- [ ] !!三方插件支持
    
- [ ] ! 完善日志
- [x] !! build args
    - [x] ! 参数及conf文件自动解析

----

## V 2.0.0

**目标：**

* 兼容drone
