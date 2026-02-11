package model

import (
	"github.com/numachen/zebra-cicd/pkg/timeutil"
)

type BuildTemplate struct {
	ID         uint              `gorm:"primaryKey" json:"id"`
	Name       string            `gorm:"size:255;uniqueIndex;not null;comment:名称" json:"name"`
	Language   string            `gorm:"size:100;not null;comment:开发语言" json:"language"`
	Creator    string            `gorm:"size:100;comment:创建人" json:"creator"`
	Updater    string            `gorm:"size:100;comment:修改人" json:"updater"`
	Dockerfile string            `gorm:"type:text;comment:Dockerfile" json:"dockerfile"`
	Pipeline   string            `gorm:"type:text;comment:Pipeline" json:"pipeline"`
	CreatedAt  timeutil.JSONTime `gorm:"comment:创建时间" json:"created_at"`
	UpdatedAt  timeutil.JSONTime `gorm:"comment:更新时间" json:"updated_at"`

	// 关联Repo表的多对多关系
	Repos []Repo `gorm:"many2many:repo_templates;" json:"repos,omitempty"`
}

type TemplateHistory struct {
	ID         uint              `gorm:"primaryKey" json:"id"`
	TemplateID uint              `gorm:"index;comment:模板ID" json:"template_id"`
	Modifier   string            `gorm:"size:100;comment:修改人" json:"modifier"`
	Dockerfile string            `gorm:"type:text;comment:Dockerfile" json:"dockerfile"`
	Pipeline   string            `gorm:"type:text;comment:Pipeline" json:"pipeline"`
	CreatedAt  timeutil.JSONTime `gorm:"comment:创建时间" json:"created_at"`

	// 关联模板
	Template BuildTemplate `json:"template"`
}
