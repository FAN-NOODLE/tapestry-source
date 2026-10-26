package delivery

import (
	"micro-scrm/internal/modules/customer"
	"micro-scrm/internal/modules/source"
)

func TransListParams(params *customer.RepeatListRequest) *customer.RepeatListParams {
	serviceParams := &customer.RepeatListParams{
		EnterpriseID: params.EnterpriseID,
		Page:         params.Page,
		Size:         params.Size,
		CustomerName: params.Search.CustomerName,
		DepartmentID: params.Search.DepartmentID,
		MemberID:     params.Search.MemberID,
	}
	return serviceParams
}

func TransImgParams(params *source.SourcesFileRequest) *source.SourcesImageRequest {
	var images []*source.SourcesImageRequestAdd
	imageItem := &source.SourcesImageRequestAdd{
		ID:           params.ID,
		EnterpriseID: params.EnterpriseID,
		Name:         params.Name,
		URL:          params.URL,
		Ext:          params.Ext,
		Size:         params.Size,
	}
	images = append(images, imageItem)

	serviceParams := &source.SourcesImageRequest{
		EnterpriseID: params.EnterpriseID,
		GroupID:      params.GroupID,
		AdminID:      params.AdminID,
		Images:       images,
	}
	return serviceParams
}

func TransVideoParams(params *source.SourcesFileRequest) *source.RequestSourceVideo {
	serviceParams := &source.RequestSourceVideo{
		ID:           params.ID,
		EnterpriseID: params.EnterpriseID,
		GroupID:      params.GroupID,
		CoverUrl:     " ",
		Name:         params.Name,
		Size:         params.Size,
		Url:          params.URL,
		Ext:          params.Ext,
		AdminID:      params.AdminID,
	}
	return serviceParams
}
