package services

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/sinker/ssop/internal/models"
	"github.com/sinker/ssop/internal/repository"
	"github.com/sinker/ssop/pkg/logger"
	"github.com/sinker/ssop/pkg/utils"
)

// DatasetService 数据集服务接口
type DatasetService interface {
	CreateDataset(dataset *models.Dataset, file io.Reader, filename string) (string, error)
	GetDatasetByID(id string) (*models.Dataset, error)
	GetDatasets(page, size int, filters map[string]interface{}) ([]*models.Dataset, int64, error)
	UpdateDataset(dataset *models.Dataset) error
	DeleteDataset(id string) error
	DownloadDataset(id string) (string, error)
}

// datasetService 数据集服务实现
type datasetService struct {
	datasetRepo repository.DatasetRepository
	storageDir  string // 数据集文件存储目录
}

// NewDatasetService 创建数据集服务
func NewDatasetService(datasetRepo repository.DatasetRepository, storageDir string) DatasetService {
	// 确保存储目录存在
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		logger.Error("Failed to create dataset storage directory", "error", err)
	}

	return &datasetService{
		datasetRepo: datasetRepo,
		storageDir:  storageDir,
	}
}

// CreateDataset 创建数据集
func (s *datasetService) CreateDataset(dataset *models.Dataset, file io.Reader, filename string) (string, error) {
	// 生成唯一ID
	if dataset.ID == "" {
		dataset.ID = utils.GenerateID("ds")
	}

	// 保存文件
	if file != nil {
		// 创建数据集目录
		datasetDir := filepath.Join(s.storageDir, dataset.ID)
		if err := os.MkdirAll(datasetDir, 0755); err != nil {
			return "", fmt.Errorf("failed to create dataset directory: %w", err)
		}

		// 保存文件
		filePath := filepath.Join(datasetDir, filename)
		outFile, err := os.Create(filePath)
		if err != nil {
			return "", fmt.Errorf("failed to create file: %w", err)
		}
		defer outFile.Close()

		// 计算文件大小
		size, err := io.Copy(outFile, file)
		if err != nil {
			return "", fmt.Errorf("failed to save file: %w", err)
		}

		// 更新数据集文件信息
		dataset.FilePath = filePath
		dataset.Size = size
	}

	// 保存数据集信息到数据库
	err := s.datasetRepo.Create(dataset)
	if err != nil {
		return "", fmt.Errorf("failed to save dataset: %w", err)
	}

	return dataset.ID, nil
}

// GetDatasetByID 根据ID获取数据集
func (s *datasetService) GetDatasetByID(id string) (*models.Dataset, error) {
	return s.datasetRepo.GetByID(id)
}

// GetDatasets 获取数据集列表
func (s *datasetService) GetDatasets(page, size int, filters map[string]interface{}) ([]*models.Dataset, int64, error) {
	return s.datasetRepo.List(page, size, filters)
}

// UpdateDataset 更新数据集
func (s *datasetService) UpdateDataset(dataset *models.Dataset) error {
	// 确保数据集存在
	existingDataset, err := s.datasetRepo.GetByID(dataset.ID)
	if err != nil {
		return fmt.Errorf("dataset not found: %w", err)
	}

	// 保留不可修改的字段
	dataset.FilePath = existingDataset.FilePath
	dataset.Size = existingDataset.Size
	dataset.CreatedAt = existingDataset.CreatedAt
	dataset.CreatedBy = existingDataset.CreatedBy
	dataset.DownloadCount = existingDataset.DownloadCount

	return s.datasetRepo.Update(dataset)
}

// DeleteDataset 删除数据集
func (s *datasetService) DeleteDataset(id string) error {
	// 获取数据集信息
	dataset, err := s.datasetRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("dataset not found: %w", err)
	}

	// 删除文件
	if dataset.FilePath != "" {
		// 删除文件
		if err := os.Remove(dataset.FilePath); err != nil && !os.IsNotExist(err) {
			logger.Error("Failed to delete dataset file", "error", err, "path", dataset.FilePath)
		}

		// 尝试删除数据集目录
		datasetDir := filepath.Dir(dataset.FilePath)
		if err := os.RemoveAll(datasetDir); err != nil {
			logger.Error("Failed to delete dataset directory", "error", err, "path", datasetDir)
		}
	}

	// 从数据库中删除
	return s.datasetRepo.Delete(id)
}

// DownloadDataset 下载数据集
func (s *datasetService) DownloadDataset(id string) (string, error) {
	// 获取数据集信息
	dataset, err := s.datasetRepo.GetByID(id)
	if err != nil {
		return "", fmt.Errorf("dataset not found: %w", err)
	}

	// 确保文件存在
	if dataset.FilePath == "" {
		return "", errors.New("dataset has no file")
	}

	if _, err := os.Stat(dataset.FilePath); os.IsNotExist(err) {
		return "", errors.New("dataset file not found")
	}

	// 增加下载计数
	if err := s.datasetRepo.IncrementDownloadCount(id); err != nil {
		logger.Error("Failed to increment download count", "error", err, "datasetId", id)
	}

	return dataset.FilePath, nil
}
