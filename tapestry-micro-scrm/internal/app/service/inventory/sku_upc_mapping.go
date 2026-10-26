package inventorySvc

import (
	"micro-scrm/internal/app"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	sourceSvc "micro-scrm/internal/app/service/source"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
)

type skuUpcBatchReader interface {
	GetEnabledActiveBatch(c app.Ctx, enterpriseId int64) (microScrmModel.InventorySkuUpcImportBatchModel, error)
}

type skuUpcMappingReader interface {
	GetUpcCodesBySkuItem(c app.Ctx, enterpriseId int64, batchNo string, skuItem string) ([]string, error)
}

type skuUpcMappingResolver struct {
	batches  skuUpcBatchReader
	mappings skuUpcMappingReader
}

func newSkuUpcMappingResolver() skuUpcMappingResolver {
	return skuUpcMappingResolver{
		batches:  microScrmDao.NewInventorySkuUpcImportBatchDao(),
		mappings: microScrmDao.NewInventorySkuUpcMappingDao(),
	}
}

func (r skuUpcMappingResolver) Resolve(c app.Ctx, enterpriseId int64, skuItem string) ([]string, bool, error) {
	batch, err := r.batches.GetEnabledActiveBatch(c, enterpriseId)
	if err != nil {
		return nil, false, err
	}
	if batch.BatchNo == "" || batch.MappingEnabled != microScrmModel.InventorySkuUpcMappingEnabled {
		return []string{}, false, nil
	}
	upcs, err := r.mappings.GetUpcCodesBySkuItem(c, enterpriseId, batch.BatchNo, skuItem)
	if err != nil {
		return nil, true, err
	}
	return upcs, true, nil
}

type skuUpcSource string

const (
	skuUpcSourceDB  skuUpcSource = "db"
	skuUpcSourcePIM skuUpcSource = "pim"
)

type pimUpcLookupFunc func(c app.Ctx, enterpriseId int64, skuItem string) ([]string, error)

func resolveSkuUpcList(c app.Ctx, resolver skuUpcMappingResolver, enterpriseId int64, skuItem string, pimLookup pimUpcLookupFunc) ([]string, skuUpcSource, bool, error) {
	dbUpcs, dbEnabled, err := resolver.Resolve(c, enterpriseId, skuItem)
	if err != nil {
		return nil, "", false, err
	}
	if dbEnabled && len(dbUpcs) > 0 {
		return dbUpcs, skuUpcSourceDB, false, nil
	}
	pimUpcs, err := pimLookup(c, enterpriseId, skuItem)
	if err != nil {
		return nil, skuUpcSourcePIM, dbEnabled, err
	}
	return pimUpcs, skuUpcSourcePIM, dbEnabled, nil
}

func (s *InventoryService) GetUpcListBySkuItem(c app.Ctx, enterpriseId int64, skuItem string) ([]string, error) {
	upcList, source, dbMissFallback, err := resolveSkuUpcList(c, newSkuUpcMappingResolver(), enterpriseId, skuItem, func(c app.Ctx, enterpriseId int64, skuItem string) ([]string, error) {
		return sourceSvc.NewBrandService().GetUpcListBySkuItem(c, enterpriseId, skuItem)
	})
	if err != nil {
		return nil, err
	}
	if dbMissFallback {
		log.Info(c).Msgf("inventory sku upc source=pim reason=db_miss enterprise_id=%d sku_item=%s", enterpriseId, skuItem)
	}
	if source == skuUpcSourcePIM {
		log.Info(c).Msgf("inventory sku upc source=pim enterprise_id=%d sku_item=%s", enterpriseId, skuItem)
	}
	return upcList, nil
}
