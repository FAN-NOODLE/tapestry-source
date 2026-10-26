<?php
/**
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/2/28
 * Time: 16:32
 */

namespace App\Console\Commands;

use App\Library\Uc\Department;
use App\Models\Scrm\Departments;
use App\Models\StoreDepartment;
use App\Services\EnterpriseService;
use CrsSdk\Boss\Enterprise;
use Exception;
use Illuminate\Console\Command;

class SyncStoreDepartmentByCoach extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'sync-store-department-by-coach';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '反向将组织架构中的部门更新至业务架构（只是coach)';

    public function handle()
    {
        $enterpriseId = env('COACH_ENTERPRISE_ID');
        if (empty($enterpriseId)) {
            return ;
        }
        $enterpriseInfo = (new Enterprise())->detail(['enterprise_id' => $enterpriseId])['data'];
        $configList = EnterpriseService::getEnterpriseConfig($enterpriseId, StoreDepartment::STORE_PARENT_DEPT_ID_KEY);
        if (empty($configList['data']) || empty($configList['data'][0])) {
            return;
        }
        $key = redisKey('SYNC_STORE_INFO_BY_DEPARTMENT_TIME', ['enterprise_id' => $enterpriseId]);
        getOrSetData($key,time(), 24 * 3 * 3600);
        $storeParentId = $configList['data'][0]['value'];
        //获取已经存在的部门
        $storeDepartment = StoreDepartment::query()->where('enterprise_id', $enterpriseId)
            ->get(['wechat_department_id', 'third_store_code'])
            ->keyBy('wechat_department_id')->toArray();
        $allStoreIds = array_keys($storeDepartment);
        $departments = Departments::query()->where('enterprise_id', $enterpriseId)
            ->where('third_parent_id', $storeParentId)
            ->get(['name', 'third_id'])->keyBy('third_id')->toArray();
        $allDpIds = array_keys($departments);
        //需要新增的部门
        $addDpIdArr = array_diff($allDpIds, $allStoreIds);
        //需要更新的部门
        $updateDpIdArr = array_unique(array_intersect($allStoreIds, $allDpIds));
        //需要删除的部门
//        $delDpIdArr = array_diff($allStoreIds, $allDpIds);
//        if (!empty($delDpIdArr)) {
            //删除
//            StoreDepartment::query()->where('enterprise_id', $enterpriseId)
//                ->whereIn('wechat_department_id', $delDpIdArr)
//                ->delete();
//        }

        $time = time();
        $insertArr = [];
        //组装要插入的数据
        foreach ($departments as $dpId => $val) {
            $storeType = 0;
            if (strstr($val['name'], 'OCF')) {
                $storeType = 1;
            } elseif (strstr($val['name'], 'OC')) {
                $storeType = 2;
            }
            if (in_array($dpId, $addDpIdArr)) {
                //新增的
                $ucDpId = '';
                try {
                    $res = (new Department())->list($enterpriseInfo['domain'], '', ['thirds' => [$val['third_id']]]);
                    if (!empty($res['list'])) {
                        $ucDpId = current($res['list'])['dpId'];
                    }
                } catch (Exception $e) {
                }
                $insertArr[] = [
                    'wechat_department_id' => $val['third_id'],
                    'uc_department_id' => $ucDpId,
                    'store_name' => $val['name'],
                    'store_type' => $storeType,
                    'third_store_code' => $val['name'],
                    'operate_name' => 'system',
                    'enterprise_id' => $enterpriseId,
                    'sync_status' => StoreDepartment::SYNC_SUCCESS,
                    'sync_at' => $time,
                    'created_at' => $time,
                    'updated_at' => $time
                ];
            }

            if (in_array($dpId, $updateDpIdArr)) {
                //更新
                if ($storeDepartment[$dpId]['third_store_code'] != $val['name']) {
                    //如果名称有变动则更新
                    $tmpUpdateArr = [
                        'store_name' => $val['name'],
                        'store_type' => $storeType,
                        'third_store_code' => $val['name'],
                        'sync_status' => StoreDepartment::SYNC_SUCCESS,
                        'sync_at' => $time,
                    ];
                    StoreDepartment::query()->where('enterprise_id', $enterpriseId)
                        ->where('wechat_department_id', $dpId)
                        ->update($tmpUpdateArr);
                }

            }
        }
        StoreDepartment::query()->insert($insertArr);

        //处理已经新增但是没有企微id的数据
        $waitUpdateDpArr = StoreDepartment::query()->where('enterprise_id', $enterpriseId)
            ->where('wechat_department_id', 0)
            ->get();
        foreach ($waitUpdateDpArr as $val) {
            try {
                $res = (new Department())->list($enterpriseInfo['domain'], '', ['dpName' => $val->name]);
                if (!empty($res['list'])) {
                    $val->uc_department_id = $res['list'][0]['dpId'];
                    $val->wechat_department_id = $res['list'][0]['dpThirdid'];
                    $val->sync_status = StoreDepartment::SYNC_SUCCESS;
                    $val->sync_at = $time;
                    $val->save();
                }
            } catch (Exception $e) {

            }
        }
    }


}
