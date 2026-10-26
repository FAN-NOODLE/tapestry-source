<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/5/19 0019
 * Time: 20:28
 */

namespace App\Http\Controllers\V1\ActivityCommon\DiscountCard;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Activity;

class SkuListController extends BaseController
{
    /**
     * sku列表
     *
     * @return \Illuminate\Http\JsonResponse
     * @author  baizhongyang <baizhongyang@vchangyi.com>
     * @date    2020/6/15 11:20
     */
    public function index()
    {
        $params                  = $this->request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['source_id']     = $this->enterpriseId;
        // 返回上架中的商品
        $params['status'] = [2];
        // 过滤仅自提商品
        $params['sale_type'] = [1, 3];
        // sku列表数据
        $skuList = (new Activity())->getSkuList($params);

        if (empty($skuList['list'])) {
            return success($skuList);
        }
        // 获取sku占用列表数据
        $params['type']     = $params['occupy_type'] ?? '';
        $params['start_at'] = $params['start_at'] ?? 0;
        $params['end_at']   = $params['end_at'] ?? 0;
        $params['rule_id']  = $params['rule_id'] ?? 0;
        $occupySpuList      = (new \CrsSdk\MicroActivity\Activity())->activitySpu($params);

        foreach ($skuList['list'] as &$value) {
            if (in_array($value['spu_id'], $occupySpuList['spu_ids'])) {
                $value['occupy'] = 1; // 被占用
            } else {
                $value['occupy'] = 0; //没有被占用
            }
        }

        return success($skuList);
    }

}
