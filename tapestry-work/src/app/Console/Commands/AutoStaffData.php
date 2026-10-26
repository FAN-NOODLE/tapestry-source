<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/11/28
 * Time: 17:55
 */

namespace App\Console\Commands;

use App\Exceptions\Error;
use App\Library\Uc\Member;
use App\Services\StaffManagerService;
use CrsSdk\Boss\Enterprise;
use Illuminate\Support\Carbon;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Log;

class AutoStaffData extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'auto:sync_staff_data';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '定时更新不同企业所有内购数据';


    public function handle()
    {
        Log::info('定时更新不同企业所有内购数据--开始', ['time' => Carbon::now()]);

        $data = (new Enterprise())->listsOfficial();
        if (!empty($data)) {
            foreach ($data['data'] as $value) {
                try {
                    //兼容不同应用标识
                    $appIdentifier = 'contact';
                    if ($value['contact_status'] == 3) {
                        $appIdentifier = 'crs-shop';
                    }

                    //为内购员工自建联系人二维码
                    (new  StaffManagerService())->autoStaffData(
                        $value['id'],
                        $value['domain'],
                        $appIdentifier
                    );
                } catch (\Exception $e) {
                    Log::info(
                        '定时更新不同企业所有内购员工联系人二维码异常',
                        ['ErrorInfo---' . $e->getMessage(), '$value---' . json_encode($value)]
                    );
                    continue;
                }
            }
        }

        Log::info('定时更新不同企业所有内购员工联系人二维码--结束', ['time' => Carbon::now()]);
    }

}