<?php
/**
 * @author  zhaopengfei<zhaopengfei@vchangyi.com>
 * @date    2021/1/11 11:33
 */

namespace App\Http\Controllers\V1\MarketingPlan;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Boss\MiniProgram;
use CrsSdk\Work\MarketingPlan;

class CreateController extends BaseController
{

    /**
     * 创建/编辑计划
     * @param MarketingPlan $marketingPlan
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/11 18:12
     */
    public function index(MarketingPlan $marketingPlan)
    {
        $input = $this->request->all();
        $input['enterprise_id'] = $this->enterpriseId;
        $input['operator_id'] = $this->enterpriseAdminId;
        $input['operator_name'] = $this->enterpriseAdmin['username'];

        $result = $marketingPlan->create($input);

        return success($result);
    }
}