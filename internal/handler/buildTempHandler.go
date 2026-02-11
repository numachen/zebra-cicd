package handler

import (
	"gorm.io/gorm"

	"github.com/numachen/zebra-cicd/internal/model"
)

type BuildTemplateRepository struct {
	db *gorm.DB
}

func NewBuildTemplateRepository(db *gorm.DB) *BuildTemplateRepository {
	return &BuildTemplateRepository{db: db}
}

// Create 创建模板
func (r *BuildTemplateRepository) Create(template *model.BuildTemplate) error {
	return r.db.Create(template).Error
}

// GetByID 根据ID获取模板
func (r *BuildTemplateRepository) GetByID(id uint) (*model.BuildTemplate, error) {
	var template model.BuildTemplate
	if err := r.db.Preload("Repos").First(&template, id).Error; err != nil {
		return nil, err
	}
	return &template, nil
}

// List 获取模板列表
func (r *BuildTemplateRepository) List() ([]model.BuildTemplate, error) {
	var templates []model.BuildTemplate
	if err := r.db.Preload("Repos").Order("id DESC").Find(&templates).Error; err != nil {
		return nil, err
	}
	return templates, nil
}

// Update 更新模板
func (r *BuildTemplateRepository) Update(template *model.BuildTemplate) error {
	return r.db.Save(template).Error
}

// Delete 删除模板
func (r *BuildTemplateRepository) Delete(id uint) error {
	return r.db.Delete(&model.BuildTemplate{}, id).Error
}

// AddRepoToTemplate 添加模板到仓库关联
func (r *BuildTemplateRepository) AddRepoToTemplate(templateID, repoID uint) error {
	template := &model.BuildTemplate{ID: templateID}
	repo := &model.Repo{ID: repoID}
	return r.db.Model(template).Association("Repos").Append(repo)
}

// RemoveRepoFromTemplate 移除仓库与模板的关联
func (r *BuildTemplateRepository) RemoveRepoFromTemplate(templateID, repoID uint) error {
	template := &model.BuildTemplate{ID: templateID}
	repo := &model.Repo{ID: repoID}
	return r.db.Model(template).Association("Repos").Delete(repo)
}

// GetTemplatesByRepoID 根据仓库ID获取关联的模板
func (r *BuildTemplateRepository) GetTemplatesByRepoID(repoID uint) ([]model.BuildTemplate, error) {
	var repo model.Repo
	if err := r.db.Preload("Templates").First(&repo, repoID).Error; err != nil {
		return nil, err
	}

	// 将 []*model.BuildTemplate 转换为 []model.BuildTemplate
	templates := make([]model.BuildTemplate, len(repo.Templates))
	for i, template := range repo.Templates {
		templates[i] = *template
	}

	return templates, nil
}
