<?php
/**
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/2/28
 * Time: 16:32
 */

namespace App\Console\Commands;

use App\Library\Uc\Member;
use App\Models\StoreWorker;
use CrsSdk\Boss\Enterprise;
use Illuminate\Console\Command;

class UpdateStoreWorkerInfo extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'auto:update-store-worker-info';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '更新store_worker表的mem_userid';

    public function handle()
    {

        // 查询所有权限
        $workerList = StoreWorker::query()->where('mem_userid', '')
            ->where('worker_status', StoreWorker::IN_SERVICE)
            ->get(['mem_uid', 'enterprise_id'])
            ->groupBy('enterprise_id')
            ->toArray();
        if (empty($workerList)) {
            echo '无数据不能更新';
            return true;
        }
        foreach ($workerList as $enterpriseId => $val) {
            $memUidArr      = array_column($val, 'mem_uid');
            try{
                $enterpriseInfo = (new Enterprise())->detail(['enterprise_id' => $enterpriseId]);
            } catch (\Throwable $e) {
                echo '企业信息获取失败:'.$enterpriseId;
                continue;
            }
            try{
                $list = (new Member())->ignorePermissionsList($enterpriseInfo['data']['domain'], '', [
                    'memUids' => $memUidArr,
                    'pageNum'  => 1,
                    'pageSize' => count($memUidArr)
                ])['list'];
            } catch (\Throwable $e) {
                echo '员工信息获取失败:'.join(',', $memUidArr);
                continue;
            }
            $arr = array_column($list, 'memUserid', 'memUid');
            foreach ($arr as $memUid => $memUserId) {
                StoreWorker::query()->where('enterprise_id', $enterpriseId)
                    ->where('mem_uid', $memUid)
                    ->update(['mem_userid' => $memUserId]);
            }

        }
        die('处理完成');
    }

}
