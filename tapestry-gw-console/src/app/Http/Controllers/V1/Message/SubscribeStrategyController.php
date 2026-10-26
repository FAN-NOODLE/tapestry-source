<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/3/17
 * Time: 18:38
 */

namespace App\Http\Controllers\V1\Message;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroAuthority\Message;
use Illuminate\Support\Facades\Validator;

class SubscribeStrategyController extends BaseController
{
    /**
     * 渠道列表
     *
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $data = (new Message())->subscribeStrategy([
            'enterprise_id' => $this->enterpriseId,
            'type'          => $this->request->input('type', 0),
            'gw'            => 1,
            'message_scene_type' => $this->request->input('message_scene_type', 1),
        ]);

        return success($data);
    }

    /**
     * 编辑配置或开关
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function update()
    {
        Validator::make($this->request->all(), [
            'message_scene_id' => 'required',
            'config'           => 'integer',
            'status'           => 'in:1,2',
        ], [
            'message_scene_id.required' => '场景id必填',
            'config.integer'            => '配置项config为正数',
            'status.in'                 => '状态开关：1启用 2停用',
        ])->validate();

        $params = $this->request->only(['message_scene_id', 'config', 'status']);

        $params['enterprise_id'] = $this->enterpriseId;

        $data = (new Message())->updateSubscribeStrategy($params);

        return success($data);
    }
}