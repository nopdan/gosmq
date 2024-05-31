<div align="center">

<img src="assets/logo.png" width=150></img>

### 昙花赛码器 - 最快的赛码器

[![GitHub Repo stars](https://img.shields.io/github/stars/nopdan/gosmq)](https://github.com/nopdan/gosmq/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/nopdan/gosmq)](https://github.com/nopdan/gosmq/network/members)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/nopdan/gosmq)](https://github.com/nopdan/gosmq/releases)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/nopdan/gosmq/build.yml)](https://github.com/nopdan/gosmq/actions/workflows/build.yml)
![GitHub repo size](https://img.shields.io/github/repo-size/nopdan/gosmq)
![GitHub](https://img.shields.io/github/license/nopdan/gosmq)

在 [assets](./assets) 文件夹下查看预览图。

</div>

## 本赛码器专用格式

`字词 tab 真实输入码 tab 重码顺序`

> 首选的 1 可以省略。

## **用法**

### `server` 命令

执行`.\smq.exe server`，启动`server`模式，自动打开浏览器。

`server`模式只会读取 `text` 文件夹下的文本和 `dict` 文件夹下的码表。

手动放码表要注意格式，推荐从剪贴板添加码表，具体操作步骤为：

1. 粘贴码表
2. 选择对应格式，并设置起顶码长等选项
3. 点击赛码

跑一次赛码会自动保存转换后的码表到 dict 文件夹内。

> 剪贴板需要浏览器授予权限

### `convert` 命令转换格式

使用 `.\smq.exe convert` 命令转换格式。

支持格式：_极速赛码表(jisu|js)_、_多多|rime(duoduo|dd|rime)_、_极点|小小(jidian|jd|xiaoxiao|xx)_、_冰凌(bingling|bl)_

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

`--verbose` 或 `-v` 输出所有数据，可以追加 `--json=false` 关闭其中某项。

#### 合并多文本的结果

使用 `--merge` 或 `-m` 合并多文本的结果，这时**输出分词结果**不再生效。

<!-- ### 示例 -->

#### _关于匹配逻辑_

若某个字符码表中匹配不到

默认情况：

- 是标点符号
  - 根据内置的符号表继续匹配
- 不是标点符号：编码设置为 `######`
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
| 红楼梦原著         |   872209 | 130ms |
| 《庆余年》         |  3464055 | 195ms |
| 那些热血飞扬的日子 | 16485176 | 440ms |
| 极品全能高手       | 24910973 | 610ms |

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
# 使用 --merge 合并所有文本结果
.\smq.exe -i .\dict\091点儿2023春.txt -t .\super\data\ --merge

# 载入码表： 091点儿2023春
# 共载入 26078 个文本，总字数 192108228，总耗时：5.9521118s
```

> 输出详情测试

```powershell
.\smq.exe -i .\dict\091点儿2023春.txt -t .\text\极品全能高手_花都大少.txt -v

# cost time: 1.2575731s
```

## 编译

> 用到的工具链：`bun`, `go 1.22.0+`

### 编译前端

```sh
cd frontend
bun install
bun run build
cd ..
```
### 编译主程序

```sh
go mod tidy
go build
```
