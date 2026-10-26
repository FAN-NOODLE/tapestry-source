package elastic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"

	"micro-scrm/internal/config"

	"micro-scrm/internal/modules/customer"
	"micro-scrm/internal/pkg/elastic"
	"micro-scrm/internal/pkg/errorIgnore"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
)

type EsRepo struct {
	cli *elastic.Client
	env string // prod /utc /test
}

type EsAddResIndex struct {
	Index struct {
		ID          string `json:"_id"`
		Index       string `json:"_index"`
		PrimaryTerm int    `json:"_primary_term"`
		SeqNo       int    `json:"_seq_no"`
		Shards      struct {
			Failed     int `json:"failed"`
			Successful int `json:"successful"`
			Total      int `json:"total"`
		} `json:"_shards"`
		Type          string      `json:"_type"`
		Version       int         `json:"_version"`
		ForcedRefresh bool        `json:"forced_refresh"`
		Result        string      `json:"result"`
		Status        int         `json:"status"`
		Error         interface{} `json:"error"`
	} `json:"index"`
}

func makeFilter(cond *customer.ListCond) []interface{} {
	var filter []interface{}
	var must []interface{}
	var mustNot []interface{}
	var should []interface{}
	if len(cond.NickName) > 0 {
		must = append(must, map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []interface{}{
					map[string]interface{}{"match_phrase": map[string]string{"name": cond.NickName}},
					map[string]interface{}{"match_phrase": map[string]string{"nick": cond.NickName}},
					map[string]interface{}{"term": map[string]string{"name.keyword": cond.NickName}},
					map[string]interface{}{"term": map[string]string{"nick.keyword": cond.NickName}},
				},
			},
		})
	}
	if len(cond.Mobile) > 0 {
		// must = append(must, map[string]interface{}{"match": map[string]string{"mobile": cond.Mobile}})
		must = append(must, map[string]interface{}{"wildcard": map[string]string{"mobile": "*" + cond.Mobile + "*"}})
	}

	if len(cond.Remark.Remark) > 0 {
		if cond.Remark.Include {
			must = append(must, map[string]interface{}{"match": map[string]string{"remark": "*" + cond.Remark.Remark + "*"}})
		} else {
			mustNot = append(mustNot, map[string]interface{}{"match": map[string]string{"remark": "*" + cond.Remark.Remark + "*"}})
		}
	}

	// todo 过滤空头像,影响了整个的统计数，去除
	// mustNot = append(mustNot, gin.H{"query_string": gin.H{"default_field": "avatar.keyword", "query": "\"\""}})

	// 过滤是否会员(0,不过滤,1会员,2非会员)
	if cond.UnionCode > 0 {
		if cond.UnionCode == 1 {
			mustNot = append(mustNot, gin.H{"query_string": gin.H{"default_field": "union_code.keyword", "query": "\"\""}})
		} else {
			must = append(must, gin.H{"query_string": gin.H{"default_field": "union_code.keyword", "query": "\"\""}})
		}
	}

	if len(cond.TagCond.TagIDList) > 0 {
		if cond.TagCond.Include {
			filter = append(filter, map[string]interface{}{"terms": map[string]interface{}{"tags": cond.TagCond.TagIDList}})
		} else {
			mustNot = append(mustNot, map[string]interface{}{"terms": map[string]interface{}{"tags": cond.TagCond.TagIDList}})
		}
	}

	if len(cond.CustomerGroupCond.GroupIDList) > 0 {
		if cond.CustomerGroupCond.Include {
			filter = append(filter, map[string]interface{}{"terms": map[string]interface{}{"groups": cond.CustomerGroupCond.GroupIDList}})
		} else {
			mustNot = append(mustNot, map[string]interface{}{"terms": map[string]interface{}{"groups": cond.CustomerGroupCond.GroupIDList}})

		}
	}

	if cond.MarketTimeEnd > 0 && cond.MarketTimeStart > 0 {
		filter = append(filter, map[string]interface{}{
			"range": map[string]interface{}{
				"latest_market_time": map[string]interface{}{"gte": cond.MarketTimeStart, "lte": cond.MarketTimeEnd},
			}})
	} else if cond.MarketTimeStart > 0 {
		filter = append(filter, map[string]interface{}{
			"range": map[string]interface{}{
				"latest_market_time": map[string]interface{}{"gte": cond.MarketTimeStart},
			}})
	} else if cond.MarketTimeEnd > 0 {
		filter = append(filter, map[string]interface{}{
			"range": map[string]interface{}{
				"latest_market_time": map[string]interface{}{"lte": cond.MarketTimeEnd},
			}})
	}

	if len(cond.MemberLevelCond.LevelIDList) > 0 {
		if cond.MemberLevelCond.Include {
			filter = append(filter, map[string]interface{}{"terms": map[string]interface{}{"grade.keyword": cond.MemberLevelCond.LevelIDList}})
		} else {
			mustNot = append(mustNot, map[string]interface{}{"terms": map[string]interface{}{"grade.keyword": cond.MemberLevelCond.LevelIDList}})
		}
	}
	if len(cond.Eids) > 0 {
		filter = append(filter, map[string]interface{}{
			"terms": map[string]interface{}{
				"id": cond.Eids,
			}})
	}
	if cond.CustomerType > 0 {
		switch cond.CustomerType {
		case customer.CustomerTypeShop: // 商城用户
			must = append(must, map[string]interface{}{"exists": map[string]string{"field": "wx_shop.app_id"}})
		case customer.CustomerTypeExternal: // 外部联系人用户
			must = append(must, map[string]interface{}{"exists": map[string]string{"field": "wx_external.app_id"}})
		case customer.CustomerTypeOpen: // 公众号用户
			must = append(must, map[string]interface{}{"exists": map[string]string{"field": "wx_open.app_id"}})
		}
	} else if len(cond.CustomerSource) > 0 {
		for _, customerType := range cond.CustomerSource {
			switch customerType {
			case customer.CustomerTypeShop:
				should = append(should, map[string]interface{}{"exists": map[string]string{"field": "wx_shop.app_id"}})
			case customer.CustomerTypeExternal: // 外部联系人用户
				should = append(should, map[string]interface{}{"exists": map[string]string{"field": "wx_external.app_id"}})
			case customer.CustomerTypeOpen: // 公众号用户
				should = append(should, map[string]interface{}{"exists": map[string]string{"field": "wx_open.app_id"}})
			case customer.CustomerTypeUnion:
				should = append(should, map[string]interface{}{"exists": map[string]string{"field": "wx_union"}})
			case customer.CustomerTypeMobile:
				should = append(should, map[string]interface{}{
					"bool": map[string]interface{}{
						"must": []interface{}{
							map[string]interface{}{"exists": map[string]string{"field": "mobile"}},
						},
						"must_not": []interface{}{map[string]interface{}{"term": map[string]string{"mobile.keyword": ""}}},
					},
				})
			case customer.CustomerTypeDeliveryMobile:
				should = append(should, map[string]interface{}{
					"bool": map[string]interface{}{
						"must": []interface{}{
							map[string]interface{}{"exists": map[string]string{"field": "delivery_mobile"}},
						},
						"must_not": []interface{}{map[string]interface{}{"term": map[string]string{"delivery_mobile.keyword": ""}}},
					},
				})
			}
		}
	}

	filterBool := make(map[string]interface{})
	if len(must) > 0 {
		filterBool["must"] = must
	}
	if len(mustNot) > 0 {
		filterBool["must_not"] = mustNot
	}

	if len(should) > 0 {
		filterBool["should"] = should
	}

	if len(filterBool) > 0 {
		filter = append(filter, map[string]interface{}{"bool": filterBool})
	}

	return filter
}

