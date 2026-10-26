<?php
/**
 * @author  zhaopengfei<zhaopengfei@vchangyi.com>
 * @date    2021/1/11 18:19
 */

namespace App\Http\Controllers\Service\V1\MarketingPlan;


use App\Http\Controllers\Service\V1\BaseController;
use App\Services\MarketingPlanService;
use Illuminate\Http\Request;

class DeleteController extends BaseController
{
    public function index(Request $request, MarketingPlanService $marketingPlanService)
    {
        $input = $this->getParams($request);

        $result = $marketingPlanService->delete($input);

        return success($result);
    }

    private function getParams(Request $request): array
    {
        $rules = [
            'task_id'       => 'required',
            'enterprise_id' => 'required',
        ];

        $message = [
            'task_id.*'       => '任务id不能为空',
            'enterprise_id.*' => '企业ID不能为空',
        ];
        $input   = $this->validate($request, $rules, $message);
        return $input;
    }
}