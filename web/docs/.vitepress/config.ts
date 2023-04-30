// .vitepress/config.js
import { defineConfig } from 'vitepress'
export default defineConfig({
    // site-level options
    title: 'Yui 文档站',
    description: 'Yui 是 OPQBot 的 golang 版本的插件实现',
    markdown: {
        lineNumbers: true
    },
    themeConfig: {
        nav: [
            { text: 'OPQBot', link: 'https://opqbot.com/' },
            { text: 'OPQBot Go SDK', link: 'https://github.com/opq-osc/OPQBot' }
        ],
        search: {
            provider: 'local'
        },
        sidebar: [
            {
                text: 'Guide',
                items: [
                    { text: '开始使用', link: '/guide/开始' },
                    { text: '管理员命令', link: '/guide/管理员命令' },
                ]

            },
            {
                text: '插件编写',
                items: [
                    { text: '环境配置', link: '/插件编写/环境配置' },
                    { text: '权限申请', link: '/插件编写/权限申请' },
                    { text: '事件', link: '/插件编写/事件' },
                    { text: '插件API', link: '/插件编写/插件API' },
                ]
            }
        ],
        footer: {
            message: 'OPQ Open Source Community',
            copyright: 'MIT Licensed | Copyright © 2023'
        },
        editLink: {
            pattern: ({ relativePath }) => {
                // @ts-ignore
                if (relativePath.startsWith('packages/')) {
                    return `hhttps://github.com/opq-osc/Yui/main/web/${relativePath}`
                } else {
                    return `https://github.com/opq-osc/Yui/main/web/docs/${relativePath}`
                }
            }
        }
    },
})