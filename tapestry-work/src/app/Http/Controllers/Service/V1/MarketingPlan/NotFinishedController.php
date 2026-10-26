<?php
/**
 * @author  zhaopengfei<zhaopengfei@vchangyi.com>
 * @date    2021/1/18 4:52
 */

namespace App\Http\Controllers\Service\V1\MarketingPlan;


use App\Http\Controllers\Service\V1\BaseController;
use App\Services\MarketingPlanService;
use Illuminate\Http\Request;

class NotFinishedController extends BaseController
{

    public function index(Request $request, MarketingPlanService $marketingPlanService)
    {
        $rules = [
            //'plan_id' => 'required|integer',
            'enterprise_id' => 'required',
            'task_id' => 'required|integer'
        ];

        $message = [
            //'plan_id.*' => '计划ID不合法',
            'enterprise_id.*' => '企业ID不合法',
            'task_id.*' => '任务ID不合法'
        ];
        $this->validate($request, $rules, $message);
        $input = $request->all();

        $list =  $marketingPlanService->notFinished($input);

        return success($list);
    }
}