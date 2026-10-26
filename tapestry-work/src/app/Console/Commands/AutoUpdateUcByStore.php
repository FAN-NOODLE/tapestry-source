<?php
/**
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/2/28
 * Time: 16:32
 */

namespace App\Console\Commands;


use App\Library\Uc\Department;
use App\Library\Uc\Member;
use App\Models\StoreDepartment;
use App\Models\StoreEmployee;
use App\Models\StoreSftpConfig;
use App\Services\EnterpriseService;
use CrsSdk\Boss\Enterprise;
use CrsSdk\Boss\EnterpriseConfig;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Redis;

class AutoUpdateUcByStore extends Command
{

    const AUTO_EXEC_PAGE_SIZE = 100;
    /**
     * sftp中 员工的启用 和离职标识
     */
    const ON_JOB = 1;
    const LEAVE_JOB = 2;

    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'auto:update_uc_by_store';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '更新业务架构至企微';


    public function handle($enterpriseId = 0)
    {
        if ($enterpriseId > 0) {
            $enterprises = [(new Enterprise())->detail(['enterprise_id' => $enterpriseId])['data']];
        } else {
            // 查询企业列表
            $enterprises = (new Enterprise())->lists(['limit' => self::AUTO_EXEC_PAGE_SIZE])['data'];
        }
        //获取配置了people hub的企业id
        $sftpConfig = StoreSftpConfig::query()->find(1);
        $brandArr     = array_values(json_decode($sftpConfig['brand'], true)); //{"CH":3353921300665408,"SW":3357415263026304}
        // 循环企业
        $now = time();
        foreach ($enterprises as $enterprise) {
            if (!in_array($enterprise['id'], $brandArr)) {
                continue;
            }
            $key = redisKey('SYNC_STORE_INFO_BY_CHAT_TIME', ['enterprise_id' => $enterprise['id']]);
            Redis::set($key, $now);
            //获取父级部门id配置
            $configList = EnterpriseService::getEnterpriseConfig($enterprise['id'], StoreDepartment::STORE_PARENT_DEPT_ID_KEY);

            $sftpConfig = current($configList['data'] ?? []);
            if (empty($sftpConfig)) {
                continue;
            }
            $storeParentDeptId = $sftpConfig['value'];
            //读取未同步部门数据
            $deptList = $this->getNotSyncDepartmentData($enterprise['id']);
            if (!empty($deptList)) {
                //获取门店部门根目录对应的uc部门id
                $res = (new Department())->list($enterprise['domain'], '', ['thirds' => [$storeParentDeptId]]);
                if (!empty($res['list']) && !empty($ucDpId = current($res['list'])['dpId'])) {
                    //同步至企微
                    $this->handleDepartmentData($deptList, $enterprise, $ucDpId);
                } else {
                    echo '配置错误，uc中找不到该配置的部门' . PHP_EOL;
                    Log::info('配置错误，uc中找不到该配置的部门', ['storeParentDeptId' => $storeParentDeptId, 'enterpriseId' => $enterprise['id']]);
                }

            }
            //读取员工文件,每次读100个
            $page = 1;
            while (true) {
                $employeeData = $this->getNotSyncEmployeeData($enterprise['id'], $page);
                if (empty($employeeData)) break;
                $this->handleEmployeeData($employeeData, $enterprise);
                $page++;
            }

        }

    }

    private function getNotSyncDepartmentData($enterpriseId)
    {
        $newDeptList = StoreDepartment::query()->where('enterprise_id', $enterpriseId)
            ->where('sync_status', 0)
            ->get(['third_store_code', 'store_name', 'uc_department_id', 'wechat_department_id'])
            ->toArray();

        return $newDeptList;
    }

    private function getNotSyncEmployeeData($enterpriseId, $page = 1)
    {
        $offset      = ($page - 1) * self::AUTO_EXEC_PAGE_SIZE;
        $newDeptList = StoreEmployee::query()->where('enterprise_id', $enterpriseId)
            ->where('sync_status', 0)
            ->limit(self::AUTO_EXEC_PAGE_SIZE)->offset($offset)
            ->get(['name', 'alias_name', 'employee_code', 'status',
                   'job_location', 'job_position', 'gender',
                   'email', 'mobile', 'third_store_code', 'wechat_user_id'])
            ->toArray();

        return $newDeptList;
    }

