<?php

namespace App\Http\Controllers\V1\Message;

use App\Http\Controllers\V1\BaseController;
use App\Jobs\MessageSendJob;
use App\Models\MessageScene;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Validator;

class SendController extends BaseController
{
    /**
     * 发送消息
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \GuzzleHttp\Exception\GuzzleException
     */
    public function index()
    {
        $input = $this->request->all();
        Validator::make($input, [
            'enterprise_id' => 'required',
            'type'          => 'required',
            'msg'           => 'required|array',
            'user_ids'      => 'array',
            'object_type'   => 'integer',
            'object_id'     => 'integer',
            'page_query'    => 'string',
            'page_url'      => 'string',
        ], [
            'enterprise_id.required' => '企业id必填',
            'type.required'          => '类型必填',
            'msg.required'           => 'msg必填',
            'msg.array'              => 'msg类型是数组',
            'user_ids.array'         => '用户ids类型是数组',
        ])->validate();

        // 1010 砍价结果提醒  对应小程序 15砍价失败提醒 16砍价成功提醒
        $scene = MessageScene::query()->where('id', $input['type'])->when(1010 == $input['type'], function ($query) use ($input) {
            $sceneTypes = (1 == $input['msg']['haggle_progress']) ? [1010, 16] : [1010, 15];
            $query->with([
                'channels' => function ($query) use ($sceneTypes) {
                    $query->whereIn('scene_type', $sceneTypes);
                }
            ]);
        }, function ($query) {
            $query->with('channels');
        })->first();

        $res = [];
        foreach ($scene['channels'] as $channel) {

            $job                      = new MessageSendJob($channel['channel'], $input['enterprise_id'], $input['user_ids'], $channel['scene_type'], $input['msg'], $input);
            $res[$channel['channel']] = $job->handle();
            // dispatch($job);

            Log::info('消息发送入队成功');
        }

        return success($res);
    }
}
