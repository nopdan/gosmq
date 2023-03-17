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

默认贪心匹配，按照词长进行匹配，与码表中的顺序无关，相同词选择靠前的。

若要按照码表顺序，指定 `--stable`

#### 输出详细数据

输出到 `result` 文件夹内，包含赛码结果、分词结果、词条数据。

<!-- ### 示例 -->

## Benchmark

> 配置：windows 10, i7-9750H(6c12t), 8g\*2 2666Mhz

> 以下采用码表《091 点儿 2023 春》，词条数 214338

```powershell
.\smq.exe go -i .\dict\091点儿2023春.txt -t .\text\心情决定事情.txt
.\smq.exe go -i .\dict\091点儿2023春.txt -t .\text\红楼梦原著.txt
.\smq.exe go -i .\dict\091点儿2023春.txt -t .\text\《庆余年》.txt
.\smq.exe go -i .\dict\091点儿2023春.txt -t .\text\那些热血飞扬的日子（整理版）.txt
.\smq.exe go -i .\dict\091点儿2023春.txt -t .\text\极品全能高手_花都大少.txt
```

| 文本               | 文本字数 |  耗时 |
| ------------------ | -------: | ----: |
| 心情决定事情       |    96253 | 120ms |
| 红楼梦原著         |   872209 | 160ms |
| 《庆余年》         |  3464055 | 250ms |
| 那些热血飞扬的日子 | 16485176 | 690ms |
| 极品全能高手       | 24910973 | 999ms |

> 以下采用码表《红辣椒五笔码表》，词条数 8896705

```powershell
.\smq.exe go -i .\dict\红辣椒五笔码表880万多多格式.txt -t .\text\心情决定事情.txt
.\smq.exe go -i .\dict\红辣椒五笔码表880万多多格式.txt -t .\text\极品全能高手_花都大少.txt
载入码表： 红辣椒五笔码表880万多多格式
```

| 文本         | 文本字数 | 耗时 |
| ------------ | -------: | ---: |
| 心情决定事情 |    96253 | 8.6s |
| 极品全能高手 | 24910973 | 9.6s |

> 多文件测试

```powershell
# 使用 hidden 隐藏输出
.\smq.exe go -i .\dict\091点儿2023春.txt -t .\super\data\ --hidden

# 载入码表： 091点儿2023春
# 共载入 26078 个文本，总字数 192108228，总耗时：10.006676s
```

> 输出详情测试

```powershell
.\smq.exe go -i .\dict\091点儿2023春.txt -t .\text\极品全能高手_花都大少.txt -v

# cost time: 2.3s
```
