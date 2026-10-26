<?php


namespace App\Http\Controllers\Service\V1\EventNotify;


use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Ddc\DdcMessage;
use App\Services\Events\ShareLinkChangeEventService;
use Illuminate\Http\Request;

class ShareLinkChangeEventNotifyController extends BaseController
{

    private $eventService;

    public function __construct()
    {
        $this->eventService = new ShareLinkChangeEventService();
    }

    public function add(Request $request)
    {
        $this->eventService->pushDdcMessage(DdcMessage::GUIDE_USER_BINDING_ADD, $this->getRequestDataList($request));
        return success([]);
    }

    public function update(Request $request)
    {
        $this->eventService->pushDdcMessage(DdcMessage::GUIDE_USER_BINDING_UPDATE, $this->getRequestDataList($request));
        return success([]);
    }

    public function delete(Request $request)
    {
        $this->eventService->pushDdcMessage(DdcMessage::GUIDE_USER_BINDING_DELETE, $this->getRequestDataList($request));
        return success([]);
    }

    private function getRequestDataList(Request $request)
    {
        $this->validate($request, [
            'list' => 'required|array',
//            'list.*.enterprise_id' => 'required|integer',
//            'list.*.third_user_id' => 'required|integer',
//            'list.*.curr_worker_id' => 'required|integer',
//            'list.*.old_worker_id' => 'integer',
//            'list.*.bind_time' => 'integer',
        ]);

        return $request->input('list');
    }
}