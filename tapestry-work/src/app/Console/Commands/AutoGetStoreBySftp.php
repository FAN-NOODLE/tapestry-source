<?php
/**
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/2/28
 * Time: 16:32
 */

namespace App\Console\Commands;

use App\Models\StoreDepartment;
use App\Models\StoreEmployee;
use App\Models\StoreSftpConfig;
use App\Models\StoreSftpErrorRecord;
use App\Models\StoreSftpLog;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Log;

class AutoGetStoreBySftp extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'auto:get_store_by_sftp';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '通过sftp获取三方门店及员工信息';

    protected $sshConn = null;

    /**
     * sftp中 员工的启用 和离职标识
     */
    const ADD = 1;
    const DEL = 0;


    public function handle()
    {
        $sftpConfig = StoreSftpConfig::query()->find(1);
        if (empty($sftpConfig)) {
            Log::info('未配置peoplehub sftp配置');

            return;
        }

        $conf          = [
            'host'     => $sftpConfig['host'],      //mftqa.tapestry.com,
            'port'     => $sftpConfig['port'],      //8022,
            'user'     => $sftpConfig['user'],      //Changyi,
            'password' => $sftpConfig['password'],  //T@pestry@123
        ];
        $this->sshConn = ssh2_connect($conf['host'], $conf['port']);

        if (!ssh2_auth_password($this->sshConn, $conf['user'], $conf['password'])) {
            echo 'sftp连接失败' . PHP_EOL;
            Log::info('sftp连接失败', $conf);
        }
        $sftpResource = ssh2_sftp($this->sshConn);
        $brandArr     = json_decode($sftpConfig['brand'], true); //{"CH":3353921300665408,"SW":3357415263026304}
        //读取文件
        $this->getCsvContent($brandArr, $sftpResource, $sftpConfig);

        ssh2_disconnect($this->sshConn);
    }


    private function getCsvContent($brandArr, $sftpResource, $sftpConfig)
    {

        $sftpFilePath = "ssh2.sftp://{$sftpResource}" . $sftpConfig['root_path'];
        if (!file_exists($sftpFilePath)) {
            echo '找不到远程sftp文件' . PHP_EOL;
            Log::info('找不到远程sftp目录', ['sftpFilePath' => $sftpFilePath]);
            return [];
        }
        $fileArr = scandir($sftpFilePath);
        foreach ($fileArr as $k => $v) {
            if ($v == '.' || $v == '..') {
                unset($fileArr[$k]);
            }
        }
        //筛选文件列表中未读文件、写入store_stfp_log表
        $existsNameArr = StoreSftpLog::query()->pluck('file_name')->toArray();
        //获取新增的文件
        $newFileArr = array_diff($fileArr, $existsNameArr);
        foreach ($newFileArr as $fileName) {
            $localFilePath = storage_path() . '/' . $fileName;
            $fileBrandArr  = explode("_", $fileName);
            if (empty($fileBrandArr[0])) {
                continue;
            }
            foreach ($brandArr as $brand => $enterpriseId) {
                if ($brand == $fileBrandArr[0]) {
                    $fileLog= StoreSftpLog::query()->create([
                        'file_name'       => $fileName,
                        'local_file_path' => $localFilePath,
                        'enterprise_id'   => $enterpriseId
                    ]);
                    if (copy($sftpFilePath . $fileName, $localFilePath)) {
                        //下载文件成功则直接处理文件
                        $this->saveCsv($localFilePath, $fileBrandArr[1], $enterpriseId, $fileLog->id);
                        continue;
                    }
                }
            }
            //处理完成把文件丢到三方规定的目录
            ssh2_sftp_rename($sftpResource, $sftpConfig['root_path'] . $fileName, $sftpConfig['recycle_path'] . $fileName);

        }
        return true;
    }


    //处理csv文件数据并入库
    private function saveCsv($localFilePath, $fileType, $enterpriseId, $fileLogId)
    {
        // 替换文件中的中文括号
        $str = str_replace('）',')',file_get_contents($localFilePath));
        file_put_contents($localFilePath,str_replace('（','(',$str));
        $handle = fopen($localFilePath, 'r');
        if (!$handle) {
            echo '未读取到sftp获取到的文件' . PHP_EOL;
            Log::info('未读取到sftp获取到的文件', ['localFilePath' => $localFilePath, 'enterpriseId' => $enterpriseId]);
            return [];
        }
        $errArr  = [];
        $dataArr = [];

        $i = 0;
        while (($row = fgetcsv($handle)) !== false) {
            // 跳过第一行标题
            if ($i === 0) {
                $i ++;
                continue;
            }
            // 过滤空行
            if (empty(array_filter($row)) || count($row) < 2) continue;
            //处理字符串前后空格
            $data = [];
            foreach ($row as $k => $v) {
                if ($fileType == 'Staff') {
                    if ($k > 9) continue;
                } else if ($fileType == 'Department') {
                    if ($k > 0) continue;
                }
//                $v = @iconv('gbk', 'utf-8', $v); //解决中文字符乱码问题
                $data[] = trim(ltrim($v, '0'));
            }

            // data 为每行的数据
            $dataArr[] = $data;
            if (count($dataArr) >= 1000) {
                $errArr = $this->saveData($fileType, $dataArr, $enterpriseId);
                $dataArr = [];
            }

        }
        $errArr = array_merge($this->saveData($fileType, $dataArr, $enterpriseId), $errArr);

        //处理有门店编码但是未创建的门店（部分门店编码通过导购数据附带生成）
        $this->createDepartmentByEmployee($enterpriseId);
        //更新日志
        $this->saveLog($enterpriseId, $localFilePath);
        //错误日志入库
        $this->saveErrorRecord($errArr, $fileLogId, $enterpriseId);
        fclose($handle);

    }

    //csv数据入数据库
    private function saveData($fileType, $dataArr, $enterpriseId)
    {
        $errArr = [];
        if ($fileType == 'Staff') {
            //导购文件
            $errArr = $this->addEmployee($dataArr, $enterpriseId);
        } else if ($fileType == 'Department') {
            //导购文件
            $errArr = $this->addDepartment($dataArr, $enterpriseId);
        }
        return $errArr;
    }

    private function addDepartment($fileData, $enterpriseId)
    {
        list($fileData, $errArr) = uniqueMultiArrayByField($fileData, [0]); // 0 代表门店编码
        foreach ($errArr as &$v) {
            array_push($v, '门店编码有重复');
        }
        $fileDeptCodeArr   = array_column($fileData, 0);
        $existDeptCodeArr  = StoreDepartment::query()
            ->where('enterprise_id',$enterpriseId)
            ->whereIn('third_store_code', $fileDeptCodeArr)
            ->pluck('third_store_code')
            ->toArray();
        $insertDeptCodeArr = array_diff($fileDeptCodeArr, $existDeptCodeArr);
        $now               = time();
        $departmentArr     = [];
        foreach ($fileData as $val) {
            $storeCode = $val[0];
            if (empty($storeCode)){
                $errArr[] = $val;
                continue;
            }
            if (in_array($storeCode, $insertDeptCodeArr)) {
                $departmentArr[] = [
                    'enterprise_id'    => $enterpriseId,
                    'third_store_code' => $storeCode,
                    'store_name'       => $storeCode,
                    'operate_name'     => 'system',
                    'created_at'       => $now,
                    'updated_at'       => $now
                ];
            }

        }
        StoreDepartment::query()->insert($departmentArr);
        return $errArr;
    }

    private function addEmployee($fileData, $enterpriseId)
    {
        list($fileData, $errArr) = uniqueMultiArrayByField($fileData, [3]); // 0 代表手机号 3代表员工编码
        foreach ($errArr as &$v) {
            array_push($v, '员工编码有重复或不合法');
        }
        $now                  = time();
        $fileEmployeeCodeArr  = array_column($fileData, 3);
        $existEmployeeCodeArr = StoreEmployee::query()
            ->where('enterprise_id',$enterpriseId)
            ->whereIn('employee_code', $fileEmployeeCodeArr)
            ->pluck('employee_code')
            ->toArray();
        //获取新增数据中在库中查不到的用户数据
        $insertEmployeeCodeArr = array_diff($fileEmployeeCodeArr, $existEmployeeCodeArr);
        //获取新增数据中在库中可以查到的用户数据
        $intersectEmployeeCodeArr = array_intersect($fileEmployeeCodeArr, $existEmployeeCodeArr);
        $addArr = [];

        //当信息无变动时不更新时间字段，可以用于判定是否有数据变更
        $updateStoreEmployeeModel = new StoreEmployee();
        $updateStoreEmployeeModel->timestamps = false;

        foreach ($fileData as $key => $val) {
            if (empty($val[3])) {
                //过滤没有员工编码的数据
                array_push($val, '员工编码不存在');
                $errArr[] = $val;
                continue;
            }
            if (!isset($val[8])) {
                //过滤没有员工编码的数据
                array_push($val, '员工在职/离职状态不存在');
                $errArr[] = $val;
                continue;
            }
            if ($val[8] == self::ADD) { //在职的员工必须要求手机号和姓名正常
                if (!preg_match("/^1[3456789]\d{9}$/", $val[0]) || empty($val[1])) {
                    array_push($val, '手机号或姓名不合法');
                    $errArr[] = $val;
                    //过滤没有电话、名称、员工编码的数据
                    continue;
                }
            }
            if (in_array($val[3], $insertEmployeeCodeArr)) {
                //过滤出之前不在业务架构中且是新增的用户数据
                $addArr[$key]['enterprise_id']    = $enterpriseId;
                $addArr[$key]['name']             = $val[1];
                $addArr[$key]['alias_name']       = $val[2];
                $addArr[$key]['employee_code']    = $val[3];
                $addArr[$key]['job_location']     = $val[4];
                $addArr[$key]['job_position']     = $val[5];
                $addArr[$key]['email']            = $val[7];
                $addArr[$key]['status']           = $val[8] == self::ADD ? StoreEmployee::STATUS_ENABLE : StoreEmployee::STATUS_DISABLE;
                $addArr[$key]['third_store_code'] = $val[9];
                $addArr[$key]['mobile']           = $val[0];
                $addArr[$key]['created_at']       = $now;
                $addArr[$key]['updated_at']       = $now;
                if ($val[6] == 'M') {
                    $addArr[$key]['gender'] = 1;
                } else if ($val[6] == 'F') {
                    $addArr[$key]['gender'] = 2;
                } else {
                    $addArr[$key]['gender'] = 0;
                }

            } else if (in_array($val[3], $intersectEmployeeCodeArr)) {
                $updateArr = [
                    'name'             => $val[1],
                    'alias_name'       => $val[2],
                    'employee_code'    => $val[3],
                    'job_location'     => $val[4],
                    'job_position'     => $val[5],
                    'email'            => $val[7],
                    'status'           => $val[8] == self::ADD ? StoreEmployee::STATUS_ENABLE : StoreEmployee::STATUS_DISABLE,
                    'mobile'           => $val[0]
                ];
                if (!empty($val[9])) {
                    $updateArr['third_store_code'] = $val[9];
                }
                //获取要编辑的用户立刻编辑
                $isUpdate = $updateStoreEmployeeModel->where('employee_code', $val[3])
                    ->where('enterprise_id',$enterpriseId)
                    ->update($updateArr);
                if ($isUpdate) {
                    //如果确实有数据更新再更新同步状态
                    StoreEmployee::query()->where('employee_code', $val[3])
                        ->where('enterprise_id',$enterpriseId)
                        ->update(['sync_status' => 0]);
                }
            }
        }
        StoreEmployee::query()->insert($addArr);
        return $errArr;
    }

    private function saveLog($enterpriseId, $filePath)
    {
        StoreSftpLog::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('local_file_path', $filePath)
            ->update(['is_success' => 1, 'success_at' => time()]);
    }

    private function createDepartmentByEmployee($enterpriseId)
    {
        $storeCodeArr = StoreEmployee::query()
            ->leftJoin('store_department', 'store_employee.third_store_code', '=', 'store_department.third_store_code')
            ->where('store_department.id', null)
            ->where('store_employee.status', StoreEmployee::STATUS_ENABLE)
            ->where('store_employee.enterprise_id',$enterpriseId)
            ->where('store_employee.third_store_code', '!=', '')
            ->pluck('store_employee.third_store_code')
            ->toArray();
        if (!empty($storeCodeArr)) {
            $addArr = [];
            $now    = time();
            $storeCodeArr = array_unique($storeCodeArr);
            foreach ($storeCodeArr as $storeCode) {
                $addArr[] = [
                    'enterprise_id'    => $enterpriseId,
                    'store_name'       => $storeCode,
                    'third_store_code' => $storeCode,
                    'operate_name'     => 'system',
                    'created_at'       => $now,
                    'updated_at'       => $now
                ];
            }
            StoreDepartment::query()->insert($addArr);
        }
    }

    private function saveErrorRecord($errArr, $fileLogId, $enterpriseId)
    {
        if (empty($errArr)) return ;
        $now = time();
        $addArr = [];
        foreach ($errArr as $val) {
            array_push($val, 0);
            $addArr[] = [
                'enterprise_id' => $enterpriseId,
                'file_log_id'   => $fileLogId,
                'data_json'      => json_encode($val),
                'created_at'       => $now,
                'updated_at'       => $now
            ];
        }
        StoreSftpErrorRecord::query()->insert($addArr);
    }


}
