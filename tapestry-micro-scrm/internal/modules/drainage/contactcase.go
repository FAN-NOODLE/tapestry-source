package drainage

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"image"
	"image/draw"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	departmentsDao "micro-scrm/internal/app/dao/departments"
	microScrmDao "micro-scrm/internal/app/dao/drainage"
	membersDao "micro-scrm/internal/app/dao/members"
	"net/http"
	"strings"
	"time"

	"github.com/nfnt/resize"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/helper"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/customer"
	"micro-scrm/internal/pkg/sdk/work"
	"micro-scrm/internal/pkg/uid"
)

type ContactService struct {
	repo Repository
}

// GetContactQrCode 根据二维码id获取二维码链接
func (s ContactService) GetContactQrCode(c ctx.Context, ID int64, enterpriseID int64) (string, string, error) {
	contact, err := s.GetContact(c, ContactDetailsRequest{
		ID:           ID,
		EnterpriseId: enterpriseID,
	})
	if err != nil {
		return "", "", err
	}
	return contact.QrCodeUrl, contact.ConfigId, nil
}

// CreateContact 创建联系我二维码
func (s ContactService) CreateContact(c ctx.Context, enterpriseID int64, ID int64, name string, channelID int64, qrCodeUrl string, configID string, cType int) error {
	contacts := Contacts{
		ID:           ID,
		Name:         name,
		ChannelId:    channelID,
		AutoFriend:   0,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
		EnterpriseId: enterpriseID,
		QrCodeUrl:    qrCodeUrl,
		MemberCount:  0,
		DeletedAt:    0,
		Type:         cType,
		ConfigId:     configID,
	}
	err := s.repo.ContactAdd(c, contacts)
	if err != nil {
		return err
	}
	return nil
}

func (s ContactService) ListChannelByQrCodeIds(ctx ctx.Context, qrCodeIds []int64, enterpriseID int64, page int, size int) (int64, map[int64]customer.ListChannelByQrCodeIdsResponse) {
	req := ContactsListResponse{
		IDs:          qrCodeIds,
		Page:         page,
		Size:         size,
		EnterpriseId: enterpriseID,
		IdsCount:     len(qrCodeIds),
	}
	total, contacts := s.List(ctx, req)
	channelContact := make(map[int64]customer.ListChannelByQrCodeIdsResponse)
	if total > 0 {
		for _, item := range *contacts {
			channelContact[item.ID] = customer.ListChannelByQrCodeIdsResponse{
				ChannelID:    item.ChannelId,
				ChannelName:  item.ChannelName,
				QrCodeID:     item.ID,
				EnterpriseID: item.EnterpriseId,
			}
		}

		return total, channelContact
	}
	return 0, nil
}

// ContactCreate 创建联系我二维码
func (s ContactService) ContactCreate(c ctx.Context, request ContactCreateRequest) (int64, error) {
	workImpl := new(work.QyWorkImpl).SetDomain(request.Domain)
	request.ID = uid.ID(ctx.Wrap(c))
	res, qyErr := workImpl.GetWxEexternalContact(ctx.New(), work.ContactCreateRequest{
		ID:           request.ID,
		CodeCame:     request.CodeCame,
		ChannelId:    request.ChannelId,
		AutoFriend:   request.AutoFriend,
		MemberIds:    request.MemberIds,
		ThirdIds:     request.ThirdIds,
		EnterpriseId: request.EnterpriseId,
		Type:         work.DrainageContactsTypeMultiple,
		UserIds:      request.UserIds,
		State:        request.State,
	})
	if res == nil {
		e, ok := qyErr.(*errorx.Error)
		// 判断是不是企业微信错误码
		if ok && e.Code == 41054 {
			return 0, errorx.NewErrorMsg(e.Code, "人员未激活")
		}
		if ok && e.Code == 40098 {
			return 0, errorx.NewErrorMsg(e.Code, "员工未实名认证")
		}
		if ok && e.Code == 41041 {
			return 0, errorx.NewErrorMsg(e.Code, "总的人员数量不能超过100")
		}
		return 0, errorx.NewErrorMsg(e.Code, fmt.Sprintf("企业微信服务异常,errcode:%d", e.Code))
	}
	contacts := Contacts{
		ID:           request.ID,
		Name:         request.CodeCame,
		ChannelId:    request.ChannelId,
		AutoFriend:   request.AutoFriend,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
		EnterpriseId: request.EnterpriseId,
		QrCodeUrl:    res.QrCode,
		MemberCount:  0,
		DeletedAt:    0,
		ConfigId:     res.ConfigId,
		Type:         request.Type,
	}
	contactDataList := make([]*ContactMembers, 0)
	memberNumber := 0
	for _, memberId := range request.MemberIds {
		contactMember := ContactMembers{
			ID:           uid.ID(c),
			ContactId:    contacts.ID,
			EnterpriseId: request.EnterpriseId,
			MemberId:     memberId,
			CreatedAt:    time.Now().Unix(),
			UpdatedAt:    time.Now().Unix(),
			DeletedAt:    0,
		}
		contactDataList = append(contactDataList, &contactMember)
		memberNumber++
	}
	for _, thirdId := range request.ThirdIds {
		contactMember := ContactMembers{
			ID:           uid.ID(c),
			ContactId:    contacts.ID,
			EnterpriseId: request.EnterpriseId,
			ThirdId:      thirdId,
			CreatedAt:    time.Now().Unix(),
			UpdatedAt:    time.Now().Unix(),
			DeletedAt:    0,
		}
		contactDataList = append(contactDataList, &contactMember)
	}

	contacts.MemberCount = memberNumber
	err := s.repo.ContactAdd(c, contacts)
	if err != nil {
		return 0, err
		// helper.Must(err)
	}
	// 批量添加联系人
	contactMembersAddError := s.repo.ContactMembersAdd(c, contactDataList)
	if contactMembersAddError != nil {
		helper.Must(contactMembersAddError)
	}
	return contacts.ID, nil
}

