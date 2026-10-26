<?php


namespace App\Http\Controllers\V1\Message;

use App\Http\Controllers\V1\BaseController;
use App\Models\SubscribeMessageRecord;
use App\Services\SubscribeMessageRecordService;
use Illuminate\Support\Facades\Validator;

class MiniSubscribeRecordController extends BaseController
{
    public function index()
    {
        Validator::make($this->request->all(), [
            'pri_tmpl_id' => 'required|array',
            'user_id'     => 'required',
            'object_id'   => 'required',
            'object_type' => 'required',
            'type'        => 'required',
            'touch_id'    => 'required',
        ], [
            'pri_tmpl_id.required' => '参数pri_tmpl_id必填，且类型是数组',
            'user_id.required'     => '参数user_id必填',
            'object_id.required'   => '参数object_id必填',
            'object_type.required' => '参数object_type必填',
            'type.required'        => '参数type必填',
            'touch_id.required'    => '参数touch_id必填',
        ])->validate();

        $priTmplId = $this->request->input('pri_tmpl_id');
        $input     = $this->request->only(['user_id', 'object_id', 'object_type', 'type', 'touch_id']);

        $params = [
            'user_id'     => $input['user_id'],
            'type'        => $input['type'] ?? SubscribeMessageRecord::TYPE_ONCE,
            'object_type' => $input['object_type'],
            'object_id'   => $input['object_id'],
            'touch_id'    => $input['touch_id'] ?? 0,
        ];

        $subscribeMessageStatus = (new SubscribeMessageRecordService())->addSubscribeMessage($params, $priTmplId);

        return success($subscribeMessageStatus);
    }

    /**
     * @return \Illuminate\Http\JsonResponse
     */
    public function userSubscribeList()
    {
        $userId     = $this->request->input('user_id', 0);
        $objectType = $this->request->input('object_type', 0);
        $data       = (new SubscribeMessageRecordService())->getUserSubscribeList($userId, $objectType);
        return success($data);
    }

    public function subscribeList()
    {
        $params = $this->request->only(['enterprise_id', 'scene_type', 'object_type', 'object_id', 'size', 'page']);
        $data   = (new SubscribeMessageRecordService())->getSubscribeList($params);
        return success($data);
    }
}
