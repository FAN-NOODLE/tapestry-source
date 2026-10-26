<?php


namespace App\Http\Controllers\V1\Message;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroAuthority\Message;

class ActivityMessageController extends BaseController
{
    public function lists()
    {
        $list = (new Message())->ActivityMessageList([
            'enterprise_id' => $this->enterpriseId,
            'page' => $this->request->input('page', 1),
            'size' => $this->request->input('size', 8),
        ]);
        return success($list);
    }

    public function add()
    {
        $params = $this->request->only([
            'message_scene_id',
            'send_time',
            'title',
            'page_type',
            'activity_start_time',
            'expected_num',
            'page_query',
            'activity_condition',
            'page_url'
        ]);
        $params['enterprise_id'] = $this->enterpriseId;
        $result = (new Message())->ActivityMessageAdd($params);
        return success($result);
    }

    public function delete()
    {
        $params = $this->request->only([
            'send_id',
        ]);
        $params['enterprise_id'] = $this->enterpriseId;
        $result = (new Message())->ActivityMessageDelete($params);
        return success($result);
    }
}