package ginkgo

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2" // 猴子补丁（万能打桩，兜底 Mock）
	"github.com/golang/mock/gomock"      // 官方接口 Mock 框架（主力 Mock）
	. "github.com/onsi/ginkgo"           // BDD 测试框架（Go 版 Jest/Mocha）
	. "github.com/onsi/gomega"           // BDD 断言库（Go 版 expect/chai）
	"github.com/stretchr/testify/assert" // 通用断言库（原生测试最常用）
)

// Ginkgo 提供 GinkgoT() 返回一个兼容 *testing.T 的对象，传给 assert 即可。
// BDD 断言库和通用断言库是功能类似，也可以混用，用GinkgoT()就可以桥接

// 这是Ginkgo 测试的入口函数，一个文件（测试套件）只出现一次！，告诉 Go 原生测试框架：“把执行权交给 Ginkgo，我这个文件用 BDD 跑。”
func TestBooks(t *testing.T) {
	RegisterFailHandler(Fail)  // 把 Ginkgo 的失败 交给 Go 原生测试系统处理
	RunSpecs(t, "Books Suite") // 运行当前文件里所有的 BDD 测试用例。t：Go 原生测试对象，"Books Suite"：测试套件名字（随便写，方便看日志）
}

// 并不一定需要每个测试用例都这么写， 对于核心的函数或者核心的业务逻辑我们建议设计好的测试用例
// 定义一个这是一个测试集合，测 Books 模块
var _ = Describe("Books", func() {
	var (
		longBook  string
		shortBook string

		pathches *gomonkey.Patches
		ctl      *gomock.Controller
	)
	// 每次用例前执行（核心：复用！）
	BeforeEach(func() {
		longBook = "long"
		shortBook = "short"

		ctl = gomock.NewController(GinkgoT())
	})
	// 每次用例后执行（清理）
	AfterEach(func() {
		longBook = ""
		shortBook = ""

		ctl.Finish()
		pathches.Reset()
	})
	// 子模块分组（清晰！）
	Describe("Add Books", func() {
		It("should be able to add a book", func() {
			//调用AddBook方法，并传入参数，期望返回的结果为true
			assert.Equal(GinkgoT(), "long", longBook)
		})
		It("should not be able to add a book", func() {
			//调用AddBook方法，并传入参数，期望返回的结果为true
			assert.Equal(GinkgoT(), "short", shortBook)
		})
	})

	Describe("Delete Books", func() {

	})
})

/*
1. proto文件可以用作http和rpc服务的生成标准写法
	我写了一个gin的服务，我还要手动去维护api文档，手动去yapi上维护
	当有了proto后，可以后期维护和迭代很简单， 改了任何代码你都可以直接生成api
	可以直接将proto生成swagger文件，然后一键导入到yapi上，这样就可以直接在yapi上查看api文档了
2. 在kratos中对proto的依赖更加重， 可以用来定义一些错误码， 并生成go源码直接使用
3. kratos甚至将配置文件都给你映射成proto文件
业内很多框架都开始逐步接受将proto文件作为核心的标准去写一系列插件去自动生成代码
proto validate

go-zero更溜，goctl，保姆式的框架 api文件 go-zero和kratos的一套设计理念
*/
