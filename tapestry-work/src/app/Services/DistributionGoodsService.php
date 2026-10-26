<?php
/**
 * Created by PhpStorm.
 * User: huangxm
 * Date: 2020-08-03
 * Time: 00:08
 */

namespace App\Services;

use App\Common\Constant;
use App\Models\Distribution;
use App\Models\DistributionGoods;
use App\Models\PartnerEnterpriseSettings;
use CrsSdk\MicroGoodsV1\Goods;

class DistributionGoodsService
{
    const DISTRIBUTION_GOODS_FIELDS
        = [
            'id', 'sku_id', 'spu_id', 'promotion_commission_type', 'promotion_commission_amount',
            'service_commission_type', 'service_commission_amount', 'distribution_start', 'distribution_end', 'enterprise_id'
        ];

    /**
     * 列表
     *
     * @param int    $enterpriseId
     * @param int    $distributionId
     * @param array  $searchParams
     * @param string $pageType
     *
     * @return array
     * @throws \Exception
     */
    public function goodsList(int $enterpriseId,array $searchParams)
    {
        $goods       = new Goods();
        $searchSpuId = [];
        if (!empty($searchParams['spu_name'])) {
            $searchSpuGoods = $goods->list([
                'page'          => 1,
                'size'          => Constant::MAX_SIZE,
                'enterprise_id' => $enterpriseId,
                'search'        => [
                    'name_type' => 'spu_name',
                    'keywords'  => $searchParams['spu_name']
                ]
            ])['list'];
            if (empty($searchSpuGoods)) {
                return [
                    'list'  => [],
                    'total' => 1,
                    'size'  => $searchParams['limit'],
                    'page'  => $searchParams['page']
                ];
            }
            $searchSpuId = array_column($searchSpuGoods, 'spu_id');
        }
        $distributionGoods = DistributionGoods::query()
            ->where('enterprise_id', $enterpriseId)
            ->when(isset($searchParams['sku_id']) && !empty($searchParams['sku_id']), function ($query) use ($searchParams) {
                return $query->where('sku_id', $searchParams['sku_id']);
            })
            ->when(!empty($searchSpuId), function ($query) use ($searchSpuId) {
                return $query->whereIn('spu_id', $searchSpuId);
            })
            ->orderBy('id', 'DESC')
            ->paginate($searchParams['limit'], self::DISTRIBUTION_GOODS_FIELDS, 'page', $searchParams['page']);

        if ($distributionGoods->total() > 0) {
            $skuIds = array_column($distributionGoods->items(), 'sku_id');
            // 查询商品服务
            $goods = $goods->skuDetail([
                'enterprise_id' => $enterpriseId,
                'sku_ids'       => $skuIds
            ])['list'];

            $return = [];
            $distributionGoodsAll = DistributionGoods::query()
                ->where('enterprise_id', $enterpriseId)
                ->orderBy('distribution_start', 'ASC')
                ->get(self::DISTRIBUTION_GOODS_FIELDS)
                ->toArray();
            foreach ($distributionGoodsAll as $distributionGood) {
                $return[$distributionGood['sku_id']][] = [
                    'distribution_start' => $distributionGood['distribution_start'],
                    'distribution_end'   => $distributionGood['distribution_end'],
                    'id'                 => $distributionGood['id']
                ];
            }

            foreach ($distributionGoods->items() as &$goodsItem) {
                $goodsItem->spu_name             = $goods[$goodsItem->sku_id]['sku_name'] ?? '该商品已被删除，请手动删除该条记录';
                $goodsItem->sku_img              = $goods[$goodsItem->sku_id]['resource']['sku_picture'][0]['url'] ?? '';
                $goodsItem->attribute_item_names = $goods[$goodsItem->sku_id]['attribute_item_names'] ?? '';
                $goodsItem->sign_price           = $goods[$goodsItem->sku_id]['sign_price'] ?? 0;
                $goodsItem->sell_price           = $goods[$goodsItem->sku_id]['on_line_price'] ?? 0;
                $distributionTime                = $return[$goodsItem->sku_id];
                $distributionTime                = array_column($distributionTime, null, 'id');
                unset($distributionTime[$goodsItem->id]);
                $goodsItem->date_range = array_values($distributionTime);
                $goodsItem->is_deleted = 0;
            }
        }

        return [
            'list'  => $distributionGoods->items(),
            'total' => $distributionGoods->total(),
            'page'  => $distributionGoods->currentPage(),
            'size'  => $distributionGoods->perPage()
        ];
    }

