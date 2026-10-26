<?php
/**
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/2/28
 * Time: 16:32
 */

namespace App\Console\Commands;

use App\Common\Constant;
use App\Services\GuideManagerService;
use CrsSdk\Boss\Enterprise;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Log;

class AutoContactMember extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'auto:sync_guide_data';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '自动更新通讯录员工分享链路';

    public function handle()
    {

        $data            = (new Enterprise())->listsOfficial();

        if (!empty($data)) {
            foreach ($data['data'] as $value) {
                if ($value['version']==Constant::NORMAL_VERSION){
                    continue;
                }

                try {
                    //兼容不同应用标识
                    $appIdentifier = 'contact';
                    if ($value['contact_status'] == 3) {
                        $appIdentifier = 'crs-shop';
                    }

                    // 同步导购数据以及分享
                    (new  GuideManagerService())->autoSyncData(
                        $value['id'],
                        $value['domain'],
                        $appIdentifier
                    );

                } catch (\Exception $e) {
           //         echo '错误行数:'. $e->getLine();
              //      echo '错误文件:'. $e->getFile();
              //      echo '原因:'. $e->getMessage();
                //    echo json_encode($value);
                    echo json_encode($value);
                    Log::info(
                        '定时任务更新导购数据分享关系',
                        ['ErrorInfo---' . $e->getMessage(), 'value---' . json_encode($value)]
                    );
                    continue;
                }
            }
        }
    }


}