// QueryUserList 查询用户列表
func (es *EsRepo) QueryUserList(c ctx.Context, cond *customer.ListRequest) (int64, []*customer.EnterpriseUser, error) {
	userList := []*customer.EnterpriseUser{}

	index := fmt.Sprintf("micro_scrm_enterprise_%d.enterprise_user", cond.EnterpriseID)
	index = GetEsIndexName(index)

	filter := makeFilter(&cond.ListCond)
	query := make(map[string]interface{})
	if len(filter) > 0 {
		query["query"] = map[string]interface{}{
			"bool": map[string]interface{}{
				"filter": filter,
			},
		}
	}

	var sortArr []interface{}
	// if cond.SortKey != "" {
	//	keysort := map[string]string{
	//		cond.SortKey: cond.SortType,
	//	}
	//	sortArr = append(sortArr, keysort)
	// }
	timeSort := map[string]string{
		"created_at": "desc",
	}
	sortArr = append(sortArr, timeSort)
	query["sort"] = sortArr
	queryj, _ := json.Marshal(query)
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(query)
	esc := es.cli.GetEsClient()
	res, err := esc.Search(
		esc.Search.WithContext(c),
		esc.Search.WithIndex(index),
		esc.Search.WithBody(&buf),
		esc.Search.WithFrom((cond.Page-1)*cond.Size),
		esc.Search.WithSize(cond.Size),
		esc.Search.WithTrackTotalHits(true),
		esc.Search.WithRestTotalHitsAsInt(true),
	)
	log.Info(c).Msgf("QueryUserList query：%s, %s,res %s", index, string(queryj), res.String())

	if errorIgnore.Ignore(err) != nil {
		return 0, nil, err
	}
	if len(res.String()) > 0 {
		resMap := make(map[string]interface{})
		if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
			log.Error(c).Msgf("Error parsing the response body: %s", err)
			return 0, nil, err
		}

		resBody, _ := json.Marshal(resMap)
		esres := customer.EsQueryRes{}
		err = json.Unmarshal(resBody, &esres)
		if err != nil || len(esres.Hits.Hits) == 0 {
			return 0, nil, err
		}

		for _, hit := range esres.Hits.Hits {
			source := hit.Source
			user := &customer.EnterpriseUser{}
			user.ID = source.ID

			if len(source.WxUnion) > 0 {
				user.WxUnion = make(map[string]string)
				for _, record := range source.WxUnion {
					user.WxUnion[record.Pid] = record.UnionID
					if record.Pid == "default" {
						user.UnionId = record.UnionID
					}
				}
			}
			if len(source.Groups) > 0 {
				for _, group := range source.Groups {
					if ftag, ok := group.(float64); ok {
						user.Groups = append(user.Groups, int64(ftag))
					}
					if _, ok := group.(int64); ok {
						user.Groups = append(user.Groups, group.(int64))
					}
					// user.Groups = append(user.Groups, group.(int64))
				}
			}
			if len(source.Tags) > 0 {
				for _, tag := range source.Tags {
					if ftag, ok := tag.(float64); ok {
						user.Tags = append(user.Tags, int64(ftag))
					}
					if _, ok := tag.(int64); ok {
						user.Tags = append(user.Tags, tag.(int64))
					}
				}
			}

			user.Avatar = source.Avatar
			user.Birthday = source.Birthday
			user.UnionCode = source.UnionCode
			user.Point = source.Point
			user.Grade = source.Grade

			user.City = source.City
			user.Country = source.Country
			user.Language = source.Language
			user.Mobile = source.Mobile
			user.DeliveryMobile = source.DeliveryMobile
			user.Name = source.Name
			user.Nick = source.Nick
			user.PayOrder = source.PayOrder
			user.LastMarketTime = int64(source.LatestMarketTime)
			user.Province = source.Province
			user.RefundFee = source.RefundFee
			user.RefundOrder = source.RefundOrder
			user.Remark = source.Remark
			user.TotalFee = source.TotalFee
			user.DtUpdate = source.DtUpdate

			user.Sex = source.Sex
			if len(source.WxExternal) > 0 {
				user.WxExternal = make(map[string]*customer.ThirdUserID)
				for _, record := range source.WxExternal {
					user.WxExternal[record.AppID] = record.ThirdUser
				}
			}
			if len(source.WxOpen) > 0 {
				user.WxOpen = make(map[string]*customer.ThirdUserID)
				for _, record := range source.WxOpen {
					user.WxOpen[record.AppID] = record.ThirdUser
				}
			}
			if len(source.WxShop) > 0 {
				user.WxShop = make(map[string]*customer.ThirdUserID)
				for _, record := range source.WxShop {
					user.WxShop[record.AppID] = record.ThirdUser
				}
			}
			userList = append(userList, user)
		}
		// userListj, _ := json.Marshal(userList)
		// log.Info(c).Msgf("QueryUserList Get Result:query=%s,result=%s", string(queryj), string(userListj))

		return esres.Hits.Total, userList, err
	}
	return 0, nil, nil
}

