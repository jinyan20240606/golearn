# 笔记

## Go fmt 包 10 个最常用方法

全部分成 打印类、拼接字符串类、扫描输入类 三类，一看就懂。

总结（最核心 5 个）
1. fmt.Println() → 打印并换行（最常用）
2. fmt.Printf() → 格式化打印
3. fmt.Sprintf() → 格式化拼字符串
4. fmt.Scanln() → 读取一行输入
5. fmt.Errorf() → 生成错误


### 一、打印类（最常用）
1. fmt.Print()
不换行，直接打印
多个参数用空格隔开
go
运行
fmt.Print("Hello", "Go") // 输出：HelloGo
1. fmt.Println()
打印并自动换行
开发最最最常用
go
运行
fmt.Println("你好", "Go") // 输出：你好 Go\n
1. fmt.Printf()
格式化打印（带占位符）
想控制输出格式必须用它
go
运行
name := "小明"
age := 20
fmt.Printf("姓名：%s，年龄：%d\n", name, age)
常用占位符：
%s 字符串
%d 整数
%f 浮点数
%t 布尔值
%v 自动识别类型
%+v 打印结构体带字段名
%#v Go 语法格式打印


### 二、字符串拼接类（不打印，只返回字符串）

1. fmt.Sprintf()
格式化生成字符串，不打印
你刚才问的就是这个
go
运行
str := fmt.Sprintf("年龄：%d", 18)
// str = "年龄：18"
1. fmt.Sprint()
直接转字符串，不用占位符
go
运行
s := fmt.Sprint("数字", 100, true)
// s = "数字100true"
1. fmt.Sprintln()
转字符串 + 自动加空格和换行
go
运行
s := fmt.Sprintln("a", 1, true)
// s = "a 1 true\n"

### 三、输入扫描类（从控制台读取用户输入）

1. fmt.Scan()
读取输入，赋值给变量
以空格 / 回车分隔
go
运行
var name string
fmt.Scan(&name) // 读取输入

2. fmt.Scanln()
读取一行输入，遇到回车停止
go
运行
var age int
fmt.Scanln(&age)

1. fmt.Scanf()
按格式读取输入，按你指定的格式，从控制台读取用户输入，精准匹配格式后赋值给变量
go
运行
var n int
fmt.Scanf("年龄：%d", &n)
// 只要你在终端按格式输入年龄：20，n变量就获得了20值

### 四、错误处理常用
1.  fmt.Errorf()
生成一个格式化错误信息
go
运行
err := fmt.Errorf("参数错误：%d", 10)

超简记忆口诀
想打印 → Print / Println / Printf
想拼字符串 → Sprint / Sprintln / Sprintf
想读输入 → Scan / Scanln / Scanf
想造错误 → Errorf