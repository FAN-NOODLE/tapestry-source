<?php


namespace App\Http\Controllers\V1\Message;

use App\Http\Controllers\V1\BaseController;
use App\Modules\Message\MessageFactory;
use App\Traits\RedisMutexLock;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Validator;

class MiniSubscribeController extends BaseController
{
    use RedisMutexLock;

    public function send()
    {
        Validator::make($this->request->input(), [
            'enterprise_id' => 'required',
            'scene_type'    => 'required',
            'msg'           => 'required|array',
            'user_ids'      => 'array',
            'object_type'   => 'required|integer',
            'object_id'     => 'required|integer',
            'page_query'    => 'string',
            'page_url'      => 'string',
        ], [
            'enterprise_id.required' => '企业id必填',
            'user_ids.array'         => '参数user_ids类型是数组',
            'scene_type.required'    => '参数scene_type必填',
            'msg.*'                  => '参数msg必填，且类型是数组',
            'object_type.required'   => '参数object_type必填',
            'object_id.required'     => '参数object_id必填',
        ])->validate();


        // 并发锁
        $lockParams = $this->request->only(['object_type', 'scene_type', 'user_id', 'object_id']);
        $key        = redisKey('SEND_SUBSCRIBE_MESSAGE', $lockParams);
        $lock       = $this->lock($key, 5000, 1000);
        if (!$lock) {
            Log::info('小程序订阅消息-重复：', $lockParams);
            return success();
        }

        try {
            $input   = $this->request->only(['enterprise_id', 'user_ids', 'scene_type', 'msg', 'object_type', 'object_id', 'page_query', 'page_url']);
            $message = MessageFactory::getInstance(
                MessageFactory::MINI_PROGRAM,
                $input['enterprise_id'],
                $input['user_ids'],
                $input['scene_type'],
                $input['msg'],
                $input);

            Log::info('小程序订阅消息', $input);
            $data = $message->send();
        } catch (\Exception $e) {
            Log::info('小程序订阅消息-失败：' . $e->getMessage());
        }

        return success($data ?? []);
    }
}