func (es *EsRepo) QueryByUnionID(c ctx.Context, cond *customer.CondUnionID) (int64, []*customer.EnterpriseUser, error) {
	userList := []*customer.EnterpriseUser{}

	index := fmt.Sprintf("micro_scrm_enterprise_%d.enterprise_user", cond.EnterpriseID)
	index = GetEsIndexName(index)

	query := make(map[string]interface{})
	var should []interface{}
	if len(cond.UnionId) > 0 {
		for _, v := range cond.UnionId {
			should = append(should, map[string]interface{}{
				"match": map[string]interface{}{
					"wx_union.union_id": v,
				},
			})
		}
	}

	if len(cond.EUid) > 0 {
		for _, v := range cond.EUid {
			should = append(should, map[string]interface{}{
				"match": map[string]interface{}{
					"_id": v,
				},
			})
		}
	}

	query["query"] = map[string]interface{}{
		"bool": map[string]interface{}{
			"should": should,
		},
	}
	queryj, _ := json.Marshal(query)
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(query)
	esc := es.cli.GetEsClient()
	res, err := esc.Search(
		esc.Search.WithContext(c),
		esc.Search.WithIndex(index),
		esc.Search.WithBody(&buf),
		esc.Search.WithTrackTotalHits(true),
		esc.Search.WithRestTotalHitsAsInt(true),
	)
	log.Info(c).Msgf("QueryUserList query：%s, %s,res %s", index, string(queryj), res.String())

	if errorIgnore.Ignore(err) != nil {
		return 0, nil, err
	}
	if len(res.String()) > 0 {
		resMap := make(map[string]interface{})
		if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
			log.Error(c).Msgf("Error parsing the response body: %s", err)
			return 0, nil, err
		}

		resBody, _ := json.Marshal(resMap)
		esres := customer.EsQueryRes{}
		err = json.Unmarshal(resBody, &esres)
		if err != nil || len(esres.Hits.Hits) == 0 {
			return 0, nil, err
		}

		for _, hit := range esres.Hits.Hits {
			source := hit.Source
			user := &customer.EnterpriseUser{}
			user.ID = source.ID

			if len(source.WxUnion) > 0 {
				user.WxUnion = make(map[string]string)
				for _, record := range source.WxUnion {
					user.WxUnion[record.Pid] = record.UnionID
					user.UnionId = record.UnionID
				}
			}
			if len(source.Groups) > 0 {
				for _, group := range source.Groups {
					if ftag, ok := group.(float64); ok {
						user.Groups = append(user.Groups, int64(ftag))
					}
					if _, ok := group.(int64); ok {
						user.Groups = append(user.Groups, group.(int64))
					}
					// user.Groups = append(user.Groups, group.(int64))
				}
			}
			if len(source.Tags) > 0 {
				for _, tag := range source.Tags {
					if ftag, ok := tag.(float64); ok {
						user.Tags = append(user.Tags, int64(ftag))
					}
					if _, ok := tag.(int64); ok {
						user.Tags = append(user.Tags, tag.(int64))
					}
				}
			}

			user.Avatar = source.Avatar
			user.Birthday = source.Birthday
			user.UnionCode = source.UnionCode
			user.Point = source.Point
			user.Grade = source.Grade

			user.City = source.City
			user.Country = source.Country
			user.Language = source.Language
			user.Mobile = source.Mobile
			user.DeliveryMobile = source.DeliveryMobile
			user.Name = source.Name
			user.Nick = source.Nick
			user.PayOrder = source.PayOrder
			user.LastMarketTime = int64(source.LatestMarketTime)
			user.Province = source.Province
			user.RefundFee = source.RefundFee
			user.RefundOrder = source.RefundOrder
			user.Remark = source.Remark
			user.TotalFee = source.TotalFee
			user.DtUpdate = source.DtUpdate

			user.Sex = source.Sex
			if len(source.WxExternal) > 0 {
				user.WxExternal = make(map[string]*customer.ThirdUserID)
				for _, record := range source.WxExternal {
					user.WxExternal[record.AppID] = record.ThirdUser
				}
			}
			if len(source.WxOpen) > 0 {
				user.WxOpen = make(map[string]*customer.ThirdUserID)
				for _, record := range source.WxOpen {
					user.WxOpen[record.AppID] = record.ThirdUser
				}
			}
			if len(source.WxShop) > 0 {
				user.WxShop = make(map[string]*customer.ThirdUserID)
				for _, record := range source.WxShop {
					user.WxShop[record.AppID] = record.ThirdUser
				}
			}
			userList = append(userList, user)
		}
		// userListj, _ := json.Marshal(userList)
		// log.Info(c).Msgf("QueryUserList Get Result:query=%s,result=%s", string(queryj), string(userListj))

		return esres.Hits.Total, userList, err
	}
	return 0, nil, nil
}

