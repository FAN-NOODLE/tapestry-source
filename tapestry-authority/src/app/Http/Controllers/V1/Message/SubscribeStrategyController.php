<?php


namespace App\Http\Controllers\V1\Message;

use App\Http\Controllers\V1\BaseController;
use App\Models\EnterpriseMessageScene;
use App\Models\EnterpriseActivityMessage;

class SubscribeStrategyController extends BaseController
{
    /**
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index()
    {
        $enterpriseId = $this->request->input('enterprise_id', 0);
        $type         = $this->request->input('type', 0);
        $messageSceneType         = $this->request->input('message_scene_type', 1);

        $enterpriseMessageScenes = EnterpriseMessageScene::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('type', $messageSceneType)
            ->when($type, function ($query) use ($type) {
                return $query->where('message_scene_id', $type);
            })
            ->with(['belongstoMessageScene','activityMessage'])
            ->get();

        foreach ($enterpriseMessageScenes as &$enterpriseMessageScene) {
            $enterpriseMessageScene['enterprise_available'] = $enterpriseMessageScene['status'];
            $enterpriseMessageScene['id']                   = $enterpriseMessageScene->belongstoMessageScene['id'];
            $enterpriseMessageScene['desc']                 = $enterpriseMessageScene->belongstoMessageScene['desc'];
            $enterpriseMessageScene['sort']                 = $enterpriseMessageScene->belongstoMessageScene['sort'];

            if ($this->request->input('gw', 0)) {
                $enterpriseMessageScene['touch_scene']       = $enterpriseMessageScene->belongstoMessageScene['touch_scene'];
                $enterpriseMessageScene['message_scene']     = $enterpriseMessageScene->belongstoMessageScene['message_scene'];
                $enterpriseMessageScene['default_send_time'] = $enterpriseMessageScene->belongstoMessageScene['default_send_time'];
                $enterpriseMessageScene['touch_page']        = json_decode($enterpriseMessageScene->belongstoMessageScene['touch_page'], 1);
                $enterpriseMessageScene['categories']        = json_decode($enterpriseMessageScene->belongstoMessageScene['available_categories'], 1);
            }

            // 手动策略 获取预计发送数量
            if ($messageSceneType == 2) {
                $enterpriseMessageScene['used_num'] = $enterpriseMessageScene->activityMessage
                    ->where('send_status', EnterpriseActivityMessage::SEND_STATUS_NO)
                    ->where('enable_status', EnterpriseActivityMessage::ENABLE_STATUS_ON)
                    ->where('enterprise_id', $enterpriseId)
                    ->sum('expected_num');
            }

            unset($enterpriseMessageScene['status'], $enterpriseMessageScene['message_scene_id']);
            unset($enterpriseMessageScene['created_at'], $enterpriseMessageScene['updated_at']);
            unset($enterpriseMessageScene->belongstoMessageScene);
            unset($enterpriseMessageScene->activityMessage);
        }

        return success(array_values($enterpriseMessageScenes->sortBy('sort')->toArray()));
    }
}
