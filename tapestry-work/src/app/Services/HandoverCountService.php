<?php
/**
 * Created by PhpStorm.
 * User: meijie
 * Date: 2019-03-23
 * Time: 23:31
 */

namespace App\Services;

use App\Models\MemberChangeRecord;
use App\Models\MemberUserRelation;
use Illuminate\Support\Facades\DB;

class HandoverCountService
{
    /**
     * 通过员工UID集合 - 获取待确认客户/待分配客户的总数
     * @param $enterpriseId
     * @param $memUids
     * @return array
     */
    public function getHandoverCount($enterpriseId, $memUids)
    {
        // 待分配
        $waitDistributed = MemberUserRelation::query()
            ->select(DB::raw('mem_uid, count(*) as total'))
            ->where(['enterprise_id' => $enterpriseId])
            ->whereIn('mem_uid', $memUids)
            ->groupBy('mem_uid')
            ->get()
            ->toArray();

        // 待确认
        $waitConfirm = MemberChangeRecord::query()
            ->select(DB::raw('old_mem_uid, count(*) as total'))
            ->where(['enterprise_id' => $enterpriseId])
            ->where('status', MemberChangeRecord::CONFIRM_HANDOVER_FALSE)
            ->whereIn('old_mem_uid', $memUids)
            ->groupBy('old_mem_uid')
            ->get()
            ->toArray();

        $data = $this->formatData($waitDistributed, $waitConfirm);
        $result = $this->formatResult($data, $memUids);

        return $result;
    }

    /**
     * 处理待分配/待确认客户的格式
     *
     * @param $waitDistributed
     * @param $waitConfirm
     * @return array
     */
    public function formatData($waitDistributed, $waitConfirm)
    {
        $data = [];

        if (empty($waitDistributed)) {
            $data['distributed'] = [];
        } else {
            foreach ($waitDistributed as $distributed) {
                $data['distributed'][$distributed['mem_uid']] = $distributed['total'];

            }
        }

        if (empty($waitConfirm)) {
            $data['confirm'] = [];
        } else {
            foreach ($waitConfirm as $confirm) {
                $data['confirm'][$confirm['old_mem_uid']] = $confirm['total'];
            }
        }

        return $data;
    }

    /**
     * 处理数据为返回的格式
     *
     * @param $data
     * @param $memUids
     * @return array
     */
    public function formatResult($data, $memUids)
    {
        $result = [];
        $confirm = $data['confirm'];
        $distributed = $data['distributed'];

        foreach ($memUids as $memUid) {
            $result['mem_uid'][$memUid] = [
                'no_distributed' => $distributed[$memUid] ?? 0,
                'wait_confirm' => $confirm[$memUid] ?? 0,
            ];
        }

        return $result;
    }
}
