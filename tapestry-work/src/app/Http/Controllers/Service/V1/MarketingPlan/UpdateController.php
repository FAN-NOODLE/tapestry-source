<?php
/**
 * @author  zhaopengfei<zhaopengfei@vchangyi.com>
 * @date    2021/1/11 11:49
 */

namespace App\Http\Controllers\Service\V1\MarketingPlan;


use App\Http\Controllers\Service\V1\BaseController;
use App\Models\MarketingPlan;
use App\Models\MarketingPlanVisibleRange;
use App\Services\MarketingPlanService;
use Illuminate\Http\Request;
use Illuminate\Validation\Rule;

class UpdateController extends BaseController
{
    public function index(Request $request, MarketingPlanService $marketingPlanService)
    {

        $rules = [
            //'plan_title' => 'required',
            'enterprise_id' => 'required',
            'operator_id' => 'integer',
            'operator_name' => 'string',
            'object_type' => [
                'required',
                Rule::in([
                    MarketingPlanVisibleRange::CHOICE_AREA,
                    MarketingPlanVisibleRange::CHOICE_DEPARTMENT
                ])
            ],
            'object_array' => 'array',
            //'plan_id' => 'required|integer',
            'task_array' => 'array'//数据效验
        ];

        $message = [
            //'plan_title.*' => '计划名称不能为空',
            'enterprise_id.*' => '企业ID不能为空',
            'object_type.*' => '可见范围表类型必须是1或者2',
            'object_array.*' => '门店或部门必须是数组',
            'task_array.*' => '任务数据必须是数组'
        ];
        $this->validate($request, $rules, $message);
        $input = $request->all();


        $result = $marketingPlanService->updatePlan($input);

        return success($result);
    }

}