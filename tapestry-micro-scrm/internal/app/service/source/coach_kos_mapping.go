package source

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"micro-scrm/internal/app/models"
	sourcemodel "micro-scrm/internal/app/models/source"
	coachkos "micro-scrm/internal/pkg/sdk/coach_kos"
)

// CoachKOSGoodsTask 描述一条待补全商品信息的素材任务。
type CoachKOSGoodsTask struct {
	EnterpriseID int64
	SourceID     int64
	ProductCodes []string
}

type coachKOSGoodsSyncPlan struct {
	CreateProductCodes []string
	EnrichProductCodes []string
	DeleteIDs          []int64
}

func expandCoachKOSProductCodes(productIDs []string) []string {
	result := make([]string, 0, len(productIDs))
	seen := make(map[string]struct{})
	for _, productID := range productIDs {
		for _, part := range strings.Split(productID, "&") {
			productCode := strings.TrimSpace(part)
			if productCode == "" {
				continue
			}
			if _, ok := seen[productCode]; ok {
				continue
			}
			seen[productCode] = struct{}{}
			result = append(result, productCode)
		}
	}
	return result
}

func coachKOSSKUCode(productCode string) string {
	return strings.ReplaceAll(strings.TrimSpace(productCode), "_", " ")
}

func planCoachKOSRelationGoods(existing []sourcemodel.RelationGoodsModel, productCodes []string) coachKOSGoodsSyncPlan {
	plan := coachKOSGoodsSyncPlan{
		CreateProductCodes: make([]string, 0),
		EnrichProductCodes: make([]string, 0),
		DeleteIDs:          make([]int64, 0),
	}
	currentProductCodes := make(map[string]string, len(productCodes))
	for _, productCode := range productCodes {
		currentProductCodes[coachKOSSKUCode(productCode)] = productCode
	}

	existingSKUCodes := make(map[string]struct{}, len(existing))
	completeSKUCodes := make(map[string]struct{}, len(existing))
	for _, goods := range existing {
		skuCode := strings.TrimSpace(goods.SkuCode)
		if _, ok := currentProductCodes[skuCode]; !ok {
			plan.DeleteIDs = append(plan.DeleteIDs, goods.ID)
			continue
		}
		existingSKUCodes[skuCode] = struct{}{}
		if strings.TrimSpace(goods.SKU) != "" {
			completeSKUCodes[skuCode] = struct{}{}
		}
	}

	for _, productCode := range productCodes {
		skuCode := coachKOSSKUCode(productCode)
		if _, ok := existingSKUCodes[skuCode]; !ok {
			plan.CreateProductCodes = append(plan.CreateProductCodes, productCode)
		}
		if _, ok := completeSKUCodes[skuCode]; !ok {
			plan.EnrichProductCodes = append(plan.EnrichProductCodes, productCode)
		}
	}
	return plan
}

func serializeCoachKOSProductIDs(productIDs []string) string {
	if len(productIDs) == 0 {
		return "[]"
	}

	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)
	_ = encoder.Encode(productIDs)
	return strings.TrimSuffix(buffer.String(), "\n")
}

func coachKOSMonthGroupName(dataDate string) (string, error) {
	date, err := time.Parse("2006-01-02", dataDate)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d月小红书素材", date.Month()), nil
}

func coachKOSCreatedAt(dataDate string, fallback int64) int64 {
	shanghai := time.FixedZone("Asia/Shanghai", 8*60*60)
	date, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(dataDate), shanghai)
	if err != nil {
		return fallback
	}
	return date.Unix()
}

const coachKOSNegativeSentimentThreshold = 0.25

func coachKOSNegativeSentimentRate(raw json.RawMessage) float64 {
	data := bytes.TrimSpace(raw)
	if len(data) == 0 || data[0] != '{' {
		return 0
	}

	var sentiment map[string]float64
	if err := json.Unmarshal(data, &sentiment); err != nil {
		return 0
	}
	return sentiment["负面"]
}

func shouldDisableCoachKOSPost(post coachkos.PostResult) bool {
	if post.IsValid == "N" {
		return true
	}
	return coachKOSNegativeSentimentRate(post.CommentSentiment) > coachKOSNegativeSentimentThreshold
}

