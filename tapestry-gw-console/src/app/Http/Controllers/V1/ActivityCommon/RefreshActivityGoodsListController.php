<?php


namespace App\Http\Controllers\V1\ActivityCommon;


use App\Common\Constant;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroActivity\Activity;

class RefreshActivityGoodsListController extends BaseController
{
    /**
     * true - 显示关闭得活动商品
     */
    const SHOW_DISABLE = true;
    const HIDE_DISABLE = false;

    public function index()
    {
        $input = $this->request->only([
            'id',
            'spu_sku_list',
        ]);

        // 请求活动详情
        $activityGoodsListParams = [
            'activity_rule_id' => $input['id'] ?? 0,
            'enterprise_id'    => $this->enterpriseId,
            'show_disable'     => true,
            'method'           => 'activity-goods-list',
            'type'             => 'all',
            'page'             => Constant::PAGE,
            'size'             => Constant::MAX_DEFAULT_SIZE
        ];

        // 获取活动商品列表
        $activityGoodsList = (new Activity())->getActivity($activityGoodsListParams)['list'] ?? [];
        if (empty($activityGoodsList)) {
            return success([
                'activity_sku' => []
            ]);
        }

        // 如果没有新加得活动商品数据，则只返回活动本身得
        $spuSkuList = $input['spu_sku_list'] ?? [];
        if (!empty($spuSkuList)) {
            $spuIds = array_column($spuSkuList, 'spu_id');
            if (empty($spuIds)) {
                goto END;
            }
            /**
             * {"enterprise_id":"102885540468992",
             * "spu_ids":[117038750032576,117041085822784,117040245681024,
             * 111553854022336,104160652896512],"sale_type":[1,3],"min_price":1}
             */
            $params = [
                'enterprise_id' => $this->enterpriseId,
                'sale_type'     => [1, 3],
                'spu_ids'       => $spuIds
            ];

            $skuList = (new \CrsSdk\MicroGoodsV1\Activity())->getSkuBySpuIds($params);
            if (empty($skuList)) {
                goto END;
            }
            $skuListBySkuId = array_column($skuList, null, 'id');

            foreach ($spuSkuList as &$spuList) {
                foreach ($spuList['sku_list'] as $k => $sku) {
                    $spuList['sku_list'][$k]['stock'] = isset($skuListBySkuId[$sku['sku_id']]) ? $skuListBySkuId[$sku['sku_id']]['stock'] : $sku['stock'];
                }
            }
            unset($spuList);
        }

        END:

        $spuData = [];
        foreach (collect($activityGoodsList)->groupBy('spu_id')->toArray() as $index => $item) {
            foreach ($item as &$i) {
                $i['isOld'] = true;
            }
            $spuData[]     = [
                'spu_id'    => $index,
                'spu_name'  => current($item)['spu_name'] ?? '',
                'url'       => current($item)['picture_url'],
                'canDelSpu' => false,
                'sku_list'  => $item
            ];
        }

        if (!empty($spuSkuList)) {
            $resList = array_merge($spuSkuList, $spuData);
        }

        return success([
            'activity_sku' => $resList ?? $spuData
        ]);
    }

}