<?php

namespace App\Http\Controllers\Html\V1\MarketingPromotion;

use App\Common\Constant;
use App\Http\Controllers\Html\V1\BaseController;
use CrsSdk\MicroActivity\Activity;

class ActivityListController extends BaseController
{
    /**
     * 获取活动列表
     *
     * @return array
     * @throws \Exception
     */
    public function index()
    {
        $params = request()->only(['keywords', 'activity_type', 'page', 'size']);

        $search = [
            'enterprise_id' => $this->enterpriseId,
            'activity_type' => [
                Constant::MARKETING_GOODS,
                Constant::MARKETING_COUPON,
                Constant::MARKETING_HELPER_POST,
                Constant::MARKETING_HELPER_BUY_GIVE,
                Constant::MARKETING_SECKILL,
                Constant::MARKETING_LOTTERY,
                Constant::MARKETING_GROUP,
                Constant::MARKETING_BARGAIN,
            ]
        ];
        if (isset($params['keywords']) && $params['keywords'] != '') {
            $search['search_value'] = $params['keywords'];
            $skuParams = [
                'enterprise_id' => $this->enterpriseId,
                'source_id' => $this->enterpriseId,
                'status' => [Constant::STATUS_UP],
                'goods_type' => Constant::ACTIVITY_GOODS,
                'search_value' => $params['keywords']
            ];
            $list = (new \CrsSdk\MicroGoodsV1\Activity())->getSkuList($skuParams);

            if (empty($list['list'])) {
                return success([
                    'data' => $list['list'] ?? [],
                    'total' => $list['total'] ?? 0,
                    'page' => $params['page'] ?? 1,
                    'size' => $params['size'] ?? 10
                ]);
            }

            $search['sku_ids'] = array_column($list['list'], 'id');
        }

        if (!empty($params['activity_type'])) {
            $search['activity_type'] = $params['activity_type'];
        }
        $search['page'] = $params['page'] ?? 1;
        $search['size'] = $params['size'] ?? 10;

        switchAppHostByDomain($this->domain);

        $list = (new Activity())->activitySkuList($search);

        return success([
            'data' => $list['list'] ?? [],
            'total' => $list['total'] ?? 0,
            'page' => $search['page'],
            'size' => $search['size']
        ]);
    }
}
