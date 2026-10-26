package source

import (
	"strings"

	"micro-scrm/internal/app"
	membersDao "micro-scrm/internal/app/dao/members"
	microBossDao "micro-scrm/internal/app/dao/micro_boss"
	membersModel "micro-scrm/internal/app/models/members"
	microBossModel "micro-scrm/internal/app/models/micro_boss"
	coachkos "micro-scrm/internal/pkg/sdk/coach_kos"
)

type coachKOSUploader struct {
	memberID int64
	adminID  int64
}

func resolveCoachKOSUploaders(c app.Ctx, enterpriseID int64, posts []coachkos.PostResult) (map[string]coachKOSUploader, error) {
	employeeIDs := make([]string, 0, len(posts))
	seenIDs := make(map[string]struct{}, len(posts))
	for _, post := range posts {
		employeeID := strings.TrimSpace(post.EmployeeID)
		if employeeID == "" {
			continue
		}
		if _, exists := seenIDs[employeeID]; exists {
			continue
		}
		seenIDs[employeeID] = struct{}{}
		employeeIDs = append(employeeIDs, employeeID)
	}
	if len(employeeIDs) == 0 {
		return map[string]coachKOSUploader{}, nil
	}

	members, err := membersDao.NewMembers().GetByEmployeeIDs(c, enterpriseID, employeeIDs)
	if err != nil {
		return nil, err
	}
	mobiles := make([]string, 0, len(members))
	seenMobiles := make(map[string]struct{}, len(members))
	for _, member := range members {
		if member.Mobile == "" {
			continue
		}
		if _, exists := seenMobiles[member.Mobile]; exists {
			continue
		}
		seenMobiles[member.Mobile] = struct{}{}
		mobiles = append(mobiles, member.Mobile)
	}
	admins, err := microBossDao.NewEnterpriseAdmin().GetActiveByMobiles(c, enterpriseID, mobiles)
	if err != nil {
		return nil, err
	}
	return matchCoachKOSUploaders(employeeIDs, members, admins), nil
}

func matchCoachKOSUploaders(employeeIDs []string, members []membersModel.MemberModel, admins []microBossModel.EnterpriseAdminModel) map[string]coachKOSUploader {
	result := make(map[string]coachKOSUploader, len(employeeIDs))
	wanted := make(map[string]struct{}, len(employeeIDs))
	for _, employeeID := range employeeIDs {
		wanted[employeeID] = struct{}{}
	}

	adminByMobile := make(map[string]int64, len(admins))
	for _, admin := range admins {
		if admin.Mobile != "" {
			if _, exists := adminByMobile[admin.Mobile]; !exists {
				adminByMobile[admin.Mobile] = admin.ID
			}
		}
	}
	assign := func(employeeID string, member membersModel.MemberModel) {
		if _, exists := wanted[employeeID]; !exists {
			return
		}
		if _, exists := result[employeeID]; exists {
			return
		}
		result[employeeID] = coachKOSUploader{
			memberID: member.ID,
			adminID:  adminByMobile[member.Mobile],
		}
	}
	// 与员工主数据的匹配规则一致：userid 优先，open_userid 作为回退。
	for _, member := range members {
		assign(member.Userid, member)
	}
	for _, member := range members {
		assign(member.OpenUserid, member)
	}
	return result
}
