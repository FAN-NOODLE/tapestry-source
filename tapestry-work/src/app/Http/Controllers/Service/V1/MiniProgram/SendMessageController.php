<?php


namespace App\Http\Controllers\Service\V1\MiniProgram;


use App\Http\Controllers\Service\V1\BaseController;
use App\Models\StoreWorker;
use App\Services\Message\MessageService;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class SendMessageController extends BaseController
{
    /**
     * 发送智慧导购应用消息
     *
     * @param Request $request
     *
     * @return array|JsonResponse
     * @author  user <qiuapeng921@vchangyi.com>
     * @date    2022/7/1 15:07
     */
    public function index(Request $request)
    {
        $enterpriseID = $request->post("enterprise_id");
        $memUserIds   = $request->post("mem_user_ids");
        $title        = $request->post("title");
        $content      = $request->post("content", []);
        $page         = $request->post("page");

        $memUidList = StoreWorker::query()
                ->where('enterprise_id', $enterpriseID)
                ->whereIn('mem_userid', $memUserIds)
                ->pluck('mem_uid')
                ->toArray() ?? [];

        $result = (new MessageService())->sendMiniMessage($enterpriseID, $content, $page, $title, $memUidList);
        Log::info('智汇导购应用消息记录', ['enterprise_id' => $enterpriseID, 'mem_uid' => $memUidList]);
        return successJson((array)$result);
    }
}