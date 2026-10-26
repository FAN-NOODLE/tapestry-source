<?php

namespace App\Http\Controllers\Service\V1\MiniProgram;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\StoreWorker;
use App\Services\Message\MessageService;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

/**
 * 对应智慧导购的应用，给BA推送消息
 */
class SendMessagePlanTask extends BaseController
{

    public function index(Request $request)
    {
        // 获取请求数据
        $enterpriseID = $request->post("enterprise_id");
        $userIDList = $request->post("user_id_list");
        $content =
            [
                [
                    'key'   => '说明',
                    'value' => '你的任务计划已开启，点击查看你的任务吧',
                ],
            ];
        $title   = '导购任务计划提醒';

        $page = 'pages/task/task?is_current=1';

        $memUidList = $this->getMemUid($userIDList ?? []);

        try {
            $result = (new MessageService())->sendMiniMessage($enterpriseID, $content, $page, $title, $memUidList);

            Log::info('导购营销任务提醒发送记录', [
                'enterprise_id' => $enterpriseID,
                'mem_uid'       => $memUidList
            ]);
        } catch (\Exception $e) {
            return success([
                "msg" => $e->getMessage(),
            ]);
        }
        return success([
            "is_success" => $result,
        ]);
    }

    private function getMemUid(array $userIDList) : array {
        $result = StoreWorker::whereIn('mem_userid', $userIDList)->pluck('mem_uid')->toArray();
        return $result ?? [];
    }
}