<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/5/12 0012
 * Time: 14:48
 */

namespace App\Http\Controllers\V1\Goods;

use App\Common\Constant;
use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroActivity\Activity;
use CrsSdk\MicroGoodsV1\Activity as GoodsActivity;
use CrsSdk\MicroGoodsV1\Goods;

class ActivitySpuListController extends BaseController
{
    /**
     * 筛选已发布的商品
     */
    const FILL_RELEASE = 1;

    /**
     * 商品被其他活动占用
     */
    const IS_ACTIVITY_GOODS_YES = 1;

    /**
     * 商品没有被其他活动占用
     */
    const IS_ACTIVITY_GOODS_NO = 0;

    /**
     * 活动spu列表只获取非自提商品
     */
    const ONLY_ONLINE = 1;

    public function index()
    {
        $param                    = $this->request->input();
        $search['enterprise_id']  = $this->enterpriseId;
        $search['search']         = $param['search']['search_value'] ?? null;
        $search['search_type']    = $param['search']['search_type'] ?? null;
        $search['page']           = $param['page'] ?? Constant::PAGE;
        $search['size']           = $param['size'] ?? Constant::DEFAULT_SIZE;
        $search['is_simple_list'] = self::FILL_RELEASE;
        $search['sale_type'] = $param['sale_type'] ?? 0;
        $search['is_back_stock'] = 1;
//        $search['only_online'] = self::ONLY_ONLINE;
        $minPrice = $param['min_price'] ?? 0;

        try {
            $allSpu = (new Goods())->spuList(array_filter($search));
        } catch (\Exception $e) {
            throw new Error($e->getCode(), $e->getMessage());
        }
        $allSpuList = $allSpu['list'] ?? [];
        $diffSpuIds = $this->getMinSpuIds($allSpuList, $minPrice);

        unset($search);
        $search['enterprise_id'] = $this->enterpriseId;
        $search['rule_id']       = $param['rule_id'] ?? 0;
        $search['type']          = $param['activity_type'] ?? 0;
        $search['start_at']      = $param['start_at'] ?? 0;
        $search['end_at']        = $param['end_at'] ?? 0;
        try {
            $activitySpuIds = (new Activity())->activitySpu($search)['spu_ids'];
        } catch (\Exception $e) {
            throw new Error($e->getCode(), $e->getMessage());
        }
        // 组装数据
        $spuList = [];
        foreach ($allSpu['list'] as $key => $spu) {
            $spuList[$key]['id']       = $spu['id'];
            $spuList[$key]['spu_name'] = $spu['spu_name'];
            $spuList[$key]['url']      = $spu['spu_source']['url'] ?? '';
            $spuList[$key]['sell_point'] = $spu['sell_point'] ?? '';
            $spuList[$key]['stock'] = $spu['stock'];
            // 如果spuId在被占用的spu中或者
            if (in_array($spu['id'], $activitySpuIds)
                || (!empty($minPrice) && in_array($spu['id'], $diffSpuIds))
            ) {
                $spuList[$key]['is_activity'] = self::IS_ACTIVITY_GOODS_YES;
            } else {
                $spuList[$key]['is_activity'] = self::IS_ACTIVITY_GOODS_NO;
            }

        }
        $res = [
            'list'  => $spuList,
            'page'  => $allSpu['page'],
            'size'  => $allSpu['size'],
            'total' => $allSpu['total']
        ];
        return success($res);
    }


    /**
     * 获取spuId全部小于最小价格的集合
     * @param array $allSpuList 商品列表
     * @param int   $minPrice 最小价格分
     *
     * @return array
     */
    private function getMinSpuIds($allSpuList = [], $minPrice = 0)
    {
        $diffSpuIds = $spuIds=[];
        // 如果是一口价需要过滤则
        if (!empty($minPrice) && !empty($allSpuList)) {
            $allSpuIds = array_column($allSpuList, 'id');
            try {
                $skuParams = [
                    'enterprise_id' => $this->enterpriseId,
                    'spu_ids'       => $allSpuIds,
                    'sale_type'     => [1, 3]
                ];
                $skuList   = (new GoodsActivity())->getSkuBySpuIds($skuParams);
            } catch (\Exception $e) {
                throw new Error($e->getCode(), $e->getMessage());
            }

            foreach ($skuList ?? [] as $item) {
                if ($item['on_line_price'] >= $minPrice) {
                    $spuIds[$item['spu_id']] = $item['spu_id'];
                }
            }

            sort($spuIds);

            $diffSpuIds = array_diff($allSpuIds, $spuIds);
        }

        return $diffSpuIds;
    }

}
