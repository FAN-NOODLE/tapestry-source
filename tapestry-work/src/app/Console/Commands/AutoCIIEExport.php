<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/11/28
 * Time: 17:55
 */

namespace App\Console\Commands;

use App\Models\ActivityReserveRecord;
use Exception;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Redis;
use Prettus\Repository\Exceptions\RepositoryException;

class AutoCIIEExport extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'auto:ciie:export';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = 'CIIE吊牌预约活动数据导出';

    private $encryptionKey = "ciie_export";
    private $encryptionMethod = 'AES-256-CBC';
    private $exportPath; // 文件保存路径
    private $chunkSize = 1000; // 分批大小

    protected $id = 4; // 活动id

    /**
     * Create a new command instance.
     */
    public function __construct()
    {
        parent::__construct();

        // 写死文件保存路径
        $this->exportPath = storage_path('app/exports/ciie_activity');

        // 确保目录存在
        if (!is_dir($this->exportPath)) {
            mkdir($this->exportPath, 0755, true);
        }
    }

    /**
     * Execute the console command.
     *
     * @return string
     *
     * @throws Exception
     */
    public function handle(): string
    {
        try {
            $filePath = $this->exportToCsv();
            $this->info("导出成功！文件路径: {$filePath}");
            return 'ok';
        } catch (Exception $e) {
            $this->error("导出失败: " . $e->getMessage());
            return 'error';
        }
    }

    /**
     * AES加密手机号（可逆加密）
     */
    private function encryptMobile($mobile)
    {
        if (empty($mobile)) {
            return '';
        }

        $key = $this->getValidEncryptionKey();
        $iv = openssl_random_pseudo_bytes(16);
        $encrypted = openssl_encrypt($mobile, $this->encryptionMethod, $key, 0, $iv);

        return base64_encode($iv . $encrypted);
    }

    /**
     * 确保加密密钥为32位
     */
    private function getValidEncryptionKey()
    {
        if (strlen($this->encryptionKey) < 32) {
            return str_pad($this->encryptionKey, 32, '0', STR_PAD_RIGHT);
        }
        if (strlen($this->encryptionKey) > 32) {
            return substr($this->encryptionKey, 0, 32);
        }
        return $this->encryptionKey;
    }

    /**
     * 对应的解密函数
     */
    public static function decryptMobile($encryptedData, $encryptionKey)
    {
        if (empty($encryptedData)) {
            return '';
        }

        if (strlen($encryptionKey) < 32) {
            $encryptionKey = str_pad($encryptionKey, 32, '0', STR_PAD_RIGHT);
        } elseif (strlen($encryptionKey) > 32) {
            $encryptionKey = substr($encryptionKey, 0, 32);
        }

        $data = base64_decode($encryptedData);
        $iv = substr($data, 0, 16);
        $encrypted = substr($data, 16);

        return openssl_decrypt($encrypted, 'AES-256-CBC', $encryptionKey, 0, $iv);
    }

    /**
     * 导出CSV主函数
     */
    public function exportToCsv(): string
    {
        $filename = 'activity_reserve_export_' . date('Ymd_His') . '.csv';
        $filePath = $this->exportPath . '/' . $filename;

        $this->info("开始导出数据到: {$filePath}");

        $file = fopen($filePath, 'w');
        if (!$file) {
            throw new Exception("无法创建文件: {$filePath}");
        }

        // 添加BOM头解决Excel中文乱码
        fwrite($file, "\xEF\xBB\xBF");

        // CSV头部
        $headers = [
            'ID', '企业ID', '活动ID', 'UnionID', 'OpenID', '用户名称', '用户头像',
            '手机号(加密)', '性别', '二维码', '预约码', '预约时间', '是否核销',
            '核销时间', '导购ID', '导购名称', '版本号', '创建时间', '更新时间'
        ];

        fputcsv($file, $headers);

        $totalCount = 0;
        $batchCount = 0;

        // 使用chunk优化内存使用
        ActivityReserveRecord::where(["activity_id" => $this->id])
            ->chunk($this->chunkSize, function ($records) use ($file, &$totalCount, &$batchCount) {
                foreach ($records as $record) {
                    $row = $this->formatRecordToRow($record);
                    fputcsv($file, $row);
                    $totalCount++;
                }

                $batchCount++;
                $this->info("已处理第 {$batchCount} 批数据，当前总数: {$totalCount}");
            });

        fclose($file);

        $this->info("导出完成！总记录数: {$totalCount}");

        return $filePath;
    }

    /**
     * 格式化单条记录为CSV行
     */
    private function formatRecordToRow($record): array
    {
        $recordArray = $record->toArray();
        $encryptedMobile = $this->encryptMobile($recordArray['mobile']);
        return [
            $recordArray['id'],
            $recordArray['enterprise_id'],
            $recordArray['activity_id'],
            $recordArray['unionid'],
            $recordArray['openid'],
            $recordArray['user_name'],
            $recordArray['avatar'],
            $encryptedMobile,
            $recordArray['gender'] == 1 ? '男士' : '女士',
            $recordArray['qr_code'],
            $recordArray['activity_code'],
            date('Y-m-d H:i:s', $recordArray['reserve_time']),
            $this->getCheckStatusText($recordArray['is_check']),
            $recordArray['check_time'] ? date('Y-m-d H:i:s', $recordArray['check_time']) : '',
            $recordArray['member_id'],
            $recordArray['member_name'],
            $recordArray['version'],
            date('Y-m-d H:i:s', $recordArray['created_at']),
            date('Y-m-d H:i:s', $recordArray['updated_at'])
        ];
    }

    private function getCheckStatusText($status)
    {
        $statusMap = [
            0 => '未核销',
            1 => '已核销',
            2 => '已过期'
        ];

        return $statusMap[$status] ?? '未知';
    }
}