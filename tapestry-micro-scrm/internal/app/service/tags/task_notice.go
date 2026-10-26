package tagsSvc

import (
	"fmt"
	"micro-scrm/internal/app"
	"micro-scrm/internal/pkg/notice"
	"time"
)

type TaskNotice struct {
}

func NewTaskNotice() *TaskNotice {
	return &TaskNotice{}
}

type TaskMessageData struct {
	Title            string        `json:"title" comment:"企微告警标题"`
	EnterpriseId     int64         `json:"enterprise_id"`
	NeedMarkCount    int           `json:"need_mark_count" comment:"需要打标的数量"`
	SuccessMarkCount int64         `json:"success_mark_count" comment:"成功打标的数量"`
	UseMemory        float64       `json:"use_memory" comment:"占用内存"`
	UseCpu           int           `json:"use_cpu" comment:"占用cpu"`
	UseTime          time.Duration `json:"use_time" comment:"执行时长"`
}

// MarkSendNotice 打标发送结果通知
func (p TaskNotice) MarkSendNotice(c app.Ctx, start bool, data TaskMessageData) {
	if start {
		content := fmt.Sprintf("%s \n"+
			">企业id：<font color='comment'>%d</font>\n"+
			">本次打标人数：<font color='comment'>%d</font>\n"+
			">RequestId：<font color='comment'>%v</font>\n",
			data.Title,
			data.EnterpriseId,
			data.NeedMarkCount,
			c.ID(),
		)
		notice.NewWechat().SendMarkMarkdownMessage(c, content)
		return
	}

	content := fmt.Sprintf("%s：\n"+
		">企业id：<font color='comment'>%d</font>\n"+
		">本次打标人数：<font color='comment'>%d</font>\n"+
		">成功打标人数：<font color='comment'>%d</font>\n"+
		">占用CPU：<font color='red'>%d</font>\n"+
		">占用内存：<font color='red'>%.2f MB</font>\n"+
		">执行时长：<font color='red'>%s</font>\n"+
		">RequestId：<font color='comment'>%v</font>\n",
		data.Title,
		data.EnterpriseId,
		data.NeedMarkCount,
		data.SuccessMarkCount,

		data.UseCpu,
		data.UseMemory,
		data.UseTime,
		c.ID(),
	)
	notice.NewWechat().SendMarkMarkdownMessage(c, content)
}
