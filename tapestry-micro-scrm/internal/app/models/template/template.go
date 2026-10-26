package templateModels

import "micro-scrm/internal/app/models"

type TemplateModel struct {
	models.MicroBaseModel
	ID                 int    `json:"id" gorm:"primaryKey;column:id;type:int(10) auto_increment"`
	TemplateName       string `json:"template_name" gorm:"column:template_name;type:varchar(255);not null;comment:模板名称"`
	TemplatePictureURL string `json:"template_picture_url" gorm:"column:template_picture_url;type:varchar(255);not null;comment:模板背景图url"`
	TemplateContent    string `json:"template_content" gorm:"column:template_content;type:varchar(255);not null;comment:模板的文案"`
	TemplateType       int32  `json:"template_type" gorm:"column:template_type;type:tinyint(1);not null;default:1;comment:模板的类型 1:贺卡"`
}

func (TemplateModel) TableName() string {
	return "template"
}
