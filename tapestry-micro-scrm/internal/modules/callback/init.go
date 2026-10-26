package callback

// NewFissionHandle 创建裂变活动处理器
func NewFissionHandle(
	f FissionService,
	m MemberService,
	c CustomerService,
	e EnterpriseService,
) *FissionHandle {
	return &FissionHandle{
		fissionService:    f,
		memberService:     m,
		customerService:   c,
		enterpriseService: e,
	}
}
