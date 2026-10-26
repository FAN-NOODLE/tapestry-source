package drainage

// NewDrainage 欢迎语
func NewDrainage(repository Repository) *DrainageService {
	return &DrainageService{
		repo: repository,
	}
}

// NewContact 联系我二维码
func NewContact(repository Repository) *ContactService {
	return &ContactService{
		repo: repository,
	}
}
