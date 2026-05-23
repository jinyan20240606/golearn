package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/olivere/elastic/v7"
)

func main() {
	ctx := context.Background()
	// 初始化一个链接
	// SetSniff必须为false，因为：嗅探（Sniff） 是 ES 客户端默认会去自动探测集群所有节点
	// 先用你给的地址连上 ES，立刻去调用 ES 接口：获取集群所有节点的地址。然后自动切换去连集群内部返回的那个 IP
	// 你给的：192.168.0.104:9200 ✅。ES 内部返回的节点地址：127.0.0.1:9200。
	// 你的机器根本访问不到这个内部 IP → 直接报错 连接失败，只有单节点，内部publish_ip默认是127.0.0.，如果正常集群环境，会返回正确的ip，只是单机时的坑
	logger := log.New(os.Stdout, "", log.LstdFlags) // 创建一个日志输出器，负责把内容打印到控制台。
	// os.Stdout输出到 控制台（标准输出），2参：""前缀为空= 日志前面不加任何字符串，log.LstdFlags：标准日志格式 = 日期 + 时间
	client, _ := elastic.NewClient(elastic.SetURL("http://192.168.0.104:9200"),
		elastic.SetSniff(false),
		// 加上这个logger：就能输出每个命令对应的es操作api了
		elastic.SetTraceLog(logger), ///开启 elastic 客户端的请求追踪它会自动打印出：你执行的每一条 ES 操作对应的真实 RESTful API！
	)

	if err != nil {
		log.Fatalf("连接 ES 失败：%v", err)
	}

	// 2. 测试是否真的连通
	info, code, err := client.Ping("http://192.168.0.104:9200").Do(ctx)
	if err != nil {
		log.Fatalf("ping ES 失败：%v", err)
	}

	log.Printf("ES 连接成功！版本：%s，状态码：%d", info.Version.Number, code)

	// 3. 构建 MatchQuery（核心！）
	// 字段：title
	// 内容：慕课网 Go 教程
	matchQuery := elastic.NewMatchQuery("title", "慕课网 Go 教程")

	// 4. 执行查询
	searchResult, err := client.Search().
		Index("your_index_name"). // 换成你的索引名
		Query(matchQuery).
		Do(ctx)

	if err != nil {
		log.Fatal(err)
	}

	// 4. 打印结果
	log.Printf("找到 %d 条数据\n", searchResult.TotalHits())

	// 遍历输出
	for _, hit := range searchResult.Hits.Hits {
		// hit.Source 默认就是字节切片类型，可以直接转化为字符串，也可以用MarshalJSON()一下
		var data map[string]interface{}
		// 01转换成map：把 ES 返回的原始 JSON 数据 → 解析成 Go 语言的 map 对象
		json.Unmarshal(hit.Source, &data)
		log.Println("结果：", data)

		// 02转成json字符串方式：jsonData, err := value.Source.MarshalJSON()
		// 这个jsonData的到是字节切片类型，查看时，需要转成字符串，fmt.Println(string(jsonData))

		// 03转成结构体
		var article struct {
			Title   string `json:"title"`
			Content string `json:"content"`
			Author  string `json:"author"`
		}
		// JSON字节 → 结构体（和map写法一样）
		err := json.Unmarshal(hit.Source, &article)
	}

	type Article struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Author  string `json:"author"`
	}

	// 3. 准备要添加的数据
	article := Article{
		Title:   "Go操作ES教程",
		Content: "使用client.Index添加数据非常简单",
		Author:  "张三",
	}

	// 4. 使用 client.Index() 添加数据到es中
	result, err := client.Index().
		Index("article_index"). // 索引名（库名）
		// 有 .Id("xxx")  →  PUT（新增/覆盖）
		// 无 .Id()       →  POST（自动生成ID新增）
		// Id("100").               // 可选：指定文档ID，不填自动生成
		BodyJson(article).       // 结构体直接放进去
		Do(context.Background()) // 执行

	if err != nil {
		panic(err)
	}
	// 5. 打印结果
	fmt.Printf("添加成功！ID: %s\n", result.Id)
	fmt.Printf("结果: %+v\n", result)

	// 2. 定义 mapping
	const goodsMapping = `
	{
		"settings": {
			"number_of_shards": 1, 
			"number_of_replicas": 0
		},
		"mappings": {
			"properties": {
			"name": {
				"type": "text",
				"analyzer": "ik_max_word"
			},
			"id": {
				"type": "integer"
			}
			}
		}
	}
	`
	// "number_of_shards": 1,     // 主分片1个（单机测试用）
	// "number_of_replicas": 0    // 副本0个（无备份）

	// 3. 创建索引
	indexName := "article_index"

	// 先判断索引是否存在
	exists, _ := client.IndexExists(indexName).Do(context.Background())
	if !exists {
		// 不存在则创建
		_, err := client.CreateIndex(indexName).
			Body(goodsMapping). // 传入 mapping
			Do(context.Background())

		if err != nil {
			panic(err)
		}
		log.Println("索引创建成功 ✅")
	} else {
		log.Println("索引已存在 ✅")
	}

}
