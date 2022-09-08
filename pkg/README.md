
支持多码表一起赛，支持多种码表格式，支持多种匹配算法，可自定义码表格式和匹配算法。

## 开始之前

下载：`go get -u github.com/cxcn/gosmq/pkg/smq`  

导入：`import "github.com/cxcn/gosmq/pkg/smq"`

## 使用

### 创建赛码器

你可以从字符串、文件路径创建赛码器

```go
func main(){
    // 从文件路径
    // 第一个参数可以填空，这时的文本名从路径推断
    s1 := smq.New("","这里填写路径")
    // 从字符串
    s2:= smq.NewFromString("文本名","这里是文本内容")
}
```

### 创建一个码表

```go
func main(){
    // 先定义一些基本参数
    dict := &dict.Dict{
        // 具体定义在 dict/struct.go 里
    }

    // 载入码表
    // 从文件路径
    dict.Load("这里填写路径")
    // 从字符串
    dict.LoadFromString("这里是码表内容")
}
```

### 添加码表到赛码器

`s.Add(dict)`

### 开始比赛

```go
func main() {
    s := smq.New(...)
    dict := &smq.Dict{...}
    dict.Load(...)
    s.Add(dict)
    res := s.Run() // 他返回一个 []*smq.Result 结构体指针数组，具体定义可查看 result.go 文件
    // 你可以输出为 json
    j,_ := json.Marshal(res)
}

```

## 自定义码表转换

```go
// 需要从 dict 生成本赛码器格式码表的字节数组
type Transformer interface {
    Read(*dict.Dict) []dict.Entry
}
```

一个例子：

```go
type newFormat struct {}

func (n newFormat) Read(dict *dict.Dict) []dict.Entry {
	ret := make([]dict.Entry, 0, 1e5)
    // 逐行读取
    scan := bufio.NewScanner(dict.Reader)
    for scan.Scan() {
        // 转格式为
        // 本赛码器格式 word code order
		ret = append(ret, dict.Entry{word, code, order})
    }
    return ret
}
```

## 自定义匹配算法

程序定义了一个接口，在创建 `dict.Dict` 时传入 `Matcher`，这时 `Algorithm` 失效

```go
type Matcher interface {
    // 插入一个词条 word code order
    Insert(string, string, int)
    // 匹配下一个词 text point -> 匹配到的词长，code，order
    Match([]rune, int) (int, string, int)
}
```
