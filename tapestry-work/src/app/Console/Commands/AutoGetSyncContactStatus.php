<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/11/28
 * Time: 17:55
 */

namespace App\Console\Commands;

use App\Library\Uc\Member;
use App\Models\SyncOrgLog;
use CrsSdk\Boss\Enterprise;
use Illuminate\Support\Carbon;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Log;

class AutoGetSyncContactStatus extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'auto:get_sync_contact_status';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '实时获取不同企业的同步状态';


    public function handle()
    {
        Log::info('实时获取不同企业的同步状态--开始', ['time' => Carbon::now()]);

        //正在同步的企业
        $epList     = SyncOrgLog::query()->where('status', SyncOrgLog::STATUS_UC_PROCESSING)
            ->get(['id', 'enterprise_id']);
        $data       = (new Enterprise())->listsOfficial();
        $domainList = array_column($data['data'], 'domain', 'id');
        foreach ($epList as $value) {
            echo '实时获取不同企业的同步状态--'. $value->enterprise_id.json_encode($domainList) . PHP_EOL;
            //调UC获取同步状态
            $ucResult = (new Member())->getSyncContactStatus($domainList[$value->enterprise_id]);
            if ($ucResult == 'ok') {
                $value->status = SyncOrgLog::STATUS_UC_OK;
                $value->save();
            }
        }

        Log::info('实时获取不同企业的同步状态--结束', ['time' => Carbon::now()]);
    }

}