<?php


namespace App\Services\Events;


use App\Exceptions\Error;
use App\Library\Ddc\DdcMessage;
use App\Library\Uc\Member;
use App\Services\StoreWorkerService;
use App\Services\UserRemoteService;
use CrsSdk\Boss\Enterprise;
use CrsSdk\MicroUser\User;
use Exception;
use Illuminate\Support\Facades\Log;

class ShareLinkChangeEventService
{

    /**
     * 批量数据每次最多操作数据量
     */
    const BATCH_SIZE_MAX = 20;

    public function pushDdcMessage($method, array $notifyDataList)
    {
        Log::info("ShareLinkChangeEventService.pushDdcMessage", [
            'method' => $method,
            'notifyDataList' => $notifyDataList
        ]);
        $size = count($notifyDataList);
        if ($size == 0) {
            throw new Error(110009, "分享链路通知数据为空");
        }
        // 按照enterprise_id分组
        $epNotifyDataMap = [];
        foreach ($notifyDataList as $notifyData) {
            if (!isset($epNotifyDataMap[$notifyData['enterprise_id']])) {
                $epNotifyDataMap[$notifyData['enterprise_id']] = [];
            }
            array_push($epNotifyDataMap[$notifyData['enterprise_id']], $notifyData);
        }

        foreach ($epNotifyDataMap as $enterpriseId => $epNotifyDataList) {
            $enterprise = (new Enterprise())->detail(['enterprise_id' => $enterpriseId])['data'];

            $totalPage = (int) ceil(count($epNotifyDataList) / self::BATCH_SIZE_MAX);
            for ($page = 1; $page <= $totalPage; $page++) {
                $offset = ($page - 1) * self::BATCH_SIZE_MAX;
                $epNotifyDataPageList = array_slice($epNotifyDataList, $offset, self::BATCH_SIZE_MAX);

                $ddcDataList = $this->batchBuildDdcData($enterprise, $epNotifyDataPageList);
                Log::info("share link data notify by enterprise pagination", [
                    'page' => $page,
                    'ddcDataList' => $ddcDataList
                ]);
                if (!empty($ddcDataList)) {
                    (new DdcMessage())->push($enterprise['hash'], $method, ['list' => $ddcDataList]);
                }

            }
        }

        Log::info("ShareLinkChangeEventService.pushDdcMessage success");
    }

    protected function batchBuildDdcData($enterprise, array $notifyDataList): array
    {
        $enterpriseId = $enterprise['id'];
        $qyDomain = $enterprise['domain'];

        $userIdList = array_unique(array_filter(array_column($notifyDataList, 'third_user_id')));
        $userUnionIdMap = UserRemoteService::listUserUnionId($enterpriseId, $userIdList);

        $workerIdList = array_unique(array_filter(array_merge(array_column($notifyDataList, 'curr_worker_id'),
            array_column($notifyDataList, 'old_worker_id'))));
        $workerMemUidMap = $this->listWorkerMemUid($enterpriseId, $workerIdList);

        $qwMemUidIdMap = [];
        if (!empty($workerMemUidMap)) {
            $qwMemUidIdMap = $this->listQyMemUserId($qyDomain, array_unique(array_values($workerMemUidMap)));
        }


        $ddcDataList = [];
        foreach ($notifyDataList as $notifyData) {
            $currWorkerQwMemUidId = '';
            if (!empty($notifyData['curr_worker_id']) && !empty($workerMemUidMap[$notifyData['curr_worker_id']])) {
                $currWorkerMemUid = $workerMemUidMap[$notifyData['curr_worker_id']];
                if (!empty($currWorkerMemUid) && !empty($qwMemUidIdMap[$currWorkerMemUid])) {
                    $currWorkerQwMemUidId = $qwMemUidIdMap[$currWorkerMemUid];
                }
            }
            if (empty($currWorkerQwMemUidId)) {
                continue;
            }

            $oldWorkerQwMemUidId = '';
            if (!empty($notifyData['old_worker_id']) && !empty($workerMemUidMap[$notifyData['old_worker_id']])) {
                $oldWorkerMemUid = $workerMemUidMap[$notifyData['old_worker_id']];
                if (!empty($oldWorkerMemUid) && !empty($qwMemUidIdMap[$oldWorkerMemUid])) {
                    $oldWorkerQwMemUidId = $qwMemUidIdMap[$oldWorkerMemUid];
                }
            }

            $ddcData = [
                'enterprise_id' => $notifyData['enterprise_id'],
                'third_user_id' => $notifyData['third_user_id'],
                'cur_user_id' => $currWorkerQwMemUidId,
                'ori_user_id' => $oldWorkerQwMemUidId,
                'unionid' => $userUnionIdMap[$notifyData['third_user_id']] ?? '',
                'bind_time' => $notifyData['bind_time'] ?? 0
            ];
            $ddcDataList[] = $ddcData;
        }
        return $ddcDataList;
    }


    protected function listWorkerMemUid(int $enterpriseId, array $workerIdList)
    {
        if (!empty($enterpriseId) && !empty($workerIdList)) {
            $params = [
                'enterprise_id' => $enterpriseId,
                'ids' => $workerIdList
            ];
            $is_page = false;
            try {
                $workerList = StoreWorkerService::list($params, $is_page)['list'];
                $result = [];
                foreach ($workerList as $worker) {
                    $result[$worker['id']] = $worker['mem_uid'] ?? '';
                }
                return $result;
            } catch (Exception $e) {
                Log::error('导购信息列表查询失败', [
                    'enterprise_id' => $enterpriseId,
                    'workerIdList' => $workerIdList,
                    'code' => $e->getCode(),
                    'message' => $e->getMessage(),
                    'exception' => $e->getTraceAsString()
                ]);
            }
        }
        return [];
    }

    protected function listQyMemUserId(string $qyDomain, array $workerMemUidList)
    {
        if (!empty($qyDomain) && !empty($workerMemUidList)) {
            try {
                $ucParams = [
                    'isCallBack' => 0,
                    'memUids' => $workerMemUidList
                ];
                $basicList = (new Member())->basicList($qyDomain, $ucParams)['list'];
                $result = [];
                foreach ($basicList as $basicInfo) {
                    $result[$basicInfo['memUid']] = $basicInfo['memUserid'];
                }
                return $result;
            } catch (Exception $e) {
                Log::error('企业微信用户id查询失败', [
                    'qyDomain' => $qyDomain,
                    'workerMemUidList' => $workerMemUidList,
                    'code' => $e->getCode(),
                    'message' => $e->getMessage(),
                    'exception' => $e->getTraceAsString()
                ]);
            }
        }
        return [];
    }

}