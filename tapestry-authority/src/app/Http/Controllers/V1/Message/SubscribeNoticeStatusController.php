<?php


namespace App\Http\Controllers\V1\Message;

use App\Common\RedisKey;
use App\Http\Controllers\V1\BaseController;
use App\Models\EnterpriseMessageScene;
use App\Models\SubscribeMessageRecord;
use Illuminate\Support\Facades\Validator;

class SubscribeNoticeStatusController extends BaseController
{
    /**
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index()
    {
        Validator::make($this->request->input(), [
            'enterprise_id' => 'required',
            'user_id'       => 'required',
        ], [
            'enterprise_id.required' => '企业id必填',
            'user_id.required'       => '用户id必填',

        ])->validate();

        $enterpriseId = $this->request->input('enterprise_id');
        $userId       = $this->request->input('user_id');

        $status = 1; // 默认提示
        $data   = [];

        $enterpriseMessageScenes = $this->getEnterpriseMessageScenes($enterpriseId);

        $types      = array_column($enterpriseMessageScenes, 'message_scene_id');
        $sceneTypes = array_column($enterpriseMessageScenes, 'scene_type');

        $unsendRecordsCount = $this->getUnsendRecordsCount($enterpriseId, $userId, $sceneTypes);

        $unsendRecordsCount = array_column($unsendRecordsCount, null, 'scene_type');

        $enterpriseMessageScenes = collect($enterpriseMessageScenes);
        foreach ($types as $type) {
            $data[$type] = $status;

            $enterpriseMessageScene = $enterpriseMessageScenes->firstWhere('message_scene_id', $type);

            if (is_null($enterpriseMessageScene)) {
                $data[$type] = 4;// 4企业消息场景不存在
                continue;
            }

            // 2企业停用 3企业不可用
            if (1 <> $enterpriseMessageScene['status']) {
                $data[$type] = $enterpriseMessageScene['status'];// 2企业停用，3行业不可用
                continue;
            }

            $sceneType = $enterpriseMessageScene['scene_type'] ?? 0;

            if (!empty($unsendRecordsCount[$sceneType])) {
                $data[$type] = 5;  // 用户订阅了还没发送
            }
        }

        return success($data);
    }


    protected function getEnterpriseMessageScenes($enterpriseId)
    {
        $enterpriseMessageScenes = $this->redis->get(RedisKey::ENTERPRISE_MESSAGE_SCENE . $enterpriseId);

        if ($enterpriseMessageScenes) {
            $enterpriseMessageScenes = json_decode($enterpriseMessageScenes, 1);
        } else {
            $enterpriseMessageScenes = EnterpriseMessageScene::query()
                ->where('enterprise_id', $enterpriseId)
                ->get()->toArray();
            $this->redis->set(RedisKey::ENTERPRISE_MESSAGE_SCENE . $enterpriseId, json_encode($enterpriseMessageScenes, 320));
        }

        return $enterpriseMessageScenes;
    }

    protected function getUnsendRecordsCount($enterpriseId, $userId, $sceneTypes)
    {
        $unsendRecordsCount = $this->redis->get(RedisKey::USER_SUBSCRIBE_NOTICE_STATUS . $userId);

        if ($unsendRecordsCount) {
            $unsendRecordsCount = json_decode($unsendRecordsCount, 1);
        } else {
            $unsendRecordsCount = SubscribeMessageRecord::query()
                ->selectRaw('scene_type, count(1) as subscribe_count')
                ->whereIn('scene_type', $sceneTypes)
                ->where('enterprise_id', $enterpriseId)
                ->where('user_id', $userId)
                ->where('send_num', 0)
                ->where('fail_reason', '')
                ->has('template')
                ->groupBy('scene_type')
                ->get()
                ->toArray();
            $this->redis->set(RedisKey::USER_SUBSCRIBE_NOTICE_STATUS . $userId, json_encode($unsendRecordsCount, 320));
        }
        return $unsendRecordsCount;
    }
}
