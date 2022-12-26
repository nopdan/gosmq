# [最快的赛码器](https://github.com/cxcn/gosmq)

![](example/preview-serve.png)


[![GitHub Repo stars](https://img.shields.io/github/stars/cxcn/gosmq)](https://github.com/cxcn/gosmq/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/cxcn/gosmq)](https://github.com/cxcn/gosmq/network/members)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/cxcn/gosmq)](https://github.com/cxcn/gosmq/releases)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/cxcn/gosmq/build.yml)](https://github.com/cxcn/gosmq/actions/workflows/build.yml)
![GitHub repo size](https://img.shields.io/github/repo-size/cxcn/gosmq)
![GitHub](https://img.shields.io/github/license/cxcn/gosmq)


在 [example](./example) 文件夹下查看更多预览图。

## web 用法

不指定参数或双击运行，自动打开浏览器，`.\smq.exe serve` 可以静默启动。  
只会读取 `text` 文件夹下的赛文和 `dict` 文件夹下的码表。  
不支持码表输出和详细结果输出。

## cli 用法

使用多个 `-i` 指定可以多个码表，其它参数是共享的。  
非默认格式会自动转换并输出到 `dict` 文件夹下。

```
Application Options:
  /t, /text:     string 文本路径
  /i, /input:    []string 码表路径
  /s, /single    bool   单字模式
  /f, /format:   string 码表格式 default|jisu,js|duoduo,dd|jidian,jd|bingling,bl
      /select:   string 自定义选重键
  /p, /push:     int    普通码表起顶码长，码长大于等于此数，首选不会追加空格
  /a, /alg:      string 匹配算法 trie,t|strie,s|longest,l
  /k, /space:    string 空格按键方式 left|right|both
  /d, /detail    bool   输出详细数据
  /o, /output    bool   输出转换后的码表
  /v, /version   bool   查看版本信息

Help Options:
  /?             Show this help message
  /h, /help      Show this help message
```

### 参数解释

#### 默认码表格式

> 自定义选重键和普通码表起顶码长对此格式无效。

`字词 tab 真实输入码 tab 重码顺序`

#### 匹配算法

最长匹配（默认）：按照词长进行匹配，相同词选编码短的，推荐多多格式选择。  
按码表顺序：字面意思，按码表从上往下匹配，推荐极速格式选择。

#### 输出详细数据

输出到 `result` 文件夹内，包含赛码结果、分词结果、词条数据。

### 示例

- 默认格式码表：`.\smq.exe -t text\text.txt -i dict\dict.txt`
- 极速格式码表：`.\smq.exe -t text\text.txt -i dict\dict.txt -f js`
- 多多格式二码顶码表：`.\smq.exe -t text\text.txt -i dict\dict.txt -f dd -p 2`
- 冰凌格式码表，单字模式，自定义选重键，指定顺序匹配，左手空格，输出详细数据，不输出转换后的码表：`.\smq.exe -t text\text.txt -i dict\dict.txt -f bl -s --select "_z;'[" -a s -k left -d --no`

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
