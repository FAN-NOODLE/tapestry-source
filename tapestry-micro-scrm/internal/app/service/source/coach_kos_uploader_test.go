package source

import (
	"testing"

	"micro-scrm/internal/app/models"
	membersModel "micro-scrm/internal/app/models/members"
	microBossModel "micro-scrm/internal/app/models/micro_boss"
)

func TestCoachKOS上传人优先匹配Userid并关联管理员(t *testing.T) {
	employeeIDs := []string{"same-code", "open-code", "missing"}
	members := []membersModel.MemberModel{
		{MicroBaseModel: models.MicroBaseModel{BaseModel: models.BaseModel{ID: 1}}, Userid: "other", OpenUserid: "same-code", Mobile: "13800000001"},
		{MicroBaseModel: models.MicroBaseModel{BaseModel: models.BaseModel{ID: 2}}, Userid: "same-code", Mobile: "13800000002"},
		{MicroBaseModel: models.MicroBaseModel{BaseModel: models.BaseModel{ID: 3}}, Userid: "third", OpenUserid: "open-code", Mobile: "13800000003"},
	}
	admins := []microBossModel.EnterpriseAdminModel{
		{MicroBaseModel: models.MicroBaseModel{BaseModel: models.BaseModel{ID: 22}}, Mobile: "13800000002"},
	}

	got := matchCoachKOSUploaders(employeeIDs, members, admins)
	if got["same-code"] != (coachKOSUploader{memberID: 2, adminID: 22}) {
		t.Fatalf("userid 应优先于 open_userid：%+v", got["same-code"])
	}
	if got["open-code"] != (coachKOSUploader{memberID: 3}) {
		t.Fatalf("open_userid 回退匹配异常：%+v", got["open-code"])
	}
	if _, exists := got["missing"]; exists {
		t.Fatalf("未匹配到员工时不应写入上传人：%+v", got["missing"])
	}
}
