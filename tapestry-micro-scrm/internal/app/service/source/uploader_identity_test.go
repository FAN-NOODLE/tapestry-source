package source

import (
	"testing"

	"micro-scrm/internal/app/models"
	membersModel "micro-scrm/internal/app/models/members"
	microBossModel "micro-scrm/internal/app/models/micro_boss"
	sourceModel "micro-scrm/internal/app/models/source"
)

func TestMergeUploaderIdentityUsesStoredMemberID(t *testing.T) {
	items := []*sourceModel.ResponseSelectList{{MemberId: 11}}
	members := []membersModel.MemberInfo{{Id: 11, Name: "企微张三", Avatar: "avatar-11"}}

	mergeUploaderIdentity(items, nil, members)

	if items[0].MemberId != 11 || items[0].MemberName != "企微张三" || items[0].MemberAvatar != "avatar-11" {
		t.Fatalf("unexpected stored member identity: %#v", items[0])
	}
}

func TestMergeUploaderIdentityMapsHistoricalAdminByMobile(t *testing.T) {
	items := []*sourceModel.ResponseSelectList{{AdminId: 21}}
	admins := []microBossModel.EnterpriseAdminModel{{
		MicroBaseModel: models.MicroBaseModel{BaseModel: models.BaseModel{ID: 21}},
		Username:       "后台李四",
		Mobile:         "13800000000",
	}}
	members := []membersModel.MemberInfo{{Id: 31, Mobile: "13800000000", Name: "企微李四", Avatar: "avatar-31"}}

	mergeUploaderIdentity(items, admins, members)

	if items[0].MemberId != 31 || items[0].MemberName != "企微李四" || items[0].MemberAvatar != "avatar-31" {
		t.Fatalf("unexpected historical admin identity: %#v", items[0])
	}
}

func TestMergeUploaderIdentitySilentlyFallsBackForMissingData(t *testing.T) {
	items := []*sourceModel.ResponseSelectList{{}, {AdminId: 22}}
	admins := []microBossModel.EnterpriseAdminModel{{
		MicroBaseModel: models.MicroBaseModel{BaseModel: models.BaseModel{ID: 22}},
		Username:       "后台王五",
	}}

	mergeUploaderIdentity(items, admins, nil)

	if items[0].MemberId != 0 || items[0].MemberName != "" || items[0].MemberAvatar != "" {
		t.Fatalf("empty upstream identity must remain empty: %#v", items[0])
	}
	if items[1].MemberId != 0 || items[1].MemberName != "后台王五" || items[1].MemberAvatar != "" {
		t.Fatalf("admin-only identity must fall back silently: %#v", items[1])
	}
}
