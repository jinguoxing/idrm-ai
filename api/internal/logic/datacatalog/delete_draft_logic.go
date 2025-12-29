package datacatalog

import (
	"context"
	"fmt"

	"idrm/api/internal/svc"
	"idrm/api/internal/types"
	"idrm/model/data-catalog/data_catalog"

	"github.com/zeromicro/go-zero/core/logx"
)

// DeleteDraftLogic 删除草稿
type DeleteDraftLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewDeleteDraftLogic 创建DeleteDraftLogic实例
func NewDeleteDraftLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDraftLogic {
	return &DeleteDraftLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// DeleteDraft 删除草稿
// 注意：只能删除草稿状态的目录，已发布的目录不能删除
func (l *DeleteDraftLogic) DeleteDraft(req *types.DeleteDraftReq) (resp *types.DeleteDraftResp, err error) {
	// 查询目录
	catalog, err := l.svcCtx.DataCatalogModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("查询目录失败: %w", err)
	}

	// 检查状态：已发布的目录不能删除
	if catalog.PublishStatus == data_catalog.StatusPublished {
		return nil, fmt.Errorf("已发布的目录不能删除")
	}

	// 使用事务删除目录及关联数据
	err = l.svcCtx.DataCatalogModel.Trans(l.ctx, func(ctx context.Context, model data_catalog.Model) error {
		// 1. 删除主目录记录
		if err := model.Delete(ctx, req.Id); err != nil {
			return fmt.Errorf("删除目录失败: %w", err)
		}

		// 2. 删除信息项
		if err := l.svcCtx.DataCatalogColumnModel.DeleteByCatalogId(ctx, req.Id); err != nil {
			return fmt.Errorf("删除信息项失败: %w", err)
		}

		// 3. 删除分类关联
		if err := l.svcCtx.DataCatalogCategoryModel.DeleteByCatalogId(ctx, req.Id); err != nil {
			return fmt.Errorf("删除分类关联失败: %w", err)
		}

		// 4. 删除挂接资源
		if err := l.svcCtx.DataResourceModel.DeleteByCatalogId(ctx, req.Id); err != nil {
			return fmt.Errorf("删除挂接资源失败: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &types.DeleteDraftResp{
		Id: req.Id,
	}, nil
}
