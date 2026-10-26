<?php
/**
 * 获取任务详情
 */
namespace App\Http\Controllers\App\V1\MarketingPlan;

use App\Http\Controllers\App\V1\BaseController as Controller;
use CrsSdk\Work\MarketingPlan;
use Illuminate\Http\Request;

class TaskDetailController extends Controller
{
    /**
     * 获取任务详情
     * User: Hui Yang
     * Email: 2361313833@qq.com
     * DateTime: 2021/1/13 17:44
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $rules = [
            'task_id' => 'required|min:1'
        ];
        $messages = [
            'task_id.*' => '任务id参数错误'
        ];
        $this->validate($request, $rules, $messages);

        $input = $request->input();
        $input['task_id'] = $request->input('task_id');
        $input['domain'] = $this->loginAdminEnterprise['domain'];
        $storeId = $this->request->post('store_id', 0);
        $input['store_id'] = $storeId;
        $input['enterprise_id'] = $this->enterpriseId;
        $input['mobile'] = $this->mobile;

        $result = (new MarketingPlan())->taskDetail($input);

        return success($result['data'] ?? $result);
    }
}
