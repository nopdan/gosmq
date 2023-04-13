
<a name="v2.1"></a>
## [v2.1](https://github.com/flowerime/gosmq/compare/v2.0...v2.1) (2023-04-13)

### Feat

* Parallel 使用 callback
* ParallelMerge

<a name="v2.0"></a>
## [v2.0](https://github.com/flowerime/gosmq/compare/v1.4...v2.0) (2023-03-19)

### Feat

- parallel
- 分文件夹输出及 html 输出
- go 命令移动到 root
- `--clean` 只统计码表中有的词
- `--merge` 合并多文本的结果

<a name="v1.4"></a>
## [v1.4](https://github.com/flowerime/gosmq/compare/v1.3...v1.4) (2023-03-17)

### Feat

- 整合 go 和 multi 命令
- parallel running
- 新的文本分割方法

### Fix

- 分词结果顺序
- 手指统计
- 词条数据统计错误
- 修复多码表并发错误
- 上上个提交漏了点

### Perf

- 优化分词结果输出性能
- 按词长稳定排序
- 用二维数组索引当量表
- 优化命令行输出
- 使用 buffer 读取文本

### Refactor

- 默认贪婪匹配
- 拆分 verbose

<a name="v1.3"></a>
## [v1.3](https://github.com/flowerime/gosmq/compare/v1.2...v1.3) (2023-03-16)

### Feat

- multi 支持多协程处理
- serve 支持递归遍历
- 增加 multi 命令

### Perf

- 优化手感统计逻辑与性能
- 优化命令行输出
- 缺字码长改为 4
- smq.Load 返回异常
- 调整 result 输出文件名
- 优化文件夹遍历

<a name="v1.2"></a>
## [v1.2](https://github.com/flowerime/gosmq/compare/v1.1...v1.2) (2023-03-15)

### Fix

- match 符号 切片越界

### Perf

- 更精准的非汉字判断
- 更改命令行参数
- 优化性能

<a name="v1.1"></a>
## [v1.1](https://github.com/flowerime/gosmq/compare/v1.0...v1.1) (2023-03-14)

### Feat

- 添加赛码的交互模式
- 双击打开 web 页面

### Fix

- 标点符号不算作打词
- 按字符数而不是字节数排序

<a name="v1.0"></a>
## [v1.0](https://github.com/flowerime/gosmq/compare/v0.30...v1.0) (2023-02-04)

- 重构 cli 工具
- 分离赛码表转换和赛码
- 细节优化

<a name="v0.30"></a>
## [v0.30](https://github.com/flowerime/gosmq/compare/v0.28...v0.30) (2022-09-02)

### Fix

- 上个提交写的 bug
- 修复 utf-16 首行读取错误
- 拇指统计错误
- 极速格式 10 重
- 极速格式

### Perf

- 又加回了默认格式
- 修改背景色

<a name="v0.28"></a>
## [v0.28](https://github.com/flowerime/gosmq/compare/v0.27...v0.28) (2022-09-02)

### Perf

- 删除自有赛码表格式
- 修改默认匹配算法
- 更改默认格式
- 完全整合两个前端

<a name="v0.27"></a>
## [v0.27](https://github.com/flowerime/gosmq/compare/v0.26...v0.27) (2022-08-29)

### Feat

- 整合了旧的 web 前端

### Fix

- dict.length
- 码表格式转换

### Perf

- 添加示例词库和文本

<a name="v0.26"></a>
## [v0.26](https://github.com/flowerime/gosmq/compare/v0.25...v0.26) (2022-08-28)

- 码表里有相同词时取码长较短的
- 整合了 cli 和 web 程序，可以双击直接启动
- 前端开源 @yyb1rd

<a name="v0.25"></a>
## [v0.25](https://github.com/flowerime/gosmq/compare/v0.24...v0.25) (2022-05-31)

### Feat

