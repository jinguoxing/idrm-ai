package datacatalog

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"idrm/api/internal/svc"
	"idrm/api/internal/types"
	"idrm/model/data-catalog/data_catalog_category"
	"idrm/model/data-catalog/data_catalog_column"
	"idrm/model/data-catalog/data_resource"

	"github.com/zeromicro/go-zero/core/logx"
)

// GetDraftLogic 获取草稿详情
type GetDraftLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetDraftLogic 创建GetDraftLogic实例
func NewGetDraftLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDraftLogic {
	return &GetDraftLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetDraft 获取草稿详情
func (l *GetDraftLogic) GetDraft(req *types.GetDraftReq) (resp *types.GetDraftResp, err error) {
	// 查询目录
	catalog, err := l.svcCtx.DataCatalogModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("查询目录失败: %w", err)
	}

	// 查询信息项
	columns, err := l.svcCtx.DataCatalogColumnModel.FindByCatalogId(l.ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("查询信息项失败: %w", err)
	}

	// 查询分类关联
	categories, err := l.svcCtx.DataCatalogCategoryModel.FindByCatalogId(l.ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("查询分类关联失败: %w", err)
	}

	// 查询挂接资源
	resources, err := l.svcCtx.DataResourceModel.FindByCatalogId(l.ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("查询挂接资源失败: %w", err)
	}

	// 判断是否为草稿副本（通过 DraftId 字段判断，当前目录有 DraftId 说明它是原目录，有草稿副本）
	// 这里简化处理：如果当前目录的 DraftId > 0，说明存在草稿副本
	draftIdStr := ""
	if catalog.DraftId > 0 {
		draftIdStr = strconv.FormatUint(catalog.DraftId, 10)
	}

	// 构建响应
	return &types.GetDraftResp{
		CatalogId:        catalog.Id,
		Title:            catalog.Title,
		Description:      catalog.Description,
		SourceDeptId:     catalog.SourceDeptId,
		SourceDeptName:   catalog.OwnerName, // 使用 OwnerName 作为部门名称
		DataDomain:       catalog.DataDomain,
		PublishStatus:    string(catalog.PublishStatus),
		DraftId:          draftIdStr,
		IsDraftCopy:      false, // 简化处理
		SharedType:       catalog.SharedType,
		SharedCondition:  catalog.SharedCondition,
		OpenType:         catalog.OpenType,
		OpenCondition:    catalog.OpenCondition,
		Columns:          l.buildColumnItems(columns),
		MountResources:   l.buildMountResources(resources),
		UpdateCycle:      catalog.UpdateCycle,
		OtherUpdateCycle: catalog.OtherUpdateCycle,
		SyncMechanism:    catalog.SyncMechanism,
		SyncFrequency:    catalog.SyncFrequency,
		CategoryNodeIds:  l.buildCategoryNodeIds(categories),
		CreatedAt:        catalog.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        formatTimePtr(catalog.UpdatedAt),
	}, nil
}

func formatTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

func (l *GetDraftLogic) buildColumnItems(columns []*data_catalog_column.DataCatalogColumn) []types.ColumnItem {
	result := make([]types.ColumnItem, 0, len(columns))
	for _, col := range columns {
		// 将字符串类型转换为 int8
		var dataFormat int8
		if col.ColumnType != "" {
			// 简化处理：使用字段类型的哈希值或映射
			dataFormat = int8(col.ColumnType[0])
		}

		result = append(result, types.ColumnItem{
			TechnicalName: col.ColumnName,
			BusinessName:  col.ColumnAlias,
			DataFormat:    dataFormat,
			DataLength:    col.ColumnLength,
			DataPrecision: int8(col.ColumnScale),
			Description:   col.ColumnDescribe,
			PrimaryFlag:   col.IsPrimaryKey,
			NullFlag:      !col.IsRequired,
			Index:         int(col.SortOrder),
		})
	}
	return result
}

func (l *GetDraftLogic) buildCategoryNodeIds(categories []*data_catalog_category.DataCatalogCategory) []string {
	result := make([]string, 0, len(categories))
	for _, cat := range categories {
		result = append(result, cat.CategoryId)
	}
	return result
}

func (l *GetDraftLogic) buildMountResources(resources []*data_resource.DataResource) []types.MountResource {
	result := make([]types.MountResource, 0, len(resources))
	for _, res := range resources {
		// 将 ResourceType 转换为 int8
		var resourceType int8
		switch res.ResourceType {
		case data_resource.ResourceTypeLogicalView:
			resourceType = 1
		case data_resource.ResourceTypeApiInterface:
			resourceType = 2
		case data_resource.ResourceTypeFile:
			resourceType = 3
		default:
			resourceType = 0
		}

		result = append(result, types.MountResource{
			ResourceId:   res.Id,
			ResourceType: resourceType,
			Name:         res.ResourceName,
		})
	}
	return result
}
