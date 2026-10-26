<?php
/**
 * Created by PhpStorm.
 * User: Yang Hui
 * Date: 2020年10月20日17:9:30
 */

namespace App\Console\Commands;

use App\Services\EnterpriseService;
use Illuminate\Console\Command;
use App\Models\WithdrawalAutomaticReview;
use App\Models\WithdrawalApplication;
use App\Services\WithdrawalServices;
use Illuminate\Support\Facades\Log;

class WithdrawalAutoAudit extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'auto:withdrawal-auto-audit';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '导购提现自动审核';

    private $cmdTimestamp = '';

    public function handle()
    {
        Log::info('导购提现自动审核withdrawal-auto-audit:执行进来了');

        set_time_limit(0);
        $this->cmdTimestamp = time();
        //获取配置信息，并检查状态，如果为空则不执行
        $configs = WithdrawalAutomaticReview::query()
            ->where('auto_status', 1)
            ->get();
        Log::info('导购提现自动审核withdrawal-auto-audit:' . json_encode($configs));
        if ($configs->isEmpty()) {
            return;
        }
        foreach ($configs as $item) {
            $this->withdrawalAction($item);
        }
        Log::info('导购提现自动审核withdrawal-auto-audit:执行完成');
        return;
    }

    public function withdrawalAction($autoConfig)
    {
        Log::info('WithdrawalAutoAudit.withdrawalAction start, autoConfig:' . json_encode($autoConfig));
        //判断是否设置了时间，未设置不启动
        if (empty($autoConfig->auto_at)) {
            Log::info('导购提现自动审核withdrawal-auto-audit:未开启自动审批: empty auto_at');
            return false;
        }

        //判断设置时间日期是不是今天,如果是今天则不执行,第二天执行
        /*$updateDate = date('Y-m-d', $autoConfig->update_at);
        if (strtotime($updateDate) == strtotime(date('Y-m-d'))) {
            echo '未到执行时间1';
            return false;
        }*/

        //检测今天是否已经执行过
        $lastDate = date('Y-m-d', $autoConfig->last_at);
        if (strtotime($lastDate) == strtotime(date('Y-m-d'))) {
            Log::info('导购提现自动审核withdrawal-auto-audit:今天已经执行');
            return false;
        }

        //查询前一天
        $prevDayStartTime = date('Y-m-d', strtotime("-1 day"));

        $conditionStart = $prevDayStartTime . ' ' . $autoConfig->auto_at;
        $conditionEnd = date('Y-m-d ') . $autoConfig->auto_at;

        // 查询企业信息
        $enterpriseData = (new EnterpriseService())->getDetailByEnterpriseId($autoConfig->enterprise_id);
        Log::info('导购提现自动审核withdrawal-auto-audit， enterpriseData：' . json_encode($enterpriseData));
        if (empty($enterpriseData) || empty($enterpriseData['hash'])) {
            return false;
        }
        $enterpriseHash = $enterpriseData['hash'];

        //查询要执行的数据，每次10条
        WithdrawalApplication::query()
            ->where('enterprise_id', $autoConfig->enterprise_id)
            ->where('app_type', $autoConfig->mem_type)
            ->where('status', WithdrawalApplication::WAIT_APPLICATION_STATUS)
            ->whereBetween('request_initiated_at', [
                strtotime($conditionStart),
                strtotime($conditionEnd)
            ])
            ->chunkById(10, function ($items) use ($autoConfig, $enterpriseHash) {
                $params['ids'] = array_column($items->toArray(), 'id');
                $params['action'] = 'pass';
                $params['enterprise_id'] = $autoConfig->enterprise_id;
                (new WithdrawalServices())->withdrawalAudit($params, $enterpriseHash);
            });
        //更新最后一次更新时间
        $autoConfig->last_at = $this->cmdTimestamp;
        $autoConfig->save();
        Log::info('WithdrawalAutoAudit.withdrawalAction end');
        return true;
    }
}