    /**
     * 小程序端列表
     *
     *
     * @param $params
     * @return array
     * @throws \Exception
     */
    public function appGoodsList($params)
    {
        $newGoodsSpu = [];
        if (!empty($params['new_goods'])) {
            $newGoodsSpu = DistributionGoods::query()
                ->where('enterprise_id', $params['enterprise_id'])
                ->where('distribution_start', '<', time())
                ->where('distribution_end', '>', time())
                ->where('created_at', '>', time() - 86400)
                ->pluck('spu_id')
                ->toArray();
            if (empty($newGoodsSpu)) {
                return [];
            }
        }
        $distributionGoods = DistributionGoods::query()
            ->where('enterprise_id', $params['enterprise_id'])
            ->where('distribution_start', '<', time())
            ->where('distribution_end', '>', time())
            ->when(!empty($newGoodsSpu), function($query) use ($newGoodsSpu) {
                $query->whereIn('spu_id', $newGoodsSpu);
            })
            ->get()
            ->toArray();

        $return = [];
        foreach ($distributionGoods as $item) {
            if (!isset($return[$item['spu_id']])) {
                $return[$item['spu_id']] = [
                    'sku_list' => [$item['sku_id']],
                    'commission' => $item['promotion_commission_amount'],
                    'new_goods' => $item['created_at'] > (time() - 86400) ? Constant::NEW_GOODS_YES : Constant::NEW_GOODS_NO
                ];
            } else {
                $return[$item['spu_id']]['sku_list'][] = $item['sku_id'];
                if ($item['promotion_commission_amount'] > $return[$item['spu_id']]['commission']) {
                    $return[$item['spu_id']]['commission'] = $item['promotion_commission_amount'];
                }
                if ($item['created_at'] > (time() - 86400)) {
                    $return[$item['spu_id']]['new_goods'] = Constant::NEW_GOODS_YES;

                }
            }
        }

        return $return;
    }

    /**
     * 获取多条分销商品记录
     * @param array $params
     * @return array
     */
    public function getDistributionGoodsList($params = [])
    {
        return DistributionGoods::query()
            ->whereIn('sku_id', $params['sku_ids'])
            ->where('enterprise_id', $params['enterprise_id'])
            ->where('deleted_at', 0)
            ->get(['id', 'sku_id', 'promotion_commission_type', 'promotion_commission_amount', 'service_commission_type', 'service_commission_amount', 'distribution_start', 'distribution_end'])
            ->toArray();
    }

    public function appGoodsSkuList($params)
    {
        $newGoodsSpu = [];
        if (!empty($params['new_goods'])) {
            $newGoodsSpu = DistributionGoods::query()
                ->where('enterprise_id', $params['enterprise_id'])
                ->where('distribution_start', '<', time())
                ->where('distribution_end', '>', time())
                ->where('created_at', '>', time() - 86400)
                ->pluck('spu_id')
                ->toArray();
            if (empty($newGoodsSpu)) {
                return [];
            }
        }
        $distributionGoods = DistributionGoods::query()
            ->where('enterprise_id', $params['enterprise_id'])
            ->where('distribution_start', '<', time())
            ->where('distribution_end', '>', time())
            ->when(!empty($newGoodsSpu), function($query) use ($newGoodsSpu) {
                $query->whereIn('spu_id', $newGoodsSpu);
            })
            ->get()
            ->toArray();

        $return = [];
        foreach ($distributionGoods as $item) {
            $new_goods = $item['created_at'] > (time() - 86400) ? Constant::NEW_GOODS_YES : Constant::NEW_GOODS_NO;
            //查询品牌上新
            if($params['new_goods'] == DistributionGoods::NEW_GOODS && $new_goods == DistributionGoods::NEW_GOODS){
                $return[$item['spu_id']]['sku_list'][$item['sku_id']] = [
                    'sku_id' => $item['sku_id'],
                    'commission' => $item['promotion_commission_amount'],
                    'new_goods' => $new_goods
                ];
            }
            if($params['new_goods'] != DistributionGoods::NEW_GOODS){
                $return[$item['spu_id']]['sku_list'][$item['sku_id']] = [
                    'sku_id' => $item['sku_id'],
                    'commission' => $item['promotion_commission_amount'],
                    'new_goods' => $new_goods
                ];
            }

        }
        $data['total'] = count($return);
        $data['list'] = $return;
        //二级分销设置
        $settingDetail = $detail = PartnerEnterpriseSettings::query()
            ->where('enterprise_id', $params['enterprise_id'])
            ->first();
        $settingDetail = !empty($settingDetail) ? $settingDetail->toArray() : [];
        $data['setting'] = $settingDetail;
        return $data;
    }
}
