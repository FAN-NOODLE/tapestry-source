<?php


namespace App\Console\Commands\Temp;

use Illuminate\Console\Command;
use App\Models\StoreWorker;
use App\Services\StoreWorkerService;
use App\Services\MemberService;
use Exception;
use Illuminate\Support\Facades\Log;
use App\Exceptions\Error;

class TempGenerateWorkerPoster extends Command
{
        /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'temp:generate-worker-poster {workerID} {domain}';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '临时手动生成导购海报';

    public function handle()
    {
        $workerID = $this->argument('workerID');
        $domain = $this->argument('domain');
        try {
            $worker = StoreWorkerService::detail([
                'worker_id'     => $workerID,
                'worker_status' => StoreWorker::IN_SERVICE
            ]);
            // 获取uc员工信息
            $params = [
                'memUids'  => [$worker['mem_uid'] ?? ''],
                'pageSize' => 100
            ];
            $memList                 = MemberService::memberListNoAuth($domain, 'contact', $params);
            $memList                 = array_combine(array_column($memList['list'], 'memUid'), $memList['list']);
            $member                  = $memList[$worker['mem_uid']];
            $worker['domain']        = $domain;
            $worker['operator_id']   = 1;
            $worker['operator_name'] = '系统';

            // 异步预备员工原始的联系我二维码
            $worker['memQrCode'] = $member['memQrCode'];
            // 生成企业微信联系我二维码
            $worker['memUserId'] = $member['memUserid'];
            StoreWorkerService::generateWorkQrCode($worker);
        } catch (\Exception $e) {
            echo '[TEMP]门店员工二维码生成异常：' . $e;
            throw new Error($e->getCode(), $e->getMessage(), [], 200);
        }
    }
}