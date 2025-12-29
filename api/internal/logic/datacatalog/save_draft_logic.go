package datacatalog

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"idrm/api/internal/svc"
	"idrm/api/internal/types"
	"idrm/model/data-catalog/data_catalog"
	"idrm/model/data-catalog/data_catalog_category"
	"idrm/model/data-catalog/data_catalog_column"
	"idrm/model/data-catalog/data_resource"

	"github.com/zeromicro/go-zero/core/logx"
)

// SaveDraftLogic 暂存数据资源目录
type SaveDraftLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewSaveDraftLogic 创建SaveDraftLogic实例
func NewSaveDraftLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveDraftLogic {
	return &SaveDraftLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// SaveDraft 暂存数据资源目录
// 支持三种场景：
// 1. 创建暂存：首次创建全新的目录草稿（CatalogID 为空）
// 2. 变更暂存：对已发布目录进行编辑，创建或更新草稿副本
// 3. 草稿更新：对现有草稿进行内容更新
func (l *SaveDraftLogic) SaveDraft(req *types.SaveDraftReq) (resp *types.SaveDraftResp, err error) {
	// ========== 1. 参数验证 ==========
	if err := l.validateRequest(req); err != nil {
		return nil, fmt.Errorf("参数验证失败: %w", err)
	}

	// ========== 2. 业务逻辑处理 ==========
	var catalogId string

	if req.CatalogId == "" {
		// 场景1: 创建新草稿
		catalogId, err = l.createNewDraft(req)
	} else {
		// 查询原目录状态
		existingCatalog, err := l.svcCtx.DataCatalogModel.FindOne(l.ctx, req.CatalogId)
		if err != nil {
			return nil, fmt.Errorf("查询目录失败: %w", err)
		}

		if existingCatalog.PublishStatus == data_catalog.StatusPublished {
			// 场景2: 已发布目录的变更暂存
			catalogId, err = l.createChangeDraft(req, existingCatalog)
		} else {
			// 场景3: 更新现有草稿
			catalogId, err = l.updateExistingDraft(req, existingCatalog)
		}
	}

	if err != nil {
		return nil, err
	}

	// ========== 3. 构建响应 ==========
	return &types.SaveDraftResp{
		CatalogId: catalogId,
		Status:    string(data_catalog.StatusUnpublished),
	}, nil
}

// validateRequest 验证请求参数
func (l *SaveDraftLogic) validateRequest(req *types.SaveDraftReq) error {
	if req.Title == "" {
		return data_catalog.NewInvalidParameterError("目录名称")
	}

	if req.SourceDeptId == "" {
		return data_catalog.NewInvalidParameterError("来源部门")
	}

	// 验证逻辑视图（最多只能有一个）
	logicalViewCount := 0
	for _, resource := range req.MountResources {
		if resource.ResourceType == 1 { // 1 表示逻辑视图
			logicalViewCount++
		}
	}

	if logicalViewCount > 1 {
		return data_resource.NewDuplicateLogicalViewError()
	}

	return nil
}

// createNewDraft 创建新草稿
func (l *SaveDraftLogic) createNewDraft(req *types.SaveDraftReq) (string, error) {
	// 检查名称唯一性
	exists, err := l.svcCtx.DataCatalogModel.CheckNameExists(l.ctx, req.SourceDeptId, req.Title, "")
	if err != nil {
		return "", fmt.Errorf("检查名称唯一性失败: %w", err)
	}
	if exists {
		return "", data_catalog.NewDuplicateNameError(req.Title)
	}

	// 生成新的目录ID
	catalogId := l.generateCatalogId()

	// 使用事务创建目录及关联数据
	err = l.svcCtx.DataCatalogModel.Trans(l.ctx, func(ctx context.Context, model data_catalog.Model) error {
		// 1. 创建主目录记录
		catalog := l.buildCatalogFromRequest(req, catalogId, "")
		if _, err := model.Insert(ctx, catalog); err != nil {
			return fmt.Errorf("创建目录失败: %w", err)
		}

		// 2. 批量插入信息项
		if len(req.Columns) > 0 {
			columns := l.buildColumnsFromRequest(req, catalogId)
			if err := l.svcCtx.DataCatalogColumnModel.Insert(ctx, columns); err != nil {
				return fmt.Errorf("创建信息项失败: %w", err)
			}
		}

		// 3. 批量插入分类关联
		if len(req.CategoryNodeIds) > 0 {
			categories := l.buildCategoriesFromRequest(req, catalogId)
			if err := l.svcCtx.DataCatalogCategoryModel.Insert(ctx, categories); err != nil {
				return fmt.Errorf("创建分类关联失败: %w", err)
			}
		}

		// 4. 批量插入挂接资源
		if len(req.MountResources) > 0 {
			resources := l.buildResourcesFromRequest(req, catalogId)
			for _, resource := range resources {
				if err := l.svcCtx.DataResourceModel.Insert(ctx, resource); err != nil {
					return fmt.Errorf("创建挂接资源失败: %w", err)
				}
			}
		}

		return nil
	})

	return catalogId, err
}

// createChangeDraft 创建变更草稿（已发布目录）
func (l *SaveDraftLogic) createChangeDraft(req *types.SaveDraftReq, existingCatalog *data_catalog.DataCatalog) (string, error) {
	// 检查是否已存在草稿副本
	if existingCatalog.DraftId > 0 {
		// 更新现有草稿副本
		draftCatalog, err := l.svcCtx.DataCatalogModel.FindOne(l.ctx, strconv.FormatUint(existingCatalog.DraftId, 10))
		if err != nil {
			return "", fmt.Errorf("查询草稿副本失败: %w", err)
		}

		return l.updateDraftContent(req, draftCatalog, existingCatalog.Id)
	}

	// 创建新的草稿副本
	draftId := l.generateCatalogId()

	err := l.svcCtx.DataCatalogModel.Trans(l.ctx, func(ctx context.Context, model data_catalog.Model) error {
		// 1. 创建草稿副本
		draftCatalog := l.buildDraftCopyFromRequest(req, draftId, existingCatalog.Id)
		if _, err := model.Insert(ctx, draftCatalog); err != nil {
			return fmt.Errorf("创建草稿副本失败: %w", err)
		}

		// 2. 更新原目录的 DraftId 字段
		draftIdUint, _ := strconv.ParseUint(draftId, 10, 64)
		existingCatalog.DraftId = draftIdUint
		if err := model.Update(ctx, existingCatalog); err != nil {
			return fmt.Errorf("更新原目录失败: %w", err)
		}

		// 3. 删除旧的草稿关联数据
		if err := l.svcCtx.DataCatalogColumnModel.DeleteByCatalogId(ctx, draftId); err != nil {
			return fmt.Errorf("删除旧信息项失败: %w", err)
		}
		if err := l.svcCtx.DataCatalogCategoryModel.DeleteByCatalogId(ctx, draftId); err != nil {
			return fmt.Errorf("删除旧分类关联失败: %w", err)
		}
		if err := l.svcCtx.DataResourceModel.DeleteByCatalogId(ctx, draftId); err != nil {
			return fmt.Errorf("删除旧挂接资源失败: %w", err)
		}

		// 4. 创建新的关联数据
		if len(req.Columns) > 0 {
			columns := l.buildColumnsFromRequest(req, draftId)
			if err := l.svcCtx.DataCatalogColumnModel.Insert(ctx, columns); err != nil {
				return fmt.Errorf("创建信息项失败: %w", err)
			}
		}

		if len(req.CategoryNodeIds) > 0 {
			categories := l.buildCategoriesFromRequest(req, draftId)
			if err := l.svcCtx.DataCatalogCategoryModel.Insert(ctx, categories); err != nil {
				return fmt.Errorf("创建分类关联失败: %w", err)
			}
		}

		if len(req.MountResources) > 0 {
			resources := l.buildResourcesFromRequest(req, draftId)
			for _, resource := range resources {
				if err := l.svcCtx.DataResourceModel.Insert(ctx, resource); err != nil {
					return fmt.Errorf("创建挂接资源失败: %w", err)
				}
			}
		}

		return nil
	})

	return draftId, err
}

// updateExistingDraft 更新现有草稿
func (l *SaveDraftLogic) updateExistingDraft(req *types.SaveDraftReq, existingCatalog *data_catalog.DataCatalog) (string, error) {
	// 检查名称唯一性（排除自己）
	exists, err := l.svcCtx.DataCatalogModel.CheckNameExists(l.ctx, req.SourceDeptId, req.Title, existingCatalog.Id)
	if err != nil {
		return "", fmt.Errorf("检查名称唯一性失败: %w", err)
	}
	if exists {
		return "", data_catalog.NewDuplicateNameError(req.Title)
	}

	return l.updateDraftContent(req, existingCatalog, "")
}

// updateDraftContent 更新草稿内容
func (l *SaveDraftLogic) updateDraftContent(req *types.SaveDraftReq, draftCatalog *data_catalog.DataCatalog, parentCatalogId string) (string, error) {
	err := l.svcCtx.DataCatalogModel.Trans(l.ctx, func(ctx context.Context, model data_catalog.Model) error {
		// 1. 更新主目录记录
		updatedCatalog := l.buildCatalogFromRequest(req, draftCatalog.Id, parentCatalogId)
		if err := model.Update(ctx, updatedCatalog); err != nil {
			return fmt.Errorf("更新目录失败: %w", err)
		}

		// 2. 删除并重新创建关联数据
		catalogId := draftCatalog.Id

		if err := l.svcCtx.DataCatalogColumnModel.DeleteByCatalogId(ctx, catalogId); err != nil {
			return fmt.Errorf("删除旧信息项失败: %w", err)
		}

		if err := l.svcCtx.DataCatalogCategoryModel.DeleteByCatalogId(ctx, catalogId); err != nil {
			return fmt.Errorf("删除旧分类关联失败: %w", err)
		}

		if err := l.svcCtx.DataResourceModel.DeleteByCatalogId(ctx, catalogId); err != nil {
			return fmt.Errorf("删除旧挂接资源失败: %w", err)
		}

		// 3. 创建新的关联数据
		if len(req.Columns) > 0 {
			columns := l.buildColumnsFromRequest(req, catalogId)
			if err := l.svcCtx.DataCatalogColumnModel.Insert(ctx, columns); err != nil {
				return fmt.Errorf("创建信息项失败: %w", err)
			}
		}

		if len(req.CategoryNodeIds) > 0 {
			categories := l.buildCategoriesFromRequest(req, catalogId)
			if err := l.svcCtx.DataCatalogCategoryModel.Insert(ctx, categories); err != nil {
				return fmt.Errorf("创建分类关联失败: %w", err)
			}
		}

		if len(req.MountResources) > 0 {
			resources := l.buildResourcesFromRequest(req, catalogId)
			for _, resource := range resources {
				if err := l.svcCtx.DataResourceModel.Insert(ctx, resource); err != nil {
					return fmt.Errorf("创建挂接资源失败: %w", err)
				}
			}
		}

		return nil
	})

	return draftCatalog.Id, err
}

// ========== 辅助方法 ==========

// generateCatalogId 生成目录ID
func (l *SaveDraftLogic) generateCatalogId() string {
	// TODO: 使用雪花算法生成ID
	return fmt.Sprintf("catalog_%d", time.Now().UnixNano())
}

// buildCatalogFromRequest 从请求构建目录对象
func (l *SaveDraftLogic) buildCatalogFromRequest(req *types.SaveDraftReq, catalogId string, parentCatalogId string) *data_catalog.DataCatalog {
	now := time.Now()
	return &data_catalog.DataCatalog{
		Id:              catalogId,
		Title:           req.Title,
		Description:     req.Description,
		SourceDeptId:    req.SourceDeptId,
		DataDomain:      req.DataDomain,
		PublishStatus:   data_catalog.StatusUnpublished,
		OnlineStatus:    data_catalog.OnlineStatusNotline,
		SharedType:      req.SharedType,
		SharedCondition: req.SharedCondition,
		OpenType:        req.OpenType,
		OpenCondition:   req.OpenCondition,
		UpdateCycle:     req.UpdateCycle,
		OtherUpdateCycle: req.OtherUpdateCycle,
		SyncMechanism:   req.SyncMechanism,
		SyncFrequency:   req.SyncFrequency,
		DataKind:        1, // 默认值
		SharedMode:      1, // 默认值
		ColumnUnshared:  false,
		Source:          1, // 默认值
		CurrentVersion:  true,
		OwnerId:         "system", // 默认值
		OwnerName:       "系统管理员",
		CreatedAt:       now,
	}
}

// buildDraftCopyFromRequest 从请求构建草稿副本对象
func (l *SaveDraftLogic) buildDraftCopyFromRequest(req *types.SaveDraftReq, draftId string, parentCatalogId string) *data_catalog.DataCatalog {
	catalog := l.buildCatalogFromRequest(req, draftId, parentCatalogId)
	// 草稿副本不需要特殊标识，通过 DraftId 关联即可
	return catalog
}

// buildColumnsFromRequest 从请求构建信息项列表
func (l *SaveDraftLogic) buildColumnsFromRequest(req *types.SaveDraftReq, catalogId string) []*data_catalog_column.DataCatalogColumn {
	columns := make([]*data_catalog_column.DataCatalogColumn, 0, len(req.Columns))
	now := time.Now()
	for i, col := range req.Columns {
		// 将 int8 类型的 DataFormat 转换为字符串
		columnType := "string" // 默认值
		if col.DataFormat > 0 {
			// 简化处理：根据 DataFormat 映射到类型字符串
			typeMap := map[int8]string{
				1: "string",
				2: "integer",
				3: "decimal",
				4: "date",
				5: "datetime",
			}
			if t, ok := typeMap[col.DataFormat]; ok {
				columnType = t
			}
		}

		columns = append(columns, &data_catalog_column.DataCatalogColumn{
			Id:             fmt.Sprintf("col_%d", time.Now().UnixNano()+int64(i)),
			CatalogId:      catalogId,
			ColumnName:     col.TechnicalName,
			ColumnAlias:    col.BusinessName,
			ColumnType:     columnType,
			ColumnLength:   col.DataLength,
			ColumnScale:    int(col.DataPrecision),
			ColumnDescribe: col.Description,
			IsPrimaryKey:   col.PrimaryFlag,
			IsRequired:     !col.NullFlag,
			SortOrder:      col.Index,
			CreatedAt:      now,
		})
	}
	return columns
}

// buildCategoriesFromRequest 从请求构建分类关联列表
func (l *SaveDraftLogic) buildCategoriesFromRequest(req *types.SaveDraftReq, catalogId string) []*data_catalog_category.DataCatalogCategory {
	categories := make([]*data_catalog_category.DataCatalogCategory, 0, len(req.CategoryNodeIds))
	now := time.Now()
	for i, catId := range req.CategoryNodeIds {
		categories = append(categories, &data_catalog_category.DataCatalogCategory{
			Id:         fmt.Sprintf("cat_%d", time.Now().UnixNano()+int64(i)),
			CatalogId:  catalogId,
			CategoryId: catId,
			CreatedAt:  now,
		})
	}
	return categories
}

// buildResourcesFromRequest 从请求构建资源列表
func (l *SaveDraftLogic) buildResourcesFromRequest(req *types.SaveDraftReq, catalogId string) []*data_resource.DataResource {
	resources := make([]*data_resource.DataResource, 0, len(req.MountResources))
	now := time.Now()
	for i, res := range req.MountResources {
		// 将 int8 类型的 ResourceType 转换为 ResourceType 枚举
		var resourceType data_resource.ResourceType
		switch res.ResourceType {
		case 1:
			resourceType = data_resource.ResourceTypeLogicalView
		case 2:
			resourceType = data_resource.ResourceTypeApiInterface
		case 3:
			resourceType = data_resource.ResourceTypeFile
		default:
			resourceType = data_resource.ResourceTypeDatabaseTable
		}

		resources = append(resources, &data_resource.DataResource{
			Id:           fmt.Sprintf("res_%d", time.Now().UnixNano()+int64(i)),
			CatalogId:    catalogId,
			ResourceType: resourceType,
			ResourceName: res.Name,
			SortOrder:    i,
			Source:       "manual",
			CreatedAt:    now,
		})
	}
	return resources
}