func (es *EsRepo) QueryByEUidPag(c ctx.Context, cond *customer.CondUnionID, form, size int) (int64, []*customer.EnterpriseUser, error) {
	userList := []*customer.EnterpriseUser{}

	index := fmt.Sprintf("micro_scrm_enterprise_%d.enterprise_user", cond.EnterpriseID)
	index = GetEsIndexName(index)

	query := make(map[string]interface{})
	var should []interface{}
	if len(cond.UnionId) > 0 {
		for _, v := range cond.UnionId {
			should = append(should, map[string]interface{}{
				"match": map[string]interface{}{
					"wx_union.union_id": v,
				},
			})
		}
	}

	if len(cond.EUid) > 0 {
		for _, v := range cond.EUid {
			should = append(should, map[string]interface{}{
				"match": map[string]interface{}{
					"_id": v,
				},
			})
		}
	}

	query["query"] = map[string]interface{}{
		"bool": map[string]interface{}{
			"should": should,
		},
	}
	queryj, _ := json.Marshal(query)
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(query)
	esc := es.cli.GetEsClient()
	res, err := esc.Search(
		esc.Search.WithContext(c),
		esc.Search.WithIndex(index),
		esc.Search.WithBody(&buf),
		esc.Search.WithFrom(form),
		esc.Search.WithSize(size),
		esc.Search.WithTrackTotalHits(true),
		esc.Search.WithRestTotalHitsAsInt(true),
	)
	log.Info(c).Msgf("QueryUserList query：%s, %s,res %s", index, string(queryj), res.String())

	if errorIgnore.Ignore(err) != nil {
		return 0, nil, err
	}
	if len(res.String()) > 0 {
		resMap := make(map[string]interface{})
		if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
			log.Error(c).Msgf("Error parsing the response body: %s", err)
			return 0, nil, err
		}

		resBody, _ := json.Marshal(resMap)
		esres := customer.EsQueryRes{}
		err = json.Unmarshal(resBody, &esres)
		if err != nil || len(esres.Hits.Hits) == 0 {
			return 0, nil, err
		}

		for _, hit := range esres.Hits.Hits {
			source := hit.Source
			user := &customer.EnterpriseUser{}
			user.ID = source.ID

			if len(source.WxUnion) > 0 {
				user.WxUnion = make(map[string]string)
				for _, record := range source.WxUnion {
					user.WxUnion[record.Pid] = record.UnionID
					user.UnionId = record.UnionID
				}
			}
			if len(source.Groups) > 0 {
				for _, group := range source.Groups {
					if ftag, ok := group.(float64); ok {
						user.Groups = append(user.Groups, int64(ftag))
					}
					if _, ok := group.(int64); ok {
						user.Groups = append(user.Groups, group.(int64))
					}
					// user.Groups = append(user.Groups, group.(int64))
				}
			}
			if len(source.Tags) > 0 {
				for _, tag := range source.Tags {
					if ftag, ok := tag.(float64); ok {
						user.Tags = append(user.Tags, int64(ftag))
					}
					if _, ok := tag.(int64); ok {
						user.Tags = append(user.Tags, tag.(int64))
					}
				}
			}

			user.Avatar = source.Avatar
			user.Birthday = source.Birthday
			user.UnionCode = source.UnionCode
			user.Point = source.Point
			user.Grade = source.Grade

			user.City = source.City
			user.Country = source.Country
			user.Language = source.Language
			user.Mobile = source.Mobile
			user.DeliveryMobile = source.DeliveryMobile
			user.Name = source.Name
			user.Nick = source.Nick
			user.PayOrder = source.PayOrder
			user.LastMarketTime = int64(source.LatestMarketTime)
			user.Province = source.Province
			user.RefundFee = source.RefundFee
			user.RefundOrder = source.RefundOrder
			user.Remark = source.Remark
			user.TotalFee = source.TotalFee
			user.DtUpdate = source.DtUpdate

			user.Sex = source.Sex
			if len(source.WxExternal) > 0 {
				user.WxExternal = make(map[string]*customer.ThirdUserID)
				for _, record := range source.WxExternal {
					user.WxExternal[record.AppID] = record.ThirdUser
				}
			}
			if len(source.WxOpen) > 0 {
				user.WxOpen = make(map[string]*customer.ThirdUserID)
				for _, record := range source.WxOpen {
					user.WxOpen[record.AppID] = record.ThirdUser
				}
			}
			if len(source.WxShop) > 0 {
				user.WxShop = make(map[string]*customer.ThirdUserID)
				for _, record := range source.WxShop {
					user.WxShop[record.AppID] = record.ThirdUser
				}
			}
			userList = append(userList, user)
		}
		// userListj, _ := json.Marshal(userList)
		// log.Info(c).Msgf("QueryUserList Get Result:query=%s,result=%s", string(queryj), string(userListj))

		return esres.Hits.Total, userList, err
	}
	return 0, nil, nil
}

