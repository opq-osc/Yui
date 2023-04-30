---
layout: home

hero:
  name: Yui
  text: OPQBot 的 golang 版本的插件实现
  tagline: 插件管理，权限控制，安全...
#  image:
#    src: /opq.png
#    alt: OPQ
  actions:
    - theme: brand
      text: 开始使用
      link: /guide/开始
    - theme: alt
      text: GitHub
      link: 
features:
  - icon: 🍭
    title: 简单，易用
    details: 小白只需要下载对应平台文件，配置OPQ链接，即可运行
  - icon: 🛡
    title: 安全
    details: 插件运行于WASI沙箱中，所有API调用均需要在Meta文件中申请
  - icon: 🚀
    title: 低开发成本
    details: 插件编译为OPQ文件后可以在所有Yui支持的平台上运行，无须再次编译
---