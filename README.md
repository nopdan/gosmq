# [最快的赛码器](https://github.com/imetool/gosmq)

![](assets/preview-serve.png)

[![GitHub Repo stars](https://img.shields.io/github/stars/imetool/gosmq)](https://github.com/imetool/gosmq/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/imetool/gosmq)](https://github.com/imetool/gosmq/network/members)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/imetool/gosmq)](https://github.com/imetool/gosmq/releases)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/imetool/gosmq/build.yml)](https://github.com/imetool/gosmq/actions/workflows/build.yml)
![GitHub repo size](https://img.shields.io/github/repo-size/imetool/gosmq)
![GitHub](https://img.shields.io/github/license/imetool/gosmq)

在 [assets](./assets) 文件夹下查看更多预览图。

## 用法

详细参数见 `.\smq.exe -h`。

执行`.\smq.exe serve`，启动`serve`模式，自动打开浏览器。

`serve`模式只会读取 `text` 文件夹下的赛文和 `dict` 文件夹下的码表。

### 参数解释

#### 默认码表格式

`字词 tab 真实输入码 tab 重码顺序`

> 首选的 1 可以省略。

#### 匹配算法

默认为按码表从上往下匹配。

贪心匹配：按照词长进行匹配，与码表中的顺序无关，相同词选编码短的。

#### 输出详细数据

输出到 `result` 文件夹内，包含赛码结果、分词结果、词条数据。

<!-- ### 示例 -->

## Benchmark

> 这是老版本的 Benchmark，新版要慢一点，不会差太多  
> 配置：windows 10, r5 3600 4.2g, 8g\*2 2933Mhz c18

| 文本         | 文本字数 |             码表 | 码表词条数 |  耗时 |
| ------------ | -------: | ---------------: | ---------- | ----: |
| 心情决定事情 |     9.6w | 091 点儿 2021 夏 | 20w        | 120ms |
| 红楼梦原著   |      87w | 091 点儿 2021 夏 | 20w        | 250ms |
| 极品全能高手 |    2164w | 091 点儿 2021 夏 | 20w        |  3.8s |
| 极品全能高手 |    2164w |   红辣椒五笔码表 | 880w       | 12.6s |
| 心情决定事情 |     9.6w |   红辣椒五笔码表 | 880w       |  6.7s |
