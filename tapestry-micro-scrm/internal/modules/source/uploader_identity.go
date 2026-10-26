package source

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	membersDao "micro-scrm/internal/app/dao/members"
	microBossDao "micro-scrm/internal/app/dao/micro_boss"
)

// resolveAdminUploaderMemberID is best-effort identity enrichment for new
// management-console uploads. Missing legacy/upstream identity is valid.
func resolveAdminUploaderMemberID(c ctx.Context, enterpriseID, adminID int64) int64 {
	if enterpriseID <= 0 || adminID <= 0 {
		return 0
	}
	admins, err := microBossDao.NewEnterpriseAdmin().GetMany(c, map[string]any{
		"id": adminID, "enterprise_id": enterpriseID,
	})
	if err != nil || len(admins) == 0 || admins[0].Mobile == "" {
		return 0
	}
	member, err := membersDao.NewMembers().GetMemberByMobile(c, enterpriseID, admins[0].Mobile)
	if err != nil {
		return 0
	}
	return member.Id
}
