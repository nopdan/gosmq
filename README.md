# go 赛码器

不指定参数或双击运行新 web 前端。

在 [example](./example) 文件夹下查看预览图。

可以使用多个 `-i` 指定多个码表，此时其它参数都是共享的，所以其它码表的格式最好先转换过来。

不指定 `-t` 只指定 `-i` 时，可以初始化码表并转换输出。

```
Application Options:
  /t, /text:     string 文本路径
  /i, /input:    []string 码表路径
  /s, /single    bool   单字模式
  /f, /format:   string 码表格式 default|jisu,js|duoduo,dd|jidian,jd|bingling,bl
      /select:   string 自定义选重键
  /p, /push:     int    普通码表起顶码长，码长大于等于此数，首选不会追加空格
  /a, /alg:      string 匹配算法 t|s|l
  /k, /space:    string 空格按键方式 left|right|both
  /d, /details   bool   详细数据
  /v, /version   bool   查看版本信息

Help Options:
  /?             Show this help message
  /h, /help      Show this help message
```

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
