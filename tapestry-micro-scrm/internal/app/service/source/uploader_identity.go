package source

import (
	"micro-scrm/internal/app"
	membersDao "micro-scrm/internal/app/dao/members"
	microBossDao "micro-scrm/internal/app/dao/micro_boss"
	membersModel "micro-scrm/internal/app/models/members"
	microBossModel "micro-scrm/internal/app/models/micro_boss"
	sourceModel "micro-scrm/internal/app/models/source"
)

func mergeUploaderIdentity(items []*sourceModel.ResponseSelectList, admins []microBossModel.EnterpriseAdminModel, members []membersModel.MemberInfo) {
	memberByID := make(map[int64]membersModel.MemberInfo, len(members))
	memberByMobile := make(map[string]membersModel.MemberInfo, len(members))
	for _, member := range members {
		memberByID[member.Id] = member
		if member.Mobile != "" {
			memberByMobile[member.Mobile] = member
		}
	}

	adminByID := make(map[int64]microBossModel.EnterpriseAdminModel, len(admins))
	for _, admin := range admins {
		adminByID[admin.ID] = admin
	}

	for _, item := range items {
		if item == nil {
			continue
		}
		if member, ok := memberByID[item.MemberId]; item.MemberId > 0 && ok {
			item.MemberName = member.Name
			item.MemberAvatar = member.Avatar
			continue
		}

		admin, ok := adminByID[item.AdminId]
		if item.AdminId <= 0 || !ok {
			continue
		}
		item.MemberName = admin.Username
		if member, ok := memberByMobile[admin.Mobile]; admin.Mobile != "" && ok {
			item.MemberId = member.Id
			item.MemberName = member.Name
			item.MemberAvatar = member.Avatar
		}
	}
}

// enrichUploaderIdentity is optional display enrichment. Missing upstream identities
// and lookup failures are intentionally ignored so source lists remain available.
func enrichUploaderIdentity(c app.Ctx, enterpriseID int64, items []*sourceModel.ResponseSelectList) {
	memberIDs := make([]int64, 0)
	adminIDs := make([]int64, 0)
	memberIDSet := make(map[int64]struct{})
	adminIDSet := make(map[int64]struct{})
	for _, item := range items {
		if item == nil {
			continue
		}
		if item.MemberId > 0 {
			if _, exists := memberIDSet[item.MemberId]; !exists {
				memberIDSet[item.MemberId] = struct{}{}
				memberIDs = append(memberIDs, item.MemberId)
			}
		} else if item.AdminId > 0 {
			if _, exists := adminIDSet[item.AdminId]; !exists {
				adminIDSet[item.AdminId] = struct{}{}
				adminIDs = append(adminIDs, item.AdminId)
			}
		}
	}

	admins := make([]microBossModel.EnterpriseAdminModel, 0)
	if len(adminIDs) > 0 {
		if found, err := microBossDao.NewEnterpriseAdmin().GetMany(c, map[string]any{"id": adminIDs, "enterprise_id": enterpriseID}); err == nil {
			admins = found
		}
	}

	members := make([]membersModel.MemberInfo, 0)
	if len(memberIDs) > 0 {
		if found, err := membersDao.NewMembers().GetUnscopedMany(c, map[string]any{"id": memberIDs, "enterprise_id": enterpriseID}); err == nil {
			for _, member := range found {
				members = append(members, membersModel.MemberInfo{
					Id: member.ID, EnterpriseId: member.EnterpriseID, Name: member.Name,
					Mobile: member.Mobile, Avatar: member.Avatar,
				})
			}
		}
	}

	mobiles := make([]string, 0, len(admins))
	for _, admin := range admins {
		if admin.Mobile != "" {
			mobiles = append(mobiles, admin.Mobile)
		}
	}
	if len(mobiles) > 0 {
		if found, err := membersDao.NewMembers().GetUnscopedMemberByMobiles(c, enterpriseID, mobiles); err == nil {
			members = append(members, found...)
		}
	}

	mergeUploaderIdentity(items, admins, members)
}
