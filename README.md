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

`serve`模式只会读取 `text` 文件夹下的文本和 `dict` 文件夹下的码表。

### 参数解释

#### 默认码表格式

`字词 tab 真实输入码 tab 重码顺序`

> 首选的 1 可以省略。

#### 匹配算法

默认为按码表顺序从上往下匹配。

贪心匹配：按照词长进行匹配，与码表中的顺序无关，相同词选编码短的。

#### 输出详细数据

输出到 `result` 文件夹内，包含赛码结果、分词结果、词条数据。

<!-- ### 示例 -->

## Benchmark

```powershell
.\smq.exe go -i .\dict\092K.txt -t .\text\心情决定事情.txt
.\smq.exe go -i .\dict\092K.txt -t .\text\红楼梦原著.txt
.\smq.exe go -i .\dict\092K.txt -t .\text\极品全能高手_花都大少.txt

.\smq.exe go -i .\dict\红辣椒五笔码表880万多多格式.txt -t .\text\心情决定事情.txt
.\smq.exe go -i .\dict\红辣椒五笔码表880万多多格式.txt -t .\text\极品全能高手_花都大少.txt
```

> 配置：windows 10, i7-9750H(6c12t), 8g\*2 2666Mhz

| 文本         | 文本字数 |           码表 | 码表词条数 |  耗时 |
| ------------ | -------: | -------------: | ---------- | ----: |
| 心情决定事情 |    96253 |           092K | 265843     | 160ms |
| 红楼梦原著   |   872209 |           092K | 265843     | 185ms |
| 极品全能高手 | 24910973 |           092K | 265843     |    1s |
| 心情决定事情 |    96253 | 红辣椒五笔码表 | 8896705    |  8.5s |
| 极品全能高手 | 24910973 | 红辣椒五笔码表 | 8896705    |  9.6s |