// QueryUserCount 计算标签分组人数
func (es *EsRepo) QueryUserCount(c ctx.Context, cond *customer.ListCond) (*customer.TagUserCount, error) {
	query := makeCountQuery(cond, false)
	queryj, _ := json.Marshal(query)

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(query)
	esc := es.cli.GetEsClient()

	index := fmt.Sprintf("micro_scrm_enterprise_%d.enterprise_user", cond.EnterpriseID)
	index = GetEsIndexName(index)
	res, err := esc.Search(
		esc.Search.WithContext(c),
		esc.Search.WithIndex(index),
		esc.Search.WithBody(&buf),
	)
	resMap := make(map[string]interface{})
	if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
		log.Error(c).Msgf("Error parsing the response body: %s", err)
		return nil, err
	}
	log.Info(c).Msgf("QueryUserCount query %s,res %s", string(queryj), res.String())
	if aggs, found := resMap["aggregations"]; found {
		aggsJ, err := json.Marshal(aggs)
		if err != nil {
			log.Error(c).Msgf("aggregations Marshal err: %s", err)
			return nil, err
		}
		ta := customer.TagUserCountAggs{}
		err = json.Unmarshal(aggsJ, &ta)
		if err != nil {
			log.Error(c).Msgf("aggregations Unmarshal TagUserCountAggs err: %s", err)
			return nil, err
		} else {
			ct := customer.TagUserCount{}
			ct.MemberCount = ta.MemberCount.Value
			ct.MobileUserCount = ta.MobileUserCount.DocCount
			ct.QwUserCount = ta.QwUserCount.DocCount
			ct.ShopUserCount = ta.ShopUserCount.DocCount
			ct.WeChatOfficialAccountCount = ta.ChatOfficialAccountCount.DocCount
			return &ct, nil
		}
	}
	return nil, err
}

// QueryUserCountByGroup 计算标签用户分组人数
func (es *EsRepo) QueryUserCountByGroup(c ctx.Context, cond *customer.ListCond) (map[string]*customer.TagUserCount, error) {
	query := makeCountQuery(cond, true)
	queryj, _ := json.Marshal(query)

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(query)
	esc := es.cli.GetEsClient()

	index := fmt.Sprintf("micro_scrm_enterprise_%d.enterprise_user", cond.EnterpriseID)
	index = GetEsIndexName(index)
	res, err := esc.Search(
		esc.Search.WithContext(c),
		esc.Search.WithIndex(index),
		esc.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, err
	}

	var resValue struct {
		Aggs struct {
			ChatOfficialAccountCount struct {
				Buckets []struct {
					Key      string `json:"key"`
					DocCount int64  `json:"doc_count"`
				} `json:"buckets"`
			} `json:"chat_official_account_count"`
			MemberCount struct {
				Buckets []struct {
					Key      string `json:"key"`
					DocCount int64  `json:"doc_count"`
				} `json:"buckets"`
			} `json:"member_count"`
			MobileUserCount struct {
				Buckets []struct {
					Key      string `json:"key"`
					DocCount int64  `json:"doc_count"`
				} `json:"buckets"`
			} `json:"mobile_user_count"`
			QwUserCount struct {
				Buckets []struct {
					Key      string `json:"key"`
					DocCount int64  `json:"doc_count"`
				} `json:"buckets"`
			} `json:"qw_user_count"`
			ShopUserCount struct {
				Buckets []struct {
					Key      string `json:"key"`
					DocCount int64  `json:"doc_count"`
				} `json:"buckets"`
			} `json:"shop_user_count"`
		} `json:"aggregations"`
	}
	if err := json.NewDecoder(res.Body).Decode(&resValue); err != nil {
		log.Error(c).Msgf("Error parsing the response body: %s", err)
		return nil, err
	}

	result := make(map[string]*customer.TagUserCount)
	for _, v := range resValue.Aggs.ChatOfficialAccountCount.Buckets {
		cnt, ok := result[v.Key]
		if !ok {
			cnt = &customer.TagUserCount{}
		}
		cnt.WeChatOfficialAccountCount = v.DocCount
	}

	for _, v := range resValue.Aggs.ShopUserCount.Buckets {
		cnt, ok := result[v.Key]
		if !ok {
			cnt = &customer.TagUserCount{}
		}
		cnt.ShopUserCount = v.DocCount
	}

	for _, v := range resValue.Aggs.QwUserCount.Buckets {
		cnt, ok := result[v.Key]
		if !ok {
			cnt = &customer.TagUserCount{}
		}
		cnt.QwUserCount = v.DocCount
	}

	for _, v := range resValue.Aggs.MobileUserCount.Buckets {
		cnt, ok := result[v.Key]
		if !ok {
			cnt = &customer.TagUserCount{}
		}
		cnt.MobileUserCount = v.DocCount
	}

	for _, v := range resValue.Aggs.MemberCount.Buckets {
		cnt, ok := result[v.Key]
		if !ok {
			cnt = &customer.TagUserCount{}
		}
		cnt.MemberCount = v.DocCount
	}
	log.Info(c).Msgf("QueryUserCount query %s,res %s", string(queryj), res.String())
	return result, nil
}

