package forms

type CategoryForm struct {
	Name           string `form:"name" json:"name" binding:"required,min=3,max=20"`
	ParentCategory int32  `form:"parent" json:"parent"`
	Level          int32  `form:"level" json:"level" binding:"required,oneof=1 2 3"` // oneof 每个空格隔开
	// 这个bool必须用零值
	IsTab *bool `form:"is_tab" json:"is_tab" binding:"required"`
}

type UpdateCategoryForm struct {
	Name string `form:"name" json:"name" binding:"required,min=3,max=20"`
	// 这个bool必须用指针
	IsTab *bool `form:"is_tab" json:"is_tab"`
}
