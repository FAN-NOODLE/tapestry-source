<?php


namespace App\Http\Controllers\V1\Message;

use App\Common\RedisKey;
use App\Http\Controllers\V1\BaseController;
use App\Models\EnterpriseMessageScene;
use Illuminate\Support\Facades\Validator;

class UpdateSubscribeStrategyController extends BaseController
{
    /**
     * @throws \Exception
     */
    public function index()
    {
        Validator::make($this->request->all(), [
            'enterprise_id'    => 'required',
            'message_scene_id' => 'required',
            'config'           => 'integer',
            'status'           => 'in:1,2',
        ], [
            'enterprise_id.required'    => '企业id必填',
            'message_scene_id.required' => '场景id必填',
            'config.integer'            => '配置项config为正数',
            'status.in'                 => '状态开关：1启用 2停用',
        ])->validate();

        $enterpriseMessageScene = EnterpriseMessageScene::query()
            ->where('enterprise_id', $this->request->input('enterprise_id'))
            ->where('message_scene_id', $this->request->input('message_scene_id'))
            ->first();

        if (is_null($enterpriseMessageScene)) {
            throw new \Exception('企业场景不存在');
        }

        $updateParams = $this->request->only(['config', 'status']);

        if (isset($updateParams['status']) && 3 == $enterpriseMessageScene->status) {
            unset($updateParams['status']);
        }

        // 如果该小程序配置行业类目不可用，则不可修改状态，仅修改配置项
        $updateParams = (3 == $enterpriseMessageScene->status)
            ? $this->request->only(['config'])
            : $this->request->only(['config', 'status']);

        if ($updateParams){
            $enterpriseMessageScene->update($updateParams);
            $this->redis->del([RedisKey::ENTERPRISE_MESSAGE_SCENE_COUNT . $this->request->input('enterprise_id')]);
            $this->redis->del([RedisKey::ENTERPRISE_MESSAGE_SCENE . $this->request->input('enterprise_id')]);
        }

        return success([]);
    }
}