func makeCountQuery(cond *customer.ListCond, withGroup bool) map[string]interface{} {
	filter := makeFilter(cond)
	query := make(map[string]interface{})
	if len(filter) > 0 {
		query["query"] = map[string]interface{}{
			"bool": map[string]interface{}{
				"filter": filter,
			},
		}
	}

	aggs := makeAggs()
	if withGroup {
		for k := range aggs {
			aggs[k]["terms"] = map[string]interface{}{"field": "groups"}
		}
	}

	query["aggs"] = aggs
	return query
}

func makeAggs() map[string]map[string]interface{} {
	aggs := make(map[string]map[string]interface{})
	aggs["member_count"] = map[string]interface{}{
		"value_count": map[string]string{
			"field": "_id",
		},
	}

	aggs["mobile_user_count"] = map[string]interface{}{
		"filter": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []interface{}{
					map[string]interface{}{
						"bool": map[string]interface{}{
							"must": []interface{}{
								map[string]interface{}{"exists": map[string]string{"field": "mobile"}},
							},
							"must_not": []interface{}{map[string]interface{}{"term": map[string]string{"mobile.keyword": ""}}},
						},
					},
					map[string]interface{}{
						"bool": map[string]interface{}{
							"must": []interface{}{
								map[string]interface{}{"exists": map[string]string{"field": "delivery_mobile"}},
							},
							"must_not": []interface{}{map[string]interface{}{"term": map[string]string{"delivery_mobile.keyword": ""}}},
						},
					},
				},
			},
		},
	}

	aggs["qw_user_count"] = map[string]interface{}{
		"filter": map[string]interface{}{"exists": map[string]string{"field": "wx_external.app_id"}},
	}
	aggs["shop_user_count"] = map[string]interface{}{
		"filter": map[string]interface{}{"exists": map[string]string{"field": "wx_shop.app_id"}},
	}
	aggs["chat_official_account_count"] = map[string]interface{}{
		"filter": map[string]interface{}{"exists": map[string]string{"field": "wx_open.app_id"}},
	}
	return aggs
}

// es 索引是否存在
func (es *EsRepo) ExistsIndex(cc ctx.Context, index string) bool {
	req := esapi.IndicesExistsRequest{Index: []string{index}}
	res, err := req.Do(cc, elastic.ShareEsClient.GetEsClient())
	log.Info(cc).Msgf("ExistsIndex -> index:%s,res:%s ,err:%v", index, res.String(), err)
	if res.StatusCode == 200 {
		return true
	}
	return false
}

// es 索引删除
func (es *EsRepo) DelsIndex(cc ctx.Context, index string) bool {
	req := esapi.IndicesDeleteRequest{Index: []string{index}}
	res, err := req.Do(cc, elastic.ShareEsClient.GetEsClient())
	log.Info(cc).Msgf("DelsIndex -> index:%s,res:%s ,err:%v", index, res.String(), err)
	if res.StatusCode == 200 {
		return true
	}
	return false
}

func (es *EsRepo) CreateIndex(cc ctx.Context, index, indexBody string) error {
	req := esapi.IndicesCreateRequest{
		Index: index,
		Body:  strings.NewReader(indexBody),
	}

	res, err := req.Do(cc, elastic.ShareEsClient.GetEsClient())
	if err != nil {
		log.Info(cc).Msgf("CreateIndex fail -> index:%s,doc:%s ,err:%v", index, indexBody, err)
		return err

	}
	defer res.Body.Close()
	if res.StatusCode == 200 {
		log.Info(cc).Msgf("CreateIndex success -> index:%s,doc:%s ,\n res:%s", index, indexBody, res.String())
	} else {
		log.Error(cc).Msgf("CreateIndex fail -> index:%s,doc:%s , \n res:%s", index, indexBody, res.String())
	}

	return err

}

