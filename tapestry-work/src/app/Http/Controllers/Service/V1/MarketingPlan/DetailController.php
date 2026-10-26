<?php
/**
 * @author  zhaopengfei<zhaopengfei@vchangyi.com>
 * @date    2021/1/11 18:18
 */

namespace App\Http\Controllers\Service\V1\MarketingPlan;


use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Services\MarketingPlanService;
use Illuminate\Http\Request;

class DetailController extends BaseController
{
    public function index(Request $request, MarketingPlanService $marketingPlanService)
    {

        $rules = [
            'task_id' => 'required|integer',
            'enterprise_id' => 'required',
        ];
        $message = [
            'task_id.*' => '任务id不能为空或必须是int类型',
            'enterprise_id.*' => '企业ID不能为空',
        ];
        $this->validate($request, $rules, $message);

        $input = $this->request->all();
        $input['type'] = $request->input('type','0');

        $planInfo = $marketingPlanService->detail($input);

        return success($planInfo);
    }

}