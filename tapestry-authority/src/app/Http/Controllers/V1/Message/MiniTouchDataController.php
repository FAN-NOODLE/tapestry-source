<?php
/** 获取订阅消息模板
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/11/9
 * Time: 10:00 AM
 */

namespace App\Http\Controllers\V1\Message;

use App\Http\Controllers\V1\BaseController;
use App\Models\EnterpriseMessageScene;
use App\Models\SubscribeMessageTouch;
use App\Models\SubscribeMessageTouchTemplate;
use Illuminate\Support\Facades\Validator;

class MiniTouchDataController extends BaseController
{

    public function index()
    {
        Validator::make($this->request->all(), [
            'enterprise_id' => 'required',
        ], [
            'enterprise_id.required' => '企业id必填',
        ])->validate();

        $enterpriseId = $this->request->input('enterprise_id');

        $expectSceneTypes = EnterpriseMessageScene::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('status', '<>', 1)
            ->pluck('scene_type')
            ->all();

        // 获取模板触摸数据
        $touchData = SubscribeMessageTouchTemplate::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('pri_tmpl_id', '<>', 1)
            ->where('pri_tmpl_id', '<>', '')
            ->where('touch_ids', '<>', '')
            ->whereNotIn('scene_type', $expectSceneTypes)
            ->get(['touch_ids', 'pri_tmpl_id', 'scene_type'])
            ->toArray();

        $list = $data = $ids = [];
        foreach ($touchData as $val) {
            $touchIds = json_decode($val['touch_ids'], true);
            foreach ($touchIds as $id) {
                $ids[]       = $id;
                $list[$id][] = $val['pri_tmpl_id'];
            }
        }

        // 如果有模板触发事件
        if (!empty($ids)) {
            $ids = array_unique($ids);
            // 查询触发事件名称
            $subscribeMessageTouch = SubscribeMessageTouch::query()
                ->whereIn('id', $ids)
                ->pluck('name', 'id')
                ->toArray();
            foreach ($list as $key => $value) {
                $data[$subscribeMessageTouch[$key]] = $value;
            }
        }

        return success($data);
    }
}