// 添加或更新es文档
func (es *EsRepo) AddEsIndexDoc(cc ctx.Context, index, id, doc string) error {
	req := esapi.IndexRequest{
		Index:      index,                  // Index name
		Body:       strings.NewReader(doc), // Document body
		DocumentID: id,                     // Document ID
		Refresh:    "true",                 // Refresh
	}
	res, err := req.Do(cc, elastic.ShareEsClient.GetEsClient())
	if err != nil {
		log.Error(cc).Msgf("AddEsIndexDoc fail -> index:%s,id:%s, doc:%s ,err:%v", index, id, doc, err)
		return err

	}
	defer res.Body.Close()

	if res.StatusCode == 201 || res.StatusCode == 200 {
		log.Info(cc).Msgf("AddEsIndexDoc success -> index:%s,id:%s, doc:%s, \n res:%s", index, id, doc, res.String())
	} else {
		log.Error(cc).Msgf("AddEsIndexDoc fail -> index:%s,id:%s, doc:%s ,\n res:%s", index, id, doc, res.String())
	}

	return err
}

func (es *EsRepo) AddDocsByBulk(cc ctx.Context, doc string) (map[string]map[string]int, error) {
	req := esapi.BulkRequest{
		Body:    strings.NewReader(doc),
		Refresh: "false",
		Timeout: 100,
	}
	res, err := req.Do(cc, elastic.ShareEsClient.GetEsClient())
	doc = strings.Replace(doc, "\n", "", -1)
	if err != nil {
		log.Error(cc).Msgf("AddDocsByBulk fail ->  doc:%s ,err:%v", doc, err)
		return nil, err

	}
	defer res.Body.Close()

	if res.StatusCode == 201 || res.StatusCode == 200 {
		// log.Info(cc).Msgf("AddDocsByBulk success ->doc:%s,  res:%s", doc, res.String())
	} else {
		log.Error(cc).Msgf("AddDocsByBulk fail -> doc:%s , res:%s", doc, res.String())
	}

	resMap := make(map[string]interface{})
	if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
		log.Error(cc).Msgf("Error parsing the response body: %s", err)
		return nil, err
	}

	var reslist []EsAddResIndex
	for k, v := range resMap {
		if k == "items" {
			var resrow []EsAddResIndex
			jsonv, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(jsonv, &resrow)
			if err != nil {
				return nil, err
			}
			reslist = append(reslist, resrow...)
		}
	}
	mapres := make(map[string]map[string]int)
	for _, v := range reslist {
		if v.Index.Error == nil {
			if _, ok := mapres[v.Index.Index]; !ok {
				mapres[v.Index.Index] = make(map[string]int)
			}
			mapres[v.Index.Index][v.Index.ID] = 1
		}
	}
	return mapres, err
}

// 删除文档
func (es *EsRepo) DeleteDoc(cc ctx.Context, index, id string) error {
	req := esapi.DeleteRequest{
		Index:      index,
		DocumentID: id,
	}

	res, err := req.Do(cc, elastic.ShareEsClient.GetEsClient())
	if err != nil {
		log.Info(cc).Msgf("DeleteDoc fail -> index:%s,id:%s,err:%v", index, id, err)
		return err
	}
	defer res.Body.Close()
	if res.StatusCode == 200 {
		log.Info(cc).Msgf("DeleteDoc success -> index:%s,id:%s,res :%v", index, id, res.String())
	} else {
		log.Error(cc).Msgf("DeleteDoc fail -> index:%s,id:%s,res :%v", index, id, res.String())
	}
	return nil
}

// 查询文档
func (es *EsRepo) SearchDoc(cc ctx.Context, index, body string) (string, error) {
	req := esapi.SearchRequest{
		Index: []string{index},
		Body:  strings.NewReader(body),
	}
	res, err := req.Do(cc, elastic.ShareEsClient.GetEsClient())
	if err != nil {
		log.Info(cc).Msgf("SearchDoc fail -> index:%s,body:%s,err:%v", index, body, err)
		return "", err
	}
	defer res.Body.Close()

	log.Info(cc).Msgf("SearchDoc success -> index:%s,body:%s,res:%s", index, body, res.String())
	return res.String(), nil
}

// 查询文档数量
func (es *EsRepo) CountDoc(cc ctx.Context, index, body string) (string, error) {
	req := esapi.CountRequest{
		Index: []string{index},
		Body:  strings.NewReader(body),
	}
	res, err := req.Do(cc, elastic.ShareEsClient.GetEsClient())
	if err != nil {
		log.Info(cc).Msgf("CountDoc fail -> index:%s,body:%s,err:%v", index, body, err)
		return "", err
	}
	defer res.Body.Close()

	log.Info(cc).Msgf("CountDoc success -> index:%s,body:%s,res:%s", index, body, res.String())
	return res.String(), nil
}

func GetEsIndexName(index string) string {
	prefix := config.CfgData.Env + "."
	return prefix + index
}
