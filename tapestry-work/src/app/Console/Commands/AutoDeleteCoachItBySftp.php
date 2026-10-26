<?php /** @noinspection PhpComposerExtensionStubsInspection */

/**
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/2/28
 * Time: 16:32
 */

namespace App\Console\Commands;

use App\Library\Uc\Member;
use App\Models\CorporationSftpCoachLog;
use App\Models\CorporationSftpCoachWhite;
use App\Models\Scrm\Departments;
use App\Models\Scrm\MemberDepartments;
use App\Models\StoreSftpConfig;
use CrsSdk\Boss\Enterprise;
use Exception;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Log;

class AutoDeleteCoachItBySftp extends Command
{

    const ADD = 1;
    const DEL = 0;
    const THIRD_IDS = [340,307];
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'auto:delete_coach_it_by_sftp';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '通过sftp获取PH文件中的Coach办公室人员，帮需要离职的人员离职';

    protected $sshConn = null;



    public function handle()
    {
        $sftpConfig = StoreSftpConfig::query()->find(6);    //固定id配置
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
        //读取文件
        $this->getCsvContent($sftpResource, $sftpConfig);

        ssh2_disconnect($this->sshConn);
    }


    private function getCsvContent($sftpResource, $sftpConfig): void
    {

        $sftpFilePath = "ssh2.sftp://{$sftpResource}" . $sftpConfig['root_path'];
        if (!file_exists($sftpFilePath)) {
            echo '找不到远程sftp文件' . PHP_EOL;
            Log::info('找不到远程sftp目录', ['sftpFilePath' => $sftpFilePath]);
            return;
        }
        $fileArr = scandir($sftpFilePath);
        list($brand, $enterpriseId) = explode(',', $sftpConfig['brand']);
        foreach ($fileArr as $fileName) {
            if ($fileName == '.' || $fileName == '..') {
                continue;
            }
            $fileBrandArr  = explode("_", $fileName);
            if (count($fileBrandArr) < 3) {
                continue;
            }
            if ($fileBrandArr[0] != $brand) {
                continue;
            }
            if ($fileBrandArr[1] != 'Staff') {
                continue;
            }
            if ($fileBrandArr[2] != date('Ymd',strtotime('-1 days'))) {
                continue;
            }
            $localFilePath = storage_path() . '/' . $fileName;

             CorporationSftpCoachLog::query()->create([
                'file_name'       => $fileName,
                'local_file_path' => $localFilePath,
                'enterprise_id'   => $enterpriseId
            ]);
            if (copy($sftpFilePath . $fileName, $localFilePath)) {
                //下载文件成功则直接处理文件
                $this->saveCsv($localFilePath, $enterpriseId);
            }

        }
    }


    //处理csv文件数据并入库
    private function saveCsv($localFilePath, $enterpriseId): void
    {
        // 替换文件中的中文括号
        $str = str_replace('）',')',file_get_contents($localFilePath));
        file_put_contents($localFilePath,str_replace('（','(',$str));
        $handle = fopen($localFilePath, 'r');
        if (!$handle) {
            echo '未读取到sftp获取到的文件' . PHP_EOL;
            Log::info('未读取到sftp获取到的文件', ['localFilePath' => $localFilePath, 'enterpriseId' => $enterpriseId]);
            return;
        }
        $deleteUserCodeArr = [];

        $i = 0;
        while (($row = fgetcsv($handle)) !== false) {
            // 跳过第一行标题
            if ($i === 0) {
                $i ++;
                continue;
            }
            // 过滤空行
            if (empty(array_filter($row)) || count($row) < 11) continue;
            //过滤有门店的人
//            if (trim($row[9]) != '')  continue;
            //过滤门店编码为OC开头的人
            if (str_starts_with(trim($row[9]), "OC")) continue;

            //处理字符串前后空格
//              $v = @iconv('gbk', 'utf-8', $v); //解决中文字符乱码问题
            if ($row[8] == self::DEL) {
                $deleteUserCodeArr[] =  trim(ltrim($row[3], '0'));
            }
        }
        $this->formatData($deleteUserCodeArr, $enterpriseId);

        //更新日志
        $this->saveLog($enterpriseId, $localFilePath);
        fclose($handle);

    }

    //csv数据入数据库
    private function formatData($deleteUserCodeArr, $enterpriseId): bool
    {

        // 读取目前在职的办公室人员
        $thirdsIds = self::THIRD_IDS;
        $dpIdArr = Departments::query()
            ->where(function($query) use ($thirdsIds){
                return $query->whereIn("third_parent_id", $thirdsIds)
                    ->orWhereIn('third_id', $thirdsIds);
            })->pluck('id')->toArray();
        $list = MemberDepartments::query()
            ->with("members")
            ->whereIn('department_id', $dpIdArr)
            ->get()->toArray();
        $saCodeArr = [];
        foreach ($list as $val) {
            if ($val['members'] && $val['members']['deleted_at'] == 0) {
                $saCodeArr[] = $val['members']['userid'];
            }
        }
        $saCodeArr = array_unique($saCodeArr);

        Log::info("文件中要删除的用户为:".count($deleteUserCodeArr), $deleteUserCodeArr);
        Log::info("在职的办公室的用户集合:".count($saCodeArr), $saCodeArr);

        // 取交集
        $deleteUserCodeArr = array_values(array_intersect($deleteUserCodeArr, $saCodeArr));
        Log::info("办公室的人和要删除的人的交集:".count($deleteUserCodeArr), $deleteUserCodeArr);
        if (empty($deleteUserCodeArr)) return true;

        // 读取白名单取差集
        $whiteUserIdArr = CorporationSftpCoachWhite::query()->where('enterprise_id', $enterpriseId)
            ->pluck('userid')
            ->toArray();
        Log::info("白名单用户集合:".count($whiteUserIdArr), $whiteUserIdArr);

        $deleteUserCodeArr = array_values(array_diff($deleteUserCodeArr, $whiteUserIdArr));
        Log::info("最终要删除的用户为:".count($deleteUserCodeArr), $deleteUserCodeArr);
        if (empty($deleteUserCodeArr)) return true;

        //删除企微用户（离职）
        $enterprise = (new Enterprise())->detail(['enterprise_id' => $enterpriseId])['data'];
        $ucMember = new Member();
        foreach ($deleteUserCodeArr as $saCode) {
            try {
                $res = $ucMember->detailByUserId($enterprise['domain'], 'communalContact',
                    ['memUserid' => $saCode]
                );
            } catch (Exception $e) {
                Log::info("uc查不到该SA，工号:" . $saCode);
                continue;
            }

            if (!empty($res['memUid'])) {
                //如果是离职 则调用离职接口
                try {
                    $ucMember->delete($enterprise['domain'], 'communalContact',
                        ['memUid' => $res['memUid']]
                    );
                } catch (Exception $e) {
                    Log::info("uc删除该SA失败，工号:" . $saCode);
                    continue;
                }
                continue;
            }
            Log::info("uc删除SA成功，工号:" . $saCode);
        }

        return true;
    }


    private function saveLog($enterpriseId, $filePath)
    {
        CorporationSftpCoachLog::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('local_file_path', $filePath)
            ->update(['is_success' => 1, 'success_at' => time()]);
    }

}
