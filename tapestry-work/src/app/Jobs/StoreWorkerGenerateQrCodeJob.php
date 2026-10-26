<?php

namespace App\Jobs;

use App\Common\Constant;
use App\Models\Scrm\Members;
use App\Services\StoreWorkerService;
use Illuminate\Support\Facades\Log;

class StoreWorkerGenerateQrCodeJob extends Job
{
    const PAGE_SIZE = 1000;
    public $queue = Constant::STORE_WORKER_GENERATE_QR_CODE_JOB_USE_WORK;
    /**
     * 异步生成门店员工的二维码和太阳码
     *
     * @return void
     */
    protected $data;
    protected $hash;
    protected $storeWorkerService;

    public function __construct($data)
    {
        $this->data = $data;
    }

    /**
     * 异步生成门店员工的二维码和太阳码
     */
    public function handle()
    {
        Log::info('JOB生成二维码海报开始:' . json_encode($this->data));
        // 获取添加成功的门店员工信息
        $memberList = Members::query()
            ->whereIn('id', $this->data['worker_ids'])
            ->where('enterprise_id', $this->data['enterprise_id'])
            ->get()
            ->toArray();
        if (empty($memberList)) {
            return false;
        }

        foreach ($memberList as $member) {
            $member['domain']        = $this->data['domain'];
            $member['operator_id']   = $this->data['operator_id'] ?? 0;
            $member['operator_name'] = $this->data['operator_name'] ?? '';
            // 异步预备员工原始联系我二维码
            $member['memQrCode'] = $member['qr_code'];
            // 生成企业微信联系我二维码
            $member['memUserId'] = $member['userid'];
            $member['type']      = $member['type'] != 4 ? 3 : 4;
            StoreWorkerService::generateWorkQrCode($member);
        }


//        $memUid = array_values(array_column($workerList, 'mem_uid'));
//        // 获取uc员工信息
//        $memList = StoreWorkerService::memberListNoAuth($this->data['domain'], 'contact', ['memUids' => $memUid,'pageSize' => 1000]);
//        $memList = array_combine(array_column($memList['list'], 'memUid'), $memList['list']);
//        foreach ($workerList as $worker) {
//            if (empty($memList[$worker['mem_uid']])) continue;
//            $member = $memList[$worker['mem_uid']];
//            $worker['domain'] = $this->data['domain'];
//            $worker['operator_id'] = $this->data['operator_id'] ?? 0;
//            $worker['operator_name'] = $this->data['operator_name'] ?? '';
//            // 异步预备员工原始的联系我二维码
//            $worker['memQrCode'] = $member['memQrCode'];
//            // 生成企业微信联系我二维码
//            $worker['memUserId'] = $member['memUserid'];
//            StoreWorkerService::generateWorkQrCode($worker);
//        }
        return true;
    }

}
