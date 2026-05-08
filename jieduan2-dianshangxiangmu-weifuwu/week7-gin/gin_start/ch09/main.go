package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	// 为什么我们通过goland运行main.go的时候并没有生成main.exe文件
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0])) // 获取当前执行文件的绝对路径，os.Args[0]是当前执行文件的路径，filepath.Dir获取路径的目录部分，filepath.Abs获取绝对路径
	fmt.Println(dir)                                 // 结果是 临时的文件夹路径下有个临时的exe文件：C:\Users\jinyan1\AppData\Local\Temp\go-build1234567890\exe\main.exe

	// 手动指定静态文件目录，默认是./static，
	router.Static("/static", "./static")
	// 01- LoadHTMLGlob：---- 后续的返回html模版时，默认就会在这个目录下找
	router.LoadHTMLGlob("templates/*")    // 可以写2行，找1级目录下的文件
	router.LoadHTMLGlob("templates/**/*") // 可以模式匹配该二级目录下所有文件，不会找1级目录下的文件
	// 02- LoadHTMLFiles：会将指定的目录下的文件加载好， 相对目录或写绝对路径--------
	// 坑点 : 这种方式写相对路径，通过go run main.go 运行时，路径是相对于当前执行文件的路径（run对应临时文件夹路径）来找的，所以会找不到文件，必须写绝对路径才能找到或者手动build手动在当前目录下执行才能找到模版文件
	// 加载多个文件
	//router.LoadHTMLFiles("templates/index.tmpl", "templates/goods.html")

	// HTML中直接写文件地址，它特点只看文件名是否相同不看所属的父级目录是否相同，相同就会冲突，只会加载到第一个
	// c.HTML(http.StatusOK, "goods/list.html" 和 c.HTML(http.StatusOK, "users/list.html"下 就会冲突，同名文件list.html
	// 为了防止冲突：：：：必须得在模版中定义下名字：{{define "goods/list.html"}}
	// 如果没有在模板中使用define定义 那么我们就可以使用默认的文件名来找
	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "myindex.tmpl", gin.H{ // 直接写文件名，默认在上面指定的目录文件下找
			"title": "慕课网",
		})
	})

	router.GET("/goods/list", func(c *gin.Context) {
		c.HTML(http.StatusOK, "goods/list.html", gin.H{
			"title": "慕课网",
		})
	})

	router.GET("/users/list", func(c *gin.Context) {
		c.HTML(http.StatusOK, "users/list.html", gin.H{
			"title": "慕课网",
		})
	})

	router.GET("/goods", func(c *gin.Context) {
		c.HTML(http.StatusOK, "goods.html", gin.H{
			"name": "微服务开发",
		})
	})

	router.Run(":8083")
}
