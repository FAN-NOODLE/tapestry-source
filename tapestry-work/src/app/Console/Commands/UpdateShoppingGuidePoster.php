<?php

namespace App\Console\Commands;

use App\Models\StoreWorker;
use App\Services\MemberService;
use App\Services\StoreWorkerService;
use CrsSdk\Boss\Enterprise;
use Illuminate\Console\Command;

class UpdateShoppingGuidePoster extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'update_shopping_guide_poster';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '更新导购海报';

    /**
     * Create a new command instance.
     */
    public function __construct()
    {
        parent::__construct();
    }

    /**
     * Execute the console command.
     *
     * @return mixed
     * @throws \Exception
     */
    public function handle()
    {
        $workerList = StoreWorker::query()
            ->where('type',1)
            ->where('contact_qr_code_id', 0)
            ->get()
            ->toArray();

        if (empty($workerList)) {
            return true;
        }

        // 保存数据
        $this->savePoster($workerList);

        return true;
    }

    /**
     * 同步更新企业联系我二维码数 据
     * @param array $workerList
     *
     * @return bool
     * @throws \Exception
     */
    public function savePoster($workerList = [])
    {
        $enterpriseIds = array_unique(array_column($workerList, 'enterprise_id'));

        $workerData = [];
        foreach ($workerList as $item) {
            $workerData[$item['enterprise_id']][] = $item;
        }

        foreach ($enterpriseIds as $enterpriseId) {
            try {
                $enterprise = (new Enterprise())->detail(['enterprise_id' => $enterpriseId])['data'];
            }catch (\Exception $exception){
                continue;
            }
            $domain = $enterprise['domain'] ?? '';
            if (empty($domain)) {
                echo '企业不存在'.$enterpriseId;
                continue;
            }

            $newWorkerList = $workerData[$enterpriseId] ?? [];
            $memUid = array_values(array_column($newWorkerList, 'mem_uid'));
            try {
                // 获取uc员工信息
                $params = [
                    'memUids'  => $memUid,
                    'pageSize' => 1000
                ];
                $memList = MemberService::memberListNoAuth($domain, 'contact', $params);
            }catch (\Exception $exception){
                echo $exception->getMessage();
                continue;
            }
            $memList = array_combine(array_column($memList['list'], 'memUid'), $memList['list']);

            foreach ($workerList as $worker) {
                $workerNo = str_replace("Y", "G", $worker['worker_no']);
                StoreWorker::query()
                    ->where('id', $worker['id'])
                    ->update([
                        'worker_no' => $workerNo
                        ]
                    );
                $member = $memList[$worker['mem_uid']] ?? [];
                if (empty($member)) {
                    continue;
                }
                $worker['domain'] = $domain;
                // 异步预备员工原始的联系我二维码
                $worker['memQrCode'] = $member['memQrCode'];
                // 生成企业微信联系我二维码
                $worker['memUserId'] = $member['memUserid'];

                StoreWorkerService::generateOldWorkQrCode($worker);
            }
        }

        return true;
    }
}
