<?php
/**
 * @author  zhaopengfei<zhaopengfei@vchangyi.com>
 * @date    2021/1/11 11:29
 */

namespace App\Http\Controllers\V1\MarketingPlan;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\MarketingPlan;

class ListController extends BaseController
{

    /**
     * 营销计划列表
     * @param MarketingPlan $marketingPlan
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/11 18:16
     */
    public function index(MarketingPlan $marketingPlan)
    {
        $input = $this->request->all();
        $input['enterprise_id'] = $this->enterpriseId;

        $list = $marketingPlan->planList($input);

        return success($list);
    }
}