    private function handleEmployeeData($employeeData, $enterprise)
    {
        $storeCodeArr = array_column($employeeData, 'third_store_code');

        $storeDepartmentArr = StoreDepartment::query()->whereIn('third_store_code', array_filter($storeCodeArr))
            ->get(['wechat_department_id', 'uc_department_id', 'third_store_code'])
            ->keyBy('third_store_code')
            ->toArray();
        foreach ($storeDepartmentArr as $key => $val) {
            $storeDepartmentArr[strtolower($key)] = $val;
        }
        foreach ($employeeData as $employee) {
            if (empty($employee['third_store_code'])) continue;
            $employee['third_store_code'] = strtolower($employee['third_store_code']);
            if (empty($storeDepartmentArr[$employee['third_store_code']]) && $employee['status'] == self::ON_JOB) {
                echo '数据错误，该门店code找不到信息' . PHP_EOL;
                Log::error('数据错误，该门店code找不到信息', ['storeParentDeptId' => $employee['third_store_code'], 'enterpriseId' => $enterprise['id']]);
                continue; //如果还没有企微部门 则跳过
            }

            try {
                $res = (new Member())->detailByUserId($enterprise['domain'], 'communalContact',
                    ['memUserid' => $employee['employee_code']]
                );
            } catch (\Exception $e) {

            }
            try {
                $addRes = [];
                if (empty($res) || !is_array($res)) {
                    if ($employee['status'] == self::ON_JOB) {
                        //同步添加企微人员
                        $addRes = (new Member())->add($enterprise['domain'], 'communalContact', [
                            'memUsername' => $employee['name'],
                            'dpIdList'    => [$storeDepartmentArr[$employee['third_store_code']]['uc_department_id']],
                            'memGender'   => $employee['gender'],
                            'memMobile'   => $employee['mobile'],
                            'memUserid'   => $employee['employee_code'],
                            'memJob'      => $employee['job_position'],
                            'memAlias'    => $employee['alias_name'],
                            "memEmail"    => $employee['email'],
                        ]);
                    }
                } else {
                    if ($employee['status'] == self::ON_JOB) {
                        //如果是在职则编辑用户的信息或者部门
                        (new Member())->edit($enterprise['domain'], 'communalContact', [
                            'memUid'      => $res['memUid'],
                            'memUsername' => $employee['name'],
                            'memGender'   => $employee['gender'],
                            'memMobile'   => $employee['mobile'],
                            'memUserid'   => $employee['employee_code'],
                            'memJob'      => $employee['job_position'],
                            'memAlias'    => $employee['alias_name'],
                            "memEmail"    => $employee['email'],
                            'dpIdList'    => [$storeDepartmentArr[$employee['third_store_code']]['uc_department_id']]
                        ]);
                    } else {
                        if (!empty($res['memUid'])) {
                            //如果是离职 则调用离职接口
                            (new Member())->delete($enterprise['domain'], 'communalContact',
                                ['memUid' => $res['memUid']]
                            );
                        }
                    }

                }

                if (isset($addRes['memUserid']) || isset($res)) {
                    if ($employee['status'] == self::ON_JOB) {
                        //更新员工企微user_id
                        $updateArr['wechat_user_id'] = $addRes['memUserid'] ?? $res['memUserid'];
                    }
                    $updateArr['sync_status'] = StoreEmployee::SYNC_SUCCESS;
                    $updateArr['sync_at'] = time();

                    StoreEmployee::query()->where('enterprise_id', $enterprise['id'])
                        ->where('employee_code', $employee['employee_code'])
                        ->update($updateArr);
                }

            } catch (\Exception $e) {
                $errorArr = [
                    'msg'           => 'UC创建/更新企微人员信息错误' . $e->getMessage(),
                    'enterprise_id' => $enterprise['id'],
                    'memUsername'   => $employee['name'],
                    'dpIdList'      => [$storeDepartmentArr[$employee['third_store_code']]['uc_department_id']],
                    'memGender'     => $employee['gender'],
                    'memMobile'     => $employee['mobile'],
                    'memUserid'     => $employee['employee_code'],
                    'memJob'        => $employee['job_position'],
                    'memAlias'      => $employee['alias_name']
                ];
                Log::error('UC错误：', $errorArr);

            }

        }
    }

    private function handleDepartmentData($deptData, $enterprise, $ucDpId)
    {
        foreach ($deptData as $deptInfo) {
            $updateArr = []; // 需要更新的数据
            $dpName    = mb_substr($deptInfo['store_name'], 0, 31);
            //同步企微
            try {
                if (!empty($deptInfo['uc_department_id'])) {
                    $addRes = (new Department())->edit($enterprise['domain'], 'communalContact', [
                        'dpId'   => $deptInfo['uc_department_id'],
                        'dpName' => $dpName ?: $deptInfo['third_store_code']
                    ]);
                } else {
                    $addRes                        = (new Department())->add($enterprise['domain'], 'communalContact', [
                        'dpName'     => $dpName ?: $deptInfo['third_store_code'],
                        'dpParentid' => $ucDpId
                    ]);
                    $updateArr['uc_department_id'] = $addRes['dpId'] ?? 0; //企微部门id
                    if (!empty($addRes['dpId'])) {
                        $deptData = (new Department())->list($enterprise['domain'], '', ['dpIds' => [$addRes['dpId']]]);
                        if (!empty($deptData['list'])) {
                            $updateArr['wechat_department_id'] = current($deptData['list'])['dpThirdid']; // uc部门id
                        }

                    }
                }
                Log::info('创建企微部门', $addRes);
            } catch (\Exception $e) {
                $errorArr = [
                    'msg'           => 'UC创建部门错误' . $e->getMessage(),
                    'enterprise_id' => $enterprise['id'],
                    'dpName'        => $deptInfo['third_store_code'],
                    'dpParentid'    => $ucDpId
                ];
                Log::error('UC错误：', $errorArr);
                continue;
            }
            //更新企微部门id至数据表中
            StoreDepartment::query()->where('enterprise_id', $enterprise['id'])
                ->where('third_store_code', $deptInfo['third_store_code'])
                ->update(array_merge($updateArr, [
                    'sync_at'     => time(),
                    'sync_status' => StoreDepartment::SYNC_SUCCESS
                ]));
        }
    }


}