func truncateCoachKOSDescription(content string) string {
	runes := []rune(content)
	if len(runes) <= 255 {
		return content
	}
	return string(runes[:254]) + "…"
}

func buildCoachKOSModels(
	enterpriseID, sourceID, groupID, now int64,
	post coachkos.PostResult,
	materialCoverURL string,
) (sourcemodel.SourceModel, sourcemodel.LinkModel, sourcemodel.LinkKOSModel) {
	status := int32(sourcemodel.StatusEnable)
	if shouldDisableCoachKOSPost(post) {
		status = sourcemodel.StatusDisable
	}

	createdAt := coachKOSCreatedAt(post.DataDate, now)
	base := models.BaseModel{ID: sourceID, CreatedAt: createdAt, UpdatedAt: now}
	description := post.Content
	if strings.TrimSpace(description) == "" {
		description = post.Hashtag
	}
	sourceModel := sourcemodel.SourceModel{
		MicroBaseModel: models.MicroBaseModel{BaseModel: base, EnterpriseID: enterpriseID},
		GroupID:        groupID,
		Type:           sourcemodel.TypeLink,
		Status:         status,
		StartTime:      int(now),
		EndTime:        int(coachKOSEndTime()),
	}
	linkModel := sourcemodel.LinkModel{
		BaseModel:   base,
		Title:       post.Name,
		Description: truncateCoachKOSDescription(description),
		CoverURL:    materialCoverURL,
		LinkType:    1,
		LinkURL:     post.URL,
		Author:      post.InfluencerName,
		Content:     "",
	}

	productIDs := serializeCoachKOSProductIDs(post.ProductID)
	kosModel := sourcemodel.LinkKOSModel{
		MicroBaseModel:   models.MicroBaseModel{BaseModel: models.BaseModel{CreatedAt: now, UpdatedAt: now}, EnterpriseID: enterpriseID},
		SourceID:         sourceID,
		SourceLinkID:     sourceID,
		PostID:           post.PostID,
		SourceURL:        post.SourceURL,
		DataDate:         post.DataDate,
		PublishedAt:      post.PublishedAt,
		Name:             post.Name,
		Content:          post.Content,
		Hashtag:          post.Hashtag,
		URL:              post.URL,
		CoverURL:         post.CoverURL,
		InfluencerID:     post.InfluencerID,
		InfluencerName:   post.InfluencerName,
		EmployeeID:       post.EmployeeID,
		EmployeeTag:      post.EmployeeTag,
		StoreID:          post.StoreID,
		StoreName:        post.StoreName,
		ProductID:        productIDs,
		ContentTag:       post.ContentTag,
		LikeCount:        post.LikeCount,
		CommentCount:     post.CommentCount,
		FavoriteCount:    post.FavoriteCount,
		ShareCount:       post.ShareCount,
		InteractionCount: post.InteractionCount,
		IsValid:          post.IsValid,
		CommentSentiment: string(post.CommentSentiment),
		PriceMentioned:   post.PriceMentioned,
		APIUpdatedAt:     post.UpdatedAt,
		SourceUpdatedAt:  post.SourceUpdatedAt,
	}
	return sourceModel, linkModel, kosModel
}

func buildCoachKOSRelationGoods(enterpriseID, sourceID, now int64, productCodes []string) []sourcemodel.RelationGoodsModel {
	goods := make([]sourcemodel.RelationGoodsModel, 0, len(productCodes))
	for _, productCode := range productCodes {
		goods = append(goods, sourcemodel.RelationGoodsModel{
			MicroBaseModel: models.MicroBaseModel{
				BaseModel:    models.BaseModel{CreatedAt: now, UpdatedAt: now},
				EnterpriseID: enterpriseID,
			},
			SourceID: sourceID,
			Type:     sourcemodel.TypeLink,
			SkuCode:  coachKOSSKUCode(productCode),
			AddType:  2,
		})
	}
	return goods
}

func coachKOSEndTime() int64 {
	shanghai := time.FixedZone("Asia/Shanghai", 8*60*60)
	return time.Date(2035, 12, 31, 23, 59, 59, 0, shanghai).Unix()
}
