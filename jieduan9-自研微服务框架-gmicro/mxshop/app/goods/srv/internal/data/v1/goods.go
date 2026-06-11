package v1

import (
	"context"
	"mxshop/app/goods/srv/internal/domain/do"
	metav1 "mxshop/pkg/common/meta/v1"

	"gorm.io/gorm"
)

type GoodsStore interface {
	Get(ctx context.Context, ID uint64) (*do.GoodsDO, error)
	ListByIDs(ctx context.Context, ids []uint64, orderby []string) (*do.GoodsDOList, error)
	List(ctx context.Context, orderby []string, opts metav1.ListMeta) (*do.GoodsDOList, error)
	Create(ctx context.Context, goods *do.GoodsDO) error
	// 因为与es操作需要同步
	// 在事务中去提交创建
	CreateInTxn(ctx context.Context, txn *gorm.DB, goods *do.GoodsDO) error
	Update(ctx context.Context, goods *do.GoodsDO) error
	// 在事务中去提交创建
	UpdateInTxn(ctx context.Context, txn *gorm.DB, goods *do.GoodsDO) error
	Delete(ctx context.Context, ID uint64) error
	// 在事务中去提交创建
	DeleteInTxn(ctx context.Context, txn *gorm.DB, ID uint64) error
	// 获取全局数据库事务对象，给上面的方法提供txn
	Begin() *gorm.DB
}
