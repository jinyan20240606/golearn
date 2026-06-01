# Go 测试核心：`testing` 包 + `*testing.T` 完全笔记（**Go 自带，无需安装**）

## 一、核心结论（先背）
1. **`testing` 是 Go 语言标准库自带的**，**不需要安装任何第三方包**
2. **测试函数必须传入 `t *testing.T`**，所有测试能力都来自这个 `t`
3. 测试函数命名规则：**必须以 `Test` 开头**，格式：
   ```go
   func TestXxx(t *testing.T) {
       // 测试逻辑
   }
   ```

---

## 二、`*testing.T` 最常用 8 个方法（99% 场景够用）
### 1. `t.Log()` / `t.Logf()`
**打印日志**（只有测试失败或加 `-v` 才显示）
```go
t.Log("开始测试用户查询")
t.Logf("用户昵称：%s", user.NickName)
```

### 2. `t.Error()` / `t.Errorf()`
**标记测试失败，继续执行后续代码**
```go
if user.NickName != "bobby18" {
    t.Error("昵称不匹配")
    t.Errorf("期望：%s，实际：%s", "bobby18", user.NickName)
}
```

### 3. `t.Fatal()` / `t.Fatalf()`
**标记测试失败，立刻终止当前测试**（后面代码不执行）
**必须用在：出错后没必要继续执行的场景**
```go
if err != nil {
    t.Fatal("获取用户失败")
    t.Fatalf("失败原因：%v", err)
}
```

### 4. `t.Fail()`
标记失败，但继续跑
```go
t.Fail()
```

### 5. `t.FailNow()`
标记失败，立刻停止
```go
t.FailNow()
```

### 6. `t.Skip()`
**跳过当前测试**
```go
t.Skip("跳过该测试用例")
```

### 7. `t.Run()`
**子测试**：一个测试函数里跑多个用例
```go
t.Run("正常查询", func(t *testing.T) { ... })
t.Run("用户不存在", func(t *testing.T) { ... })
```

### 8. `t.Helper()`
**封装断言工具函数时用**（让报错行指向测试用例，而非工具函数）
```go
func assertEqual(t *testing.T, expect, actual string) {
    t.Helper() // 关键
    if expect != actual {
        t.Errorf("期望：%s，实际：%s", expect, actual)
    }
}
```

---

## 三、最核心区别：**Error vs Fatal**（笔记必记）
| 方法 | 效果 | 场景 |
|---|---|---|
| **t.Error** | 失败 + **继续执行** | 断言不影响后续流程 |
| **t.Fatal** | 失败 + **立即停止** | 出错了就没必要往下跑 |

**示例：**
```go
user, err := getUser()
if err != nil {
    t.Fatalf("获取失败：%v", err) // 必须停
}
if user.Name != "test" {
    t.Error("名字错误") // 可以继续跑
}
```

---

## 四、完整测试模板（万能模板）
```go
package demo

import "testing"

func TestUserService(t *testing.T) {
    t.Log("测试开始")

    // 1. 准备数据

    // 2. 调用方法
    // user, err := ...

    // 3. 断言
    if err != nil {
        t.Fatalf("调用失败：%v", err)
    }
    if user.Name != "test" {
        t.Errorf("名称不匹配，期望：test，实际：%s", user.Name)
    }

    t.Log("测试通过")
}
```

---

## 五、极简总结（笔记最后一行）
- **`testing` 是 Go 官方自带测试库**
- **`t *testing.T` 是测试核心对象**
- **记住 4 个方法 = 搞定所有测试**
  - `t.Log`  打印
  - `t.Fatal` 失败并终止
  - `t.Error` 失败但继续
  - `t.Run`   子测试

---

## go test 命令 **所有用法**（极简、完整、笔记专用）
```go test [flags] [packages]```

### 一、基础运行
- `go test`              运行当前包测试
- `go test -v`           详细输出
- `go test ./...`        运行当前目录+所有子包
- `go test ./service`    运行指定包
- `go test -run 正则`     只运行匹配的单元测试
- `go test -timeout 10s` 设置超时

### 二、失败控制
- `go test -failfast`    第一个失败立即停止
- `go test -v -failfast` 详细+快速失败

### 三、覆盖率
- `go test -cover`                 显示覆盖率
- `go test -coverprofile=cover.out` 生成覆盖率文件
- `go tool cover -html=cover.out`   打开HTML覆盖率报告

### 四、基准测试（性能）

- Go 自带的性能测试工具
  - 用来：
  - 测代码运行速度
  - 测内存分配
  - 测CPU 消耗
- 基准测试规则
  - 函数名必须以 Benchmark 开头
  - 参数是 b *testing.B
  - 循环 b.N 次（系统自动控制）
  - 放在 _test.go 文件里

- `go test -bench .`               运行所有基准测试
- `go test -bench BenchmarkX`      运行指定基准测试
- `go test -benchmem`              显示内存分配
- `go test -benchtime 10s`        运行时长
- `go test -count 5`               重复运行5次
- `go test -cpu 1,2,4`             指定CPU核心数
- `go test -run ^$ -bench .`       只跑基准测试，不跑单元测试

### 五、性能分析（pprof）
- `-cpuprofile cpu.out`    CPU分析
- `-memprofile mem.out`    内存分析
- `-blockprofile block.out` 阻塞分析
- `-mutexprofile mutex.out` 互斥锁分析

### 六、缓存/清理
- `go test -count=1`        禁用缓存，强制执行
- `go clean -testcache`     清空测试缓存

### 七、输出/编译
- `go test -o test.bin`     编译成二进制文件（不执行）
- `go test -c`              只编译，不运行
- `go test -short`          运行简短测试
