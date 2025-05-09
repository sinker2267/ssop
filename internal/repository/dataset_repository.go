package repository

import (
	_ "encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/sinker/ssop/internal/models"
	"gorm.io/gorm"
)

// DatasetRepository 数据集仓库接口
type DatasetRepository interface {
	Create(dataset *models.Dataset) error
	GetByID(id string) (*models.Dataset, error)
	List(page, size int, filters map[string]interface{}) ([]*models.Dataset, int64, error)
	Update(dataset *models.Dataset) error
	Delete(id string) error
	IncrementDownloadCount(id string) error
}

// datasetRepository 数据集仓库实现
type datasetRepository struct {
	db *gorm.DB
}

// NewDatasetRepository 创建数据集仓库
func NewDatasetRepository(db *gorm.DB) DatasetRepository {
	return &datasetRepository{db: db}
}

// Create 创建数据集
func (r *datasetRepository) Create(dataset *models.Dataset) error {
	return r.db.Create(dataset).Error
}

// GetByID 根据ID获取数据集
func (r *datasetRepository) GetByID(id string) (*models.Dataset, error) {
	var dataset models.Dataset
	err := r.db.Where("id = ?", id).First(&dataset).Error
	if err != nil {
		return nil, err
	}
	return &dataset, nil
}

// List 获取数据集列表
func (r *datasetRepository) List(page, size int, filters map[string]interface{}) ([]*models.Dataset, int64, error) {
	var datasets []*models.Dataset
	var total int64

	query := r.db.Model(&models.Dataset{})

	// 应用过滤条件
	if filters != nil {
		// 数据类型过滤
		if typeVal, ok := filters["type"]; ok && typeVal != "" {
			query = query.Where("type = ?", typeVal)
		}

		// 时间范围过滤
		if startDate, ok := filters["startDate"]; ok && startDate != "" {
			query = query.Where("end_time >= ?", startDate)
		}
		if endDate, ok := filters["endDate"]; ok && endDate != "" {
			query = query.Where("start_time <= ?", endDate)
		}

		// 区域过滤
		if region, ok := filters["region"]; ok && region != "" {
			// 区域格式: "minLat,minLng,maxLat,maxLng"
			parts := strings.Split(region.(string), ",")
			if len(parts) == 4 {
				// 区域搜索的逻辑: 两个区域有重叠
				// 使用JSON函数提取边界值并比较
				// 这里简化处理，实际应根据数据库和存储方式具体实现
				query = query.Where("region_bounds LIKE ?", "%"+region.(string)+"%")
			}
		}

		// 关键词搜索
		if keyword, ok := filters["keyword"]; ok && keyword != "" {
			query = query.Where("name LIKE ? OR description LIKE ? OR tags LIKE ?", 
				fmt.Sprintf("%%%s%%", keyword), 
				fmt.Sprintf("%%%s%%", keyword),
				fmt.Sprintf("%%%s%%", keyword))
		}
	}

	// 计算总数
	query.Count(&total)

	// 分页
	if page > 0 && size > 0 {
		offset := (page - 1) * size
		query = query.Offset(offset).Limit(size)
	}

	// 排序(默认按创建时间倒序)
	query = query.Order("created_at DESC")

	// 执行查询
	err := query.Find(&datasets).Error
	if err != nil {
		return nil, 0, err
	}

	return datasets, total, nil
}

// Update 更新数据集
func (r *datasetRepository) Update(dataset *models.Dataset) error {
	return r.db.Save(dataset).Error
}

// Delete 删除数据集
func (r *datasetRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Dataset{}).Error
}

// IncrementDownloadCount 增加下载计数
func (r *datasetRepository) IncrementDownloadCount(id string) error {
	return r.db.Model(&models.Dataset{}).Where("id = ?", id).
		UpdateColumn("download_count", gorm.Expr("download_count + ?", 1)).
		UpdateColumn("updated_at", time.Now()).Error
} 