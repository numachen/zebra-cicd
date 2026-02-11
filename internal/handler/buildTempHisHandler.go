package handler

import (
	"gorm.io/gorm"

	"github.com/numachen/zebra-cicd/internal/model"
	"github.com/numachen/zebra-cicd/internal/types"
)

type TemplateHistoryRepository struct {
	db *gorm.DB
}

func NewTemplateHistoryRepository(db *gorm.DB) *TemplateHistoryRepository {
	return &TemplateHistoryRepository{db: db}
}

// Create 创建模板历史记录
func (r *TemplateHistoryRepository) Create(history *model.TemplateHistory) error {
	return r.db.Create(history).Error
}

// GetHistoryByTemplateID 根据模板ID获取历史记录
// GetHistoryByTemplateID 根据模板ID获取历史记录
func (r *TemplateHistoryRepository) GetHistoryByTemplateID(templateID uint) ([]types.TemplateHistoryResponse, error) {
	var histories []model.TemplateHistory
	if err := r.db.Where("template_id = ?", templateID).Order("created_at DESC").Find(&histories).Error; err != nil {
		return nil, err
	}

	// 转换为响应类型
	responses := make([]types.TemplateHistoryResponse, len(histories))
	for i, history := range histories {
		responses[i] = types.TemplateHistoryResponse{
			ID:         history.ID,
			TemplateID: history.TemplateID,
			Modifier:   history.Modifier,
			Dockerfile: history.Dockerfile,
			Pipeline:   history.Pipeline,
			CreatedAt:  history.CreatedAt,
		}
	}

	return responses, nil
}
