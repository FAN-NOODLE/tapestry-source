<?php


namespace App\Http\Controllers\V1\Message;

use App\Http\Controllers\V1\BaseController;
use App\Models\SubscribeMessageRecord;

class MiniSubscribeRecordStatusController extends BaseController
{
    /**
     * @return \Illuminate\Http\JsonResponse
     */
    public function index()
    {
        $userId     = $this->request->input('user_id', 0);
        $objectId   = $this->request->input('object_id', 0);
        $objectIds  = $this->request->input('object_ids', []);
        $objectType = $this->request->input('object_type', 0);

        if ($objectId) {
            $objectIds = [$objectId];
        }

        $unsendRecords = SubscribeMessageRecord::query()
            ->selectRaw('object_id,count(1) as subscribe_status')
            ->where('user_id', $userId)
            ->where('object_type', $objectType)
            ->whereIn('object_id', $objectIds)
            ->where('send_num', 0)
            ->where('fail_reason','')
            ->has('template')
            ->groupBy('object_id')
            ->get()->toArray();

        $unsendRecords = array_column($unsendRecords, null, 'object_id');

        $res = [];
        foreach ($objectIds as $inputObjectId) {
            $subscribeStatus = 0;
            if (isset($unsendRecords[$inputObjectId])) {
                $subscribeStatus = $unsendRecords[$inputObjectId]['subscribe_status'] ? 1 : 0;
            }
            $res[] = [
                'object_id'        => $inputObjectId,
                'subscribe_status' => $subscribeStatus,
            ];
        }

        $res = $objectId ? current($res) : $res;

        return success($res);
    }
}
