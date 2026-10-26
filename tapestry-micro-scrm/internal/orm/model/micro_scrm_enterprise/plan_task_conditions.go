package micro_scrm_enterprise

import (
	"fmt"

	"micro-scrm/internal/orm/model"
)

type PlanTaskConditionModel struct {
	model.BranchModel

	PlanId      int64  `gorm:"column:plan_id;default:0;NOT NULL" json:"plan_id"`           // 条件对应的子任务对应的计划的ID
	TaskId      int64  `gorm:"column:task_id;default:0;NOT NULL" json:"task_id"`           // 条件对应的子任务的ID
	Conditions  string `gorm:"column:conditions;NOT NULL" json:"conditions"`               // 条件内容，json结构;包含：多个维度的结构
	Content     string `gorm:"column:content;NOT NULL" json:"content"`                     // 条件内容，json结构；包含可以推送给顾客的所有信息：文本/图片/视频/小程序/链接
	ContentType int    `gorm:"column:content_type;default:0;NOT NULL" json:"content_type"` // 消息类型；1文本；2图片；3语音；4视频；5位置；6链接
}
type Content struct {
	Text struct {
		Content string `json:"content"`
	} `json:"text"`
	Image struct {
		ImgList []struct {
			PicUrl  string `json:"pic_url"`
			MediaId string `json:"media_id"`
		} `json:"imgList"`
	} `json:"image"`
	Voice struct {
		MediaId string `json:"media_id"`
		Format  string `json:"format"`
	} `json:"voice"`
	Video struct {
		MediaId       string `json:"media_id"`
		ThumbMediaId  string `json:"thumb_media_id"`
		VideoUrl      string `json:"video_url"`
		VideoImageUrl string `json:"video_image_url"`
	} `json:"video"`
	Location struct {
		Location_X string `json:"location_x"`
		Location_Y string `json:"location_y"`
		Scale      string `json:"scale"`
		Label      string `json:"label"`
	} `json:"location"`
	Link struct {
		Id          int64  `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Url         string `json:"url"` // 以link_url为主
		LinkUrl     string `json:"link_url"`
		PicUrl      string `json:"pic_url"`
		CoverUrl    string `json:"cover_url"` // PicUrl的地址
	} `json:"link"`
	Applets struct {
		Id            int64  `json:"id"`
		Type          int    `json:"type"`
		AdminId       int64  `json:"admin_id"`
		Title         string `json:"title"`
		MiniProgramId int64  `json:"mini_program_id"`
		AppId         string `json:"app_id"`
		AppName       string `json:"app_name"`
		Path          string `json:"path"`
		CoverUrl      string `json:"cover_url"`
	} `json:"applets"`
	Card struct {
		Id                 int      `json:"id"`
		TemplateName       string   `json:"template_name"`
		TemplatePictureURL []string `json:"template_picture_url"`
		TemplateContent    string   `json:"template_content"`
	} `json:"card"`
}

type ConditionDetail struct {
	Tag     []TagInfo   `json:"tag"`
	Logical LogicalInfo `json:"logical"` // 逻辑条件
}
type LogicalInfo struct {
	Remark     string `json:"remark"`     // 备注说明
	Nubmer     int64  `json:"number"`     // 数量
	Expression string `json:"expression"` // 表达式 //等于(eq)、大于(gt)、大于等于(ge)、小于(lt)、小于等于(le)
}
type TagInfo struct {
	GroupName string `json:"group_name"`
	GroupID   int64  `json:"group_id"`
	QwGroupId string `json:"qw_group_id"`
	TagName   string `json:"tag_name"`
	TagID     int64  `json:"tag_id"`
	QwTagID   string `json:"qw_tag_id"`
}

func (PlanTaskConditionModel) TableName() string {
	return "plan_task_conditions"
}

func (p PlanTaskConditionModel) AllTableName(enterpriseID int64) string {
	return fmt.Sprintf(" %s%d.%s ", p.DBPrefix(), enterpriseID, p.TableName())
}

// func (c *PlanTaskConditionModel) AfterFind(db gorm.DB) (err error) {
// 	_ = json.Unmarshal([]byte(c.Conditions), &c.ConditionDetail)
// 	_ = json.Unmarshal([]byte(c.Content), &c.ContentDetail)

// 	return nil
// }
