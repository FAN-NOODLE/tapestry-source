<?php
/**
 * @author  zhaopengfei<zhaopengfei@vchangyi.com>
 * @date    2021/1/11 11:32
 */

namespace App\Http\Controllers\V1\MarketingPlan;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\MarketingPlan;

class DetailController  extends BaseController
{

    /**
     * 营销计划详情
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

        $detail  = $marketingPlan->detail($input);
        if(empty($detail)){
            notice(2200003);
        }

        return success($detail);
    }
}