func (s ContactService) ContactUpdate(c ctx.Context, request ContactCreateRequest) (int64, error) {
	contacts := Contacts{
		ID:           request.ID,
		Name:         request.CodeCame,
		ChannelId:    request.ChannelId,
		AutoFriend:   request.AutoFriend,
		EnterpriseId: request.EnterpriseId,
		UpdatedAt:    time.Now().Unix(),
	}
	contactsMemberDelete := Contacts{
		ID:        request.ID,
		DeletedAt: time.Now().Unix(),
	}
	contactMembersDeleteError := s.repo.ContactMembersDelete(c, contactsMemberDelete)
	if contactMembersDeleteError != nil {
		helper.Must(contactMembersDeleteError)
	}
	memberNumber := 0
	contactMembersDatas := make([]*ContactMembers, 0)
	for _, memberId := range request.MemberIds {
		contactMember := ContactMembers{
			ID:           uid.ID(c),
			ContactId:    contacts.ID,
			EnterpriseId: request.EnterpriseId,
			MemberId:     memberId,
			CreatedAt:    time.Now().Unix(),
			UpdatedAt:    time.Now().Unix(),
			DeletedAt:    0,
		}
		contactMembersDatas = append(contactMembersDatas, &contactMember)
		memberNumber++
	}
	for _, thirdId := range request.ThirdIds {
		contactMember := ContactMembers{
			ID:           uid.ID(c),
			ContactId:    contacts.ID,
			EnterpriseId: request.EnterpriseId,
			ThirdId:      thirdId,
			CreatedAt:    time.Now().Unix(),
			UpdatedAt:    time.Now().Unix(),
			DeletedAt:    0,
		}
		contactMembersDatas = append(contactMembersDatas, &contactMember)
	}
	contacts.MemberCount = memberNumber
	err := s.repo.ContactUpdate(c, contacts)
	if err != nil {
		return 0, err
		// helper.Must(err)
	}
	// 批量添加联系人
	contactMembersAddError := s.repo.ContactMembersAdd(c, contactMembersDatas)
	if contactMembersAddError != nil {
		return 0, contactMembersAddError
		// helper.Must(contactMembersAddError)
	}
	return contacts.ID, nil
}

func (s ContactService) GetContact(c ctx.Context, request ContactDetailsRequest) (*Contacts, error) {
	contacts, _ := s.repo.ContactSelect(c, request)
	return contacts, nil
}

// ContactQrcodeList  联系我二维码详情
func (s ContactService) ContactQrcodeList(c ctx.Context, request ContactDetailsRequest) (*Contacts, error) {
	contacts, erro := s.repo.ContactSelect(c, request)
	if erro != nil {
		return nil, erro
	}
	// 查询联系我二维码成员
	memberData, _ := s.repo.ContactMemberSelect(c, request)
	contacts.MemberDatas = memberData
	if memberData == nil {
		contacts.MemberDatas = []MemberData{}
	}
	departmentData, _ := s.repo.ContactDepartmentSelect(c, request)
	contacts.DepartmentData = departmentData
	if departmentData == nil {
		contacts.DepartmentData = []DepartmentData{}
	}
	return contacts, nil
}

// ContactDetails 联系我二维码详情
func (s ContactService) ContactDetails(c ctx.Context, request ContactDetailsRequest) (*Contacts, error) {
	contacts, erro := s.repo.ContactSelect(c, request)
	if erro != nil {
		return nil, erro
	}
	// 查询联系我二维码成员
	memberData, _ := s.repo.ContactMemberSelect(c, request)
	contacts.MemberDatas = memberData
	if memberData == nil {
		contacts.MemberDatas = []MemberData{}
	}
	departmentData, _ := s.repo.ContactDepartmentSelect(c, request)
	contacts.DepartmentData = departmentData
	if departmentData == nil {
		contacts.DepartmentData = []DepartmentData{}
	}
	return contacts, nil
}
func (s ContactService) ContactDelete(c ctx.Context, request ContactDeleteRequest) error {
	err := s.repo.ContactDelete(c, request)
	/*if err != nil {
		helper.Must(err)
	}*/
	return err
}

// List 联系我二维码列表
func (s ContactService) List(c ctx.Context, request ContactsListResponse) (int64, *[]Contacts) {
	offset := (request.Page - 1) * request.Size
	limit := request.Size
	option := &repos.FindPageOption{
		TplPrefix: "t_drainage_contact/contacts-select-list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]Contacts{}
		},
	}
	if request.Name != "" {
		var build strings.Builder
		build.WriteString("%")
		build.WriteString(request.Name)
		build.WriteString("%")
		request.Name = build.String()
	}

	total, list, err := s.repo.FindPage(c, repos.PageVO{Offset: offset, Limit: limit}, request, option)
	if err != nil {
		return 0, nil
	}
	contactsList := list.(*[]Contacts)
	var ids []int64
	for _, contact := range *contactsList {
		ids = append(ids, contact.ID)
	}
	scopeList, err := microScrmDao.NewDrainageContactMemberDao().GetMany(c, "enterprise_id=? AND contact_id in ?", request.EnterpriseId, ids)
	if err != nil {
		return 0, nil
	}
	var memberIds []int64
	var contactMemberMap = make(map[int64][]int64)
	var contactThirdMap = make(map[int64][]int)
	var thirdIds []int

	for _, scope := range scopeList {
		if scope.MemberId > 0 {
			memberIds = append(memberIds, scope.MemberId)
			contactMemberMap[scope.ContactId] = append(contactMemberMap[scope.ContactId], scope.MemberId)
		}
		if scope.ThirdId > 0 {
			thirdIds = append(thirdIds, scope.ThirdId)
			contactThirdMap[scope.ContactId] = append(contactThirdMap[scope.ContactId], scope.ThirdId)
		}
	}
	var memberMap = make(map[int64]string)
	if len(memberIds) > 0 {
		memberList, _ := membersDao.NewMembers().GetMany(c, "enterprise_id = ? AND id in ?", request.EnterpriseId, memberIds)
		for _, member := range memberList {
			memberMap[member.ID] = member.Name
		}
	}
	var departmentMap = make(map[int]string)
	if len(thirdIds) > 0 {
		departmentList, _ := departmentsDao.NewDepartments().GetMany(c, "enterprise_id = ? AND third_id in ?", request.EnterpriseId, thirdIds)
		for _, department := range departmentList {
			departmentMap[department.ThirdID] = department.Name
		}
	}
	for i := range *contactsList {
		contact := &(*contactsList)[i] // 获取原始元素的指针
		if contactMemberIds, ok := contactMemberMap[contact.ID]; ok {
			for _, memberId := range contactMemberIds {
				if memberName, ok1 := memberMap[memberId]; ok1 {
					contact.MemberDatas = append(contact.MemberDatas, MemberData{
						MemberName: memberName,
					})
				}
			}
		}
		if contactThirdIds, ok := contactThirdMap[contact.ID]; ok {
			for _, thirdId := range contactThirdIds {
				if departmentName, ok1 := departmentMap[thirdId]; ok1 {
					contact.DepartmentData = append(contact.DepartmentData, DepartmentData{
						DepartmentName: departmentName,
					})
				}
			}
		}
	}
	return total, contactsList
}

func (s ContactService) DownLoadQrCode(c ctx.Context, request QrCodeRequest) (string, error) {
	// 获取远端图片
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	res, imageUrlErr := client.Get(request.ImageUrl)

	var width = request.Width
	var height = request.Height

	if imageUrlErr != nil {
		return "", imageUrlErr
	}
	defer res.Body.Close()

	// 读取获取的[]byte数据
	data, _ := ioutil.ReadAll(res.Body)

	imageBase64 := base64.StdEncoding.EncodeToString(data)
	ddd, _ := base64.StdEncoding.DecodeString(imageBase64) // 成图片文件并把文件写入到buffer
	bbb := bytes.NewBuffer(ddd)                            // 必须加一个buffer 不然没有read方法就会报错

	img, _, err := image.Decode(bbb)
	fmt.Println(err)

	// 调用resize库进行图片缩放(高度填0，resize.Resize函数中会自动计算缩放图片的宽高比)
	m := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)

	newImg := image.NewNRGBA(image.Rect(0, 0, width, height))        // 创建一个新RGBA图像
	draw.Draw(newImg, newImg.Bounds(), m, m.Bounds().Min, draw.Over) // 画上缩放后的图片

	/*
		// 保存文件
		imgfile, _ := os.Create("003.jpg")
		defer imgfile.Close()
		jpeg.Encode(imgfile, newImg, &jpeg.Options{100})
	*/

	emptyBuff := bytes.NewBuffer(nil)                  // 开辟一个新的空buff
	jpeg.Encode(emptyBuff, newImg, nil)                // img写入到buff
	dist := make([]byte, len(emptyBuff.Bytes())*2)     // 开辟存储空间
	base64.StdEncoding.Encode(dist, emptyBuff.Bytes()) // buff转成base64
	index := bytes.IndexByte(dist, 0)                  // 这里要注意，因为申请的固定长度数组，所以没有被填充完的部分需要去掉，负责输出可能出错
	baseImage := dist[0:index]

	return "data:image/jpg;base64," + string(baseImage), nil
}
