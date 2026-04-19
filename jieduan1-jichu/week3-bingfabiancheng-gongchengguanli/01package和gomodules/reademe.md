# 03周并发编程与工程管理

## package
1. go语言的中代码组织是通过package组织的，一个package可以包含多个go文件，一个go文件可以包含多个函数。
2. package用来组织源码，是多个go源码的集合，代码复用的基础，如fmt包，math包，os包等等。
3. 每个源码文件开始都必须要声明所属的package，如package main。
4. python中不需要去声明package因为它内部默认是按照文件名自动声明的，而 php，c#+，java，c#，go都需要去声明namcespace或 package。
5. 注意要点
   1. 同一个文件夹下的所有源码文件(子文件夹除外)，package名字可以随意命名，但必须声明相同的package，否则会报错。
   2. 在同一个目录下所有文件中的代码是透明可以互相访问的，但是跨包文件夹的必须通过包名访问


## modules机制

go modules包管理机制： 可以添加依赖，删除未使用的依赖项


1. Go 1.11 引入了 Go modules 机制，用于管理依赖包。常用命令如下：
   1. `go mod`: 查看可用命令help信息
   1. `go mod init`： 初始化项目，创建 go.mod 文件
   2. 确保你的项目 有 go.mod：`go mod init 你的项目名`
   3. 安装第三方包：`go get github.com/gin-gonic/gin`
   4. 安装完后，2 个变化
      1. go.mod 文件：自动记录了项目依赖的包及版本信息
      2. go.sum 文件：记录了项目依赖的包及版本信息
      3. 包自动下载到`GOPATH/pkg/mod`
   5. 根据 go.mod 自动初始化下载所有包:   `go mod tidy` 
      1. 会依据【代码里的 import】进行扫描，但必须有 【go.mod】 作为模块环境进行必要条件。
         1. 没有 go.mod 文件，直接报错，无法工作，必须先 go mod init。
      2. 它会自动扫描你代码里 import 的包，然后去修正gomod文件：
         1. 添加缺少的依赖
         2. 删除没用的依赖
         3. 自动下载所有需要的包
         4. 自动修正版本、清理冲突
   6. 清空本地缓存（重装用）：`go clean -modcache`
   7. `go list -m all` : 查看当前项目下的管理的所有依赖包
   8. `go list -m -versions github.com/gin-gonic/gin` : 查看包所有可用的版本
   9.  `go get github.com/gin-gonic/gin@v1.7.7`: 安装或升级到指定版本
   10. 升级某个包：使用 `go get -u` 升级到最新的修订版本和次要版本
       1.  `go get -u=patch`：升级到最新的修订版本
   11. 降级某个包：使用 `go get -d` 降级
   12. `go mod edit -replace github.com/gin-gonic/gin=github.com/gin-gonic/gin@v1.7.7`：将包替换为指定版本
       1.  执行后，会在go.mod 文件中添加一条替换规则
       2.  意思就是项目代码中引用的包名是左边的gin，但是实际下载或编译使用的是右边的v1.7.7版本，就是做个中转映射
       3.  正常稳定开发 不需要用，只有版本不兼容、bug、调试、冲突时才用
   13. `go mod tidy -v`：查看依赖包的版本
   14. `go mod graph`：查看依赖包的依赖关系
   15. `go mod vendor`：将依赖包复制到 vendor 目录
   16. `go mod vendor -v`：查看依赖包的版本`

2. Go modules 下载的包 默认装在哪里？
   1. 所有第三方包，默认统一安装在：`$GOPATH/pkg/mod`
      1. $GOPATH: Go的安装路径 自动自带系统路径变量
   2. Windows 默认路径: `C:\Users\你的用户名\go\pkg\mod`
   3. Linux/Mac 默认路径: `/Users/你的用户名/go/pkg/mod`
3. 一般在import三方依赖时飘红的编辑器的提示，点击快捷安装按钮即可或者自己使用命令
   1. `go get 包地址 ` 如`go get github.com/gin-gonic/gin`
4. go modules下载代理设置国内镜像源
   1. `go env` 进入环境变量展示
   2. 设置环境变量的字段
      1. `go env -w GO111MODULE=on`  // 开启 / 关闭 Go Modules 模式,空值 / on = 开启 modules（现在默认就是）
      2. `go env -w GOPROXY=https://goproxy.cn,direct` // 设置国内代理


## 编码规范

### 代码规范

1. 命名规范
   1. 包名：尽量和目录保持一致，尽量采取有意义且简短的包名，如：user、blog、product，不要和标准库重名，包名称一般全部使用小写字母，多个单词：全小写，多词直接连，不下划线、不驼峰，因为要短简洁
   2. 文件名：使用小写字母，多个单词用下划线连接，如：user_info.go，这叫蛇形命名法
   3. 变量名：用小驼峰
      1. 蛇形：python，php
      2. 驼峰：在java，c。go，js中，变量名一般使用小驼峰命名法，如：userInfo
         1. java中一般会使用全单词，go中一般使用首字母缩写，推崇简短
      3. 函数方法：也用小驼峰，要求公共导出的话，可以使用大驼峰
   4. 结构体命名：也是小驼峰，要求公共导出的话，可以使用大驼峰
   5. 接口命名：与结构体一样，单词一般以er结尾，和首字母I 开头代表interface
   6. 常量命名：一般全部大写，多个单词使用下划线
2. 最核心的 Go 公开 / 私有规则（超级重要）
   1. 首字母大写 = 公开（public）
   2. 首字母小写 = 私有（private）
   3. 这是 Go 唯一的访问控制方式！

### 注释规范

---- 可以装下自动注释插件

1. go中有2种注释
   1. 单行 // 
   2. 大段注释 /* */
2. 变量加注释：变量上方加
3. 结构体上方加注释：// 结构体名字：描述
4. 包注释
   1. 上方加注释
   2. ```
   // user 包：描述
   // author：内容
   // datetime：时间
   ```
6. 接口注释
7. 函数注释：上方加函数名字：描述
   1. 函数参数注释：参数名：参数描述
   2. 返回值
8. 代码逻辑的注释

### import规范

1. import3类包：
   1. go自带的包
   2. 第三方包
   3. 自定义包
2. 这几个引入时，每一类需要分组引入，用空白行分组隔开
