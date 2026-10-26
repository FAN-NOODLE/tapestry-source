<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/11/28
 * Time: 17:55
 */

namespace App\Console\Commands;

use App\Exceptions\Error;
use App\Jobs\SyncContactJob;
use App\Library\Uc\Member;
use App\Models\SyncOrgLog;
use App\Models\Task;
use App\Services\TaskService;
use CrsSdk\Boss\Auth;
use CrsSdk\Boss\Enterprise;
use Illuminate\Support\Carbon;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Log;

class AutoSyncContact extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'auto:sync_contact';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '每天定时一次自动同步通讯录';


    public function handle()
    {
        Log::info('定时同步不同企业通讯录开始', ['time' => Carbon::now()]);

        $data = (new Enterprise())->listsOfficial([]);
        if (!empty($data)) {
            foreach ($data['data'] as $value) {
                try {
                    //调UC同步通讯录
                    (new Member())->syncContact($value['domain']);
                } catch (Error $e) {
                    Log::info(
                        '定时同步不同企业通讯录异常',
                        ['ErrorInfo---' . $e->getMessage(), '$value---' . json_encode($value)]
                    );
                    continue;
                }
            }
        }

        Log::info('定时同步不同企业通讯录结束', ['time' => Carbon::now()]);
    }

}