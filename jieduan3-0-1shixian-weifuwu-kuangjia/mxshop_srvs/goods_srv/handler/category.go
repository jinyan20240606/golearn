package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"mxshop_srvs/goods_srv/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"mxshop_srvs/goods_srv/global"
	"mxshop_srvs/goods_srv/proto"

	"google.golang.org/protobuf/types/known/emptypb"
)

// //商品分类
func (s *GoodsServer) GetAllCategorysList(context.Context, *emptypb.Empty) (*proto.CategoryListResponse, error) {
	/*
		希望返回如下这种结构给前端，拼装好这种结构我们选择在server层做，不交给web-服务层做，
		因为这种结构用gorm来做是非常简单的，web-服务一般不与数据库交互，web层拿原始结构体拼接是很麻烦的
		[
			{
				"id":xxx,
				"name":"",
				"level":1,
				"is_tab":false,
				"parent":13xxx,
				"sub_category":[
					"id":xxx,
					"name":"",
					"level":1,
					"is_tab":false,
					"sub_category":[]
				]
			}
		]
	*/
	var categorys []model.Category
	// 预加载关联的子分类字段
	// 条件：只查 Level=1 的数据 → 也就是【一级分类】
	// 这一行是核心！预加载（最关键）。Preload("SubCategory") → 加载二级分类，Preload("SubCategory.SubCategory") → 加载二级分类 + 三级分类

	global.DB.Where(&model.Category{Level: 1}).Preload("SubCategory.SubCategory").Find(&categorys)
	b, _ := json.Marshal(&categorys)
	// 可以直接强转的：[] byte ↔ string
	return &proto.CategoryListResponse{JsonData: string(b)}, nil
}

// //获取子分类
func (s *GoodsServer) GetSubCategory(ctx context.Context, req *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	categoryListResponse := proto.SubCategoryListResponse{}

	// 先查该商品分类是否存在
	var category model.Category
	if result := global.DB.First(&category, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	// 然后再查询子分类
	categoryListResponse.Info = &proto.CategoryInfoResponse{
		Id:             category.ID,
		Name:           category.Name,
		Level:          category.Level,
		IsTab:          category.IsTab,
		ParentCategory: category.ParentCategoryID,
	}

	var subCategorys []model.Category
	var subCategoryResponse []*proto.CategoryInfoResponse
	//preloads := "SubCategory"
	//if category.Level == 1 {
	//	preloads = "SubCategory.SubCategory"
	//}
	global.DB.Where(&model.Category{ParentCategoryID: req.Id}).Find(&subCategorys)
	for _, subCategory := range subCategorys {
		subCategoryResponse = append(subCategoryResponse, &proto.CategoryInfoResponse{
			Id:             subCategory.ID,
			Name:           subCategory.Name,
			Level:          subCategory.Level,
			IsTab:          subCategory.IsTab,
			ParentCategory: subCategory.ParentCategoryID,
		})
	}

	categoryListResponse.SubCategorys = subCategoryResponse
	return &categoryListResponse, nil
}
func (s *GoodsServer) CreateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	category := model.Category{}
	cMap := map[string]interface{}{}
	cMap["name"] = req.Name
	cMap["level"] = req.Level
	cMap["is_tab"] = req.IsTab
	if req.Level != 1 {
		//去查询父类目是否存在
		cMap["parent_category_id"] = req.ParentCategory
	}
	tx := global.DB.Model(&model.Category{}).Create(cMap)
	fmt.Println(tx)
	return &proto.CategoryInfoResponse{Id: int32(category.ID)}, nil
}

/*
*
真正完整逻辑：
查出该分类的所有子分类（递归 / 循环）
删除所有子分类
删除该分类下的所有商品
最后删除自己

目前实现是简版

正确的级联删除应该做到：例如1级：手机 (id: 1)

	└── 2级：安卓手机 (id: 2, 父id:1)
	     └── 3级：小米手机 (id: 3, 父id:2)

----
删除 id=1 时，必须自动删除：
安卓手机（id=2）
小米手机（id=3）
该分类下所有商品
*/
func (s *GoodsServer) DeleteCategory(ctx context.Context, req *proto.DeleteCategoryRequest) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.Category{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	return &emptypb.Empty{}, nil
}

func (s *GoodsServer) UpdateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*emptypb.Empty, error) {
	var category model.Category

	if result := global.DB.First(&category, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.ParentCategory != 0 {
		category.ParentCategoryID = req.ParentCategory
	}
	if req.Level != 0 {
		category.Level = req.Level
	}
	if req.IsTab {
		category.IsTab = req.IsTab
	}

	global.DB.Save(&category)

	return &emptypb.Empty{}, nil
}
