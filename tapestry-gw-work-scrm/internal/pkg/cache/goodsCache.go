package cache

// 商品spu缓存
type GoodsDetailSpuCache struct {
	Cache
}

func (c GoodsDetailSpuCache) Key(unique string) string {
	return "spu_resource_" + unique
}

func (c GoodsDetailSpuCache) MinTime() int  {
	return 60
}

func (c GoodsDetailSpuCache) MaxTime() int {
	return 120
}

// 小程序商品足迹列表缓存
type FootGoodsListCache struct {
	Cache
}

func (c FootGoodsListCache) Key(unique string) string {
	return "footprint:list:" + unique
}

func (c FootGoodsListCache) MinTime() int {
	return 60
}

func (c FootGoodsListCache) MaxTime() int  {
	return 120
}