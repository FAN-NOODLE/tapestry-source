<?php


namespace App\Http\Controllers\V1\Wx;


use App\Http\Controllers\V1\BaseController;
use App\WxServices\DynamicMessagesWxService;
use Illuminate\Http\Request;

class DynamicMsgUpdateController extends BaseController
{
    public function index(Request $request)
    {
        $params       = $this->getParams($request);
        $enterpriseId = $params['enterprise_id'];
        $activityId   = $params['activity_id'];
        $state        = $params['state'];

        $msgParams = $request->input('params', 0);

        $result = (new DynamicMessagesWxService())->updateActivityMessage($enterpriseId, $activityId, $state, $msgParams);

        return success($result ?? []);
    }

    protected function getParams(Request $request)
    {
        // 接受参数
        $input = $request->all();

        // 参数验证规则
        $rules = [
            'enterprise_id'       => 'required',
            'activity_id'         => 'required',
            'state'               => 'required',
            'params'              => 'array',
        ];

        // 提示语
        $message = [
            'enterprise_id.required'       => '企业Id不能为空',
            'activity_id.required'         => '活动消息Id不能为空',
            'state.required'               => '活动消息状态不能为空',
            'params.array'                 => '参数为数组格式',
        ];

        $this->validate($request, $rules, $message);

        return $input;
    }
}
