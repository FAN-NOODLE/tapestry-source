package marketing

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
)

func (s Service) ShortUrlAdd(c ctx.Context, request ShortUrlAddReq) error {
	err := s.repo.ShortUrlAdd(c, request)
	if err != nil {
		return err
	}
	return err
}
func (s Service) ShortUrlDel(c ctx.Context, request ShortUrlDelReq) error {
	err := s.repo.ShortUrlDel(c, request)
	if err != nil {
		log.Error(c).Msgf("短链删除失败：%+v", err)
		return err
	}
	return err
}

func (s Service) TwoLevelIds(c ctx.Context, request ShortUrlAddReq) ([]ShortUrlDataReq, error) {
	Ids, err := s.repo.TwoLevelIds(c, request)
	if err != nil {
		log.Error(c).Msgf("TwoLevelIdsErr：%+v", err)
		return nil, err
	}
	return Ids, err
}

func (s Service) ShortUrlData(c ctx.Context, request ShortUrlDataReq) (ShortUrlDataRes, error) {
	shortUrlData, err := s.repo.ShortUrlData(c, request)
	if err != nil {
		log.Error(c).Msgf("短链添加失败：%+v", err)
		return shortUrlData, err
	}
	return shortUrlData, err
}

func (s Service) ShortUrlList(c ctx.Context, request ShortUrlListReq) (int64, *[]*ShortUrlDataRes, error) {
	offset := (request.CurrentPage - 1) * request.PageNum
	limit := request.PageNum
	option := &repos.FindPageOption{
		TplPrefix: "t_marketings/short-url-list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]*ShortUrlDataRes{}
		},
	}

	total, list, err := s.repo.FindPage(c, repos.PageVO{Offset: offset, Limit: limit}, request, option)
	if err != nil {
		return 0, nil, err
	}

	resList := list.(*[]*ShortUrlDataRes)

	return total, resList, err
}

func (s Service) AllShortUrlList(c ctx.Context, request ShortUrlListReq) (*[]*ShortUrlDataRes, error) {
	list, err := s.repo.AllShortUrlList(c, request)
	if err != nil {
		return nil, err
	}

	return list, err
}
