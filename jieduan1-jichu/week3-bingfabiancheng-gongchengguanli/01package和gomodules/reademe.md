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
   1. 确保你的项目 有 go.mod：`go mod init 你的项目名`
   2. 安装第三方包：`go get github.com/gin-gonic/gin`
   3. 安装完后，2 个变化
      1. go.mod 文件：自动记录了项目依赖的包及版本信息
      2. go.sum 文件：记录了项目依赖的包及版本信息
      3. 包自动下载到`GOPATH/pkg/mod`
   4. 根据 go.mod 自动初始化下载所有包:   `go mod tidy` 
      1. 最常用！别人的项目拷过来，直接跑这个，所有依赖自动装齐
   5. 清空本地缓存（重装用）：`go clean -modcache`
   6. `go list -m all` : 查看当前项目下的管理的所有依赖包
   7. `go list -m -versions github.com/gin-gonic/gin` : 查看包所有可用的版本
   8. `go get github.com/gin-gonic/gin@v1.7.7`: 安装或升级到指定版本

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