- 新的 web 前端 by [@yyb1rd](https://github.com/yyb1rd)
- 输出详细数据

### Fix

- 不输出详细数据会覆盖空数据的 bug
- 首选词统计错误
- 换行符错误

<a name="v0.24"></a>
## [v0.24](https://github.com/flowerime/gosmq/compare/v0.23...v0.24) (2022-04-30)

- 增加 冰凌码表格式
- 修复 极点格式转换错误
- 修复 带 BOM 文件首行读取错误
- cli 支持多码表
- 增加 首选词统计
- 优化 `go-pretty` 改到 v6 版本，体积减少大半

<a name="v0.23"></a>
## [v0.23](https://github.com/flowerime/gosmq/compare/v0.22...v0.23) (2022-04-29)

- 使用 `embed` 初始化 `puncts` 数据
- 优化 `longest`
- 使用前缀树重写了 `order` ，速度提升了几个数量级

<a name="v0.22"></a>
## [v0.22](https://github.com/flowerime/gosmq/compare/v0.14...v0.22) (2022-04-29)

几乎完全重写了

- `Result` 结构体更加规范，方便外部程序解析
- 多种码表格式支持，极速赛码表、多多、极点
- 自定义码表格式，提供 Transfer 接口
- 多种匹配算法，字典树、顺序匹配、最长匹配
- 自定义匹配算法，提供 Matcher 接口
- 左右空格可以分别统计了
- 可以使用专用赛码表格式，速度更快

<a name="v0.14"></a>
## [v0.14](https://github.com/flowerime/gosmq/compare/v0.13...v0.14) (2021-11-30)

- 合并了 cli 和 web
- 修复了读取文件的 bug

<a name="v0.13"></a>
## [v0.13](https://github.com/flowerime/gosmq/compare/v0.12...v0.13) (2021-11-30)

- 添加了一个 web 前端
- 添加了一个 web-server 前端

<a name="v0.12"></a>
## [v0.12](https://github.com/flowerime/gosmq/compare/v0.11...v0.12) (2021-11-28)

- 支持更多文件编码格式（GBK，UTF-16 等，其它的自行测试）
- 改变了输出结果的格式

<a name="v0.11"></a>
## [v0.11](https://github.com/flowerime/gosmq/compare/v0.10...v0.11) (2021-09-28)

- 重写 html 模版
- 修复 手指频率统计错误
- 重写 默认符号
- 优化 按键统计

<a name="v0.10"></a>
## [v0.10](https://github.com/flowerime/gosmq/compare/v0.8...v0.10) (2021-09-23)

- 添加 html 结果输出，感谢 @yyb1rd
- 添加 多码表同时测试
- 优化 性能

<a name="v0.8"></a>
## [v0.8](https://github.com/flowerime/gosmq/compare/v0.7...v0.8) (2021-08-14)

添加 选重统计

<a name="v0.7"></a>
## [v0.7](https://github.com/flowerime/gosmq/compare/v0.6...v0.7) (2021-08-12)

添加 按键频率统计

<a name="v0.6"></a>
## [v0.6](https://github.com/flowerime/gosmq/compare/v0.5...v0.6) (2021-08-10)

- 添加 当量 大小跨排 小指干扰 错手
- 添加 指法统计开关

<a name="v0.5"></a>
## [v0.5](https://github.com/flowerime/gosmq/compare/v0.4...v0.5) (2021-08-09)

- 添加 赛码表输出
- 添加 只跑单字
- 美化输出结果

<a name="v0.4"></a>
## [v0.4](https://github.com/flowerime/gosmq/compare/v0.3...v0.4) (2021-08-07)

添加指法统计

<a name="v0.3"></a>
## [v0.3](https://github.com/flowerime/gosmq/compare/v0.2...v0.3) (2021-08-05)

采用 hash 字典树进行字符串匹配

<a name="v0.2"></a>
## [v0.2](https://github.com/flowerime/gosmq/compare/v0.1...v0.2) (2021-08-04)

支持 rime 格式码表

<a name="v0.1"></a>
## v0.1 (2021-08-01)

initial release
