<?php
/**
 * @author  zhaopengfei<zhaopengfei@vchangyi.com>
 * @date    2021/1/11 18:19
 */

namespace App\Http\Controllers\Service\V1\MarketingPlan;


use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Services\MarketingPlanService;
use Illuminate\Http\Request;

class ListController extends BaseController
{
    public function index(Request $request, MarketingPlanService $marketingPlanService)
    {
        $rules = [
            'enterprise_id' => 'required',
        ];

        $message = [
            'enterprise_id.*' => '企业ID不能为空',
        ];
        $this->validate($request, $rules, $message);
        $input = $request->all();

        $planInfo = $marketingPlanService->list($input);

        return success($planInfo);
    }
}