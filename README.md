# [最快的赛码器](https://github.com/imetool/gosmq)

![](assets/preview-serve.png)

[![GitHub Repo stars](https://img.shields.io/github/stars/imetool/gosmq)](https://github.com/imetool/gosmq/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/imetool/gosmq)](https://github.com/imetool/gosmq/network/members)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/imetool/gosmq)](https://github.com/imetool/gosmq/releases)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/imetool/gosmq/build.yml)](https://github.com/imetool/gosmq/actions/workflows/build.yml)
![GitHub repo size](https://img.shields.io/github/repo-size/imetool/gosmq)
![GitHub](https://img.shields.io/github/license/imetool/gosmq)

在 [assets](./assets) 文件夹下查看更多预览图。

## ！请在命令行中运行本程序

## 本赛码器专用格式

`字词 tab 真实输入码 tab 重码顺序`

> 首选的 1 可以省略。

## **用法**

### `serve` 命令

执行`.\smq.exe serve`，启动`serve`模式，自动打开浏览器。

`serve`模式只会读取 `text` 文件夹下的文本和 `dict` 文件夹下的码表。

### `gen` 格式转换

使用 `.\smq.exe gen` 命令转换格式。

支持格式：_极速赛码表(jisu)_ 和 [_其他_](https://github.com/imetool/dtool#%E5%AD%97%E8%AF%8D%E7%A0%81%E8%A1%A8%E8%BD%AC%E6%8D%A2)

### 主命令参数解释

#### 文本和码表

`-i` 和 `-t` 用法一样

examples:

`-i 文件1 -i 文件2` 载入两个码表  
`-i 文件夹1 -i 文件夹2` 载入文件夹 1 和文件夹 2 内所有码表  
`-i 文件1 -i 文件夹1` 载入文件 1 和文件夹 1 内所有码表

> 载入多个文本，不保证输出结果有序。

#### 匹配算法

默认贪心匹配，按照词长进行匹配，与码表中的顺序无关，相同词选择**码长较短**的。

指定 `--stable`，按照码表顺序匹配，相同词选择**靠前**的。

指定 `--single` 或 `-s`，单字模式，只取码表中的单字，同一个字**码长较短**的。

#### 输出详细数据

- `--split`: 输出分词数据
- `--stat`: 输出词条数据
- `--json`: 输出 json 数据
- `--html`: 保存 html 结果

`--verbose` 或 `-v` 输出所有数据，可以追加 `--json=false` 关闭其中某项。

#### 合并多文本的结果

使用 `--merge` 或 `-m` 合并多文本的结果，这时**输出分词结果**不再生效。

<!-- ### 示例 -->

#### _关于匹配逻辑_

若某个字符码表中匹配不到

默认情况：

- 是标点符号
  - 根据内置的符号表继续匹配
- 不是标点符号：编码设置为 `####`
  - 是汉字：记为缺字

指定 `--clean` 或 `-c`: 跳过该字符

## Benchmark

> 配置：windows 10, i7-9750H(6c12t), 8g\*2 2666Mhz

> 以下采用码表《091 点儿 2023 春》，词条数 214338

```powershell
.\smq.exe -i .\dict\091点儿2023春.txt -t .\text\心情决定事情.txt
.\smq.exe -i .\dict\091点儿2023春.txt -t .\text\红楼梦原著.txt
.\smq.exe -i .\dict\091点儿2023春.txt -t .\text\《庆余年》.txt
.\smq.exe -i .\dict\091点儿2023春.txt -t .\text\那些热血飞扬的日子（整理版）.txt
.\smq.exe -i .\dict\091点儿2023春.txt -t .\text\极品全能高手_花都大少.txt
```

| 文本               | 文本字数 |  耗时 |
| ------------------ | -------: | ----: |
| 心情决定事情       |    96253 | 120ms |
| 红楼梦原著         |   872209 | 160ms |
| 《庆余年》         |  3464055 | 250ms |
| 那些热血飞扬的日子 | 16485176 | 630ms |
| 极品全能高手       | 24910973 | 870ms |

> 以下采用码表《红辣椒五笔码表》，词条数 8896705

```powershell
.\smq.exe -i .\dict\红辣椒五笔码表880万多多格式.txt -t .\text\心情决定事情.txt
.\smq.exe -i .\dict\红辣椒五笔码表880万多多格式.txt -t .\text\极品全能高手_花都大少.txt
载入码表： 红辣椒五笔码表880万多多格式
```

| 文本         | 文本字数 | 耗时 |
| ------------ | -------: | ---: |
| 心情决定事情 |    96253 | 8.6s |
| 极品全能高手 | 24910973 | 9.6s |

> 多文件测试

```powershell
# 使用 hidden 隐藏输出
.\smq.exe -i .\dict\091点儿2023春.txt -t .\super\data\ --hidden

# 载入码表： 091点儿2023春
# 共载入 26078 个文本，总字数 192108228，总耗时：10.006676s
```

> 输出详情测试

```powershell
.\smq.exe -i .\dict\091点儿2023春.txt -t .\text\极品全能高手_花都大少.txt -v

# cost time: 1.35s
```
