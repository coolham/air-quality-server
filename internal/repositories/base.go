package repositories

import (
	"context"
	"fmt"
	"reflect"

	"air-quality-server/internal/utils"

	"gorm.io/gorm"
)

// BaseRepository 基础仓储接口
type BaseRepository[T any] interface {
	Create(ctx context.Context, entity *T) error
	GetByID(ctx context.Context, id interface{}) (*T, error)
	Update(ctx context.Context, id interface{}, updates interface{}) error
	Delete(ctx context.Context, id interface{}) error
	List(ctx context.Context, req *ListRequest) (*ListResponse[T], error)
	Count(ctx context.Context, conditions map[string]interface{}) (int64, error)
}

// ListRequest 列表请求
type ListRequest struct {
	Page       int                    `json:"page"`
	PageSize   int                    `json:"page_size"`
	OrderBy    string                 `json:"order_by"`
	Order      string                 `json:"order"` // asc, desc
	Conditions map[string]interface{} `json:"conditions"`
}

// ListResponse 列表响应
type ListResponse[T any] struct {
	Data     []T   `json:"data"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
}

// baseRepository 基础仓储实现
type baseRepository[T any] struct {
	db     *gorm.DB
	logger utils.Logger
}

// NewBaseRepository 创建基础仓储
func NewBaseRepository[T any](db *gorm.DB, logger utils.Logger) BaseRepository[T] {
	return &baseRepository[T]{
		db:     db,
		logger: logger,
	}
}

// Create 创建实体
func (r *baseRepository[T]) Create(ctx context.Context, entity *T) error {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		r.logger.Error("创建实体失败", utils.ErrorField(err))
		return fmt.Errorf("创建实体失败: %w", err)
	}
	return nil
}

// GetByID 根据ID获取实体
func (r *baseRepository[T]) GetByID(ctx context.Context, id interface{}) (*T, error) {
	var entity T
	if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("实体不存在")
		}
		r.logger.Error("获取实体失败", utils.ErrorField(err))
		return nil, fmt.Errorf("获取实体失败: %w", err)
	}
	return &entity, nil
}

// Update 更新实体
func (r *baseRepository[T]) Update(ctx context.Context, id interface{}, updates interface{}) error {
	var entity T
	if err := r.db.WithContext(ctx).Model(&entity).Where("id = ?", id).Updates(updates).Error; err != nil {
		r.logger.Error("更新实体失败", utils.ErrorField(err))
		return fmt.Errorf("更新实体失败: %w", err)
	}
	return nil
}

// Delete 删除实体
func (r *baseRepository[T]) Delete(ctx context.Context, id interface{}) error {
	var entity T
	if err := r.db.WithContext(ctx).Delete(&entity, id).Error; err != nil {
		r.logger.Error("删除实体失败", utils.ErrorField(err))
		return fmt.Errorf("删除实体失败: %w", err)
	}
	return nil
}

// List 获取实体列表
func (r *baseRepository[T]) List(ctx context.Context, req *ListRequest) (*ListResponse[T], error) {
	var entities []T
	var total int64

	// 构建查询
	query := r.db.WithContext(ctx).Model(new(T))

	// 添加条件
	if req.Conditions != nil {
		for key, value := range req.Conditions {
			if value != nil {
				query = query.Where(fmt.Sprintf("%s = ?", key), value)
			}
		}
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("计算总数失败", utils.ErrorField(err))
		return nil, fmt.Errorf("计算总数失败: %w", err)
	}

	// 设置排序
	if req.OrderBy != "" {
		order := "ASC"
		if req.Order == "desc" {
			order = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", req.OrderBy, order))
	}

	// 设置分页
	if req.Page > 0 && req.PageSize > 0 {
		offset := (req.Page - 1) * req.PageSize
		query = query.Offset(offset).Limit(req.PageSize)
	}

	// 执行查询
	if err := query.Find(&entities).Error; err != nil {
		r.logger.Error("查询实体列表失败", utils.ErrorField(err))
		return nil, fmt.Errorf("查询实体列表失败: %w", err)
	}

	return &ListResponse[T]{
		Data:     entities,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// Count 统计实体数量
func (r *baseRepository[T]) Count(ctx context.Context, conditions map[string]interface{}) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(new(T))

	// 添加条件
	if conditions != nil {
		for key, value := range conditions {
			if value != nil {
				query = query.Where(fmt.Sprintf("%s = ?", key), value)
			}
		}
	}

	if err := query.Count(&count).Error; err != nil {
		r.logger.Error("统计实体数量失败", utils.ErrorField(err))
		return 0, fmt.Errorf("统计实体数量失败: %w", err)
	}

	return count, nil
}

// Transaction 事务执行
func (r *baseRepository[T]) Transaction(ctx context.Context, fn func(*gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(fn)
}

// GetTableName 获取表名
func (r *baseRepository[T]) GetTableName() string {
	var entity T
	entityType := reflect.TypeOf(entity)
	if entityType.Kind() == reflect.Ptr {
		entityType = entityType.Elem()
	}

	// 尝试获取TableName方法
	if method, ok := entityType.MethodByName("TableName"); ok {
		entityValue := reflect.New(entityType)
		result := method.Func.Call([]reflect.Value{entityValue})
		if len(result) > 0 {
			return result[0].String()
		}
	}

	// 默认使用类型名的小写形式
	return fmt.Sprintf("%ss", entityType.Name())
}
