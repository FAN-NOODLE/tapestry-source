<?php


namespace App\Services;


use App\Common\Constant;
use App\Jobs\Import\GoodsBatchChange;
use Illuminate\Http\UploadedFile;
use Illuminate\Redis\Connections\Connection;
use Illuminate\Support\Facades\Redis;
use Maatwebsite\Excel\Excel;
use Maatwebsite\Excel\Readers\LaravelExcelReader;
use PHPExcel_Worksheet;

class GoodsBatchChangeService
{
    /**
     * 无导出进程 可正常导出
     */
    public const STATUS_NO_PROCESS = 1;
    /**
     * 正在导出
     */
    public const STATUS_RUNNING = 1;
    /**
     * 导出完成
     */
    public const STATUS_COMPLETE = 2;
    /**
     * 当前没有导出或导出数据已删除
     */
    public const NO_STATUS = 3;

    /**
     * 未上架更改价格
     */
    public const TYPE_WAIT_SALES_PRICE = 1;

    /**
     * 更改商品名称（spu纬度）
     */
    public const TYPE_GOODS_NAME_SPU = 2;

    /**
     * 更改商品名称（sku纬度）
     */
    public const TYPE_GOODS_NAME_SKU = 3;

    /**
     * 待上架商品更改库存(快递)
     */
    public const TYPE_WAIT_SALES_STOCK_EXPRESS = 4;

    /**
     * 待上架商品更改库存(自提)
     */
    public const TYPE_WAIT_SALES_STOCK_EXTRACT = 5;

    /**
     * 已上架商品更改库存(快递)
     */
    public const TYPE_ALREADY_SALES_STOCK_EXPRESS = 6;

    /**
     * 已上架商品更改库存(自提)
     */
    public const TYPE_ALREADY_SALES_STOCK_EXTRACT = 7;

    /**
     * 上传类型的集合字符串
     */
    public const TYPE_COLLECT_STR =
        GoodsBatchChangeService::TYPE_WAIT_SALES_PRICE . ','
        . GoodsBatchChangeService::TYPE_GOODS_NAME_SPU . ','
        . GoodsBatchChangeService::TYPE_GOODS_NAME_SKU . ','
        . GoodsBatchChangeService::TYPE_WAIT_SALES_STOCK_EXPRESS . ','
        . GoodsBatchChangeService::TYPE_WAIT_SALES_STOCK_EXTRACT . ','
        . GoodsBatchChangeService::TYPE_ALREADY_SALES_STOCK_EXPRESS . ','
        . GoodsBatchChangeService::TYPE_ALREADY_SALES_STOCK_EXTRACT;

    /**
     * 更改表头
     */
    public const TYPE_TITLE_COLLECT = [
        self::TYPE_WAIT_SALES_PRICE => ['skuId', '商品价格'],
        self::TYPE_GOODS_NAME_SPU => ['spuId', '商品名称'],
        self::TYPE_GOODS_NAME_SKU => ['skuId', '商品名称'],
        self::TYPE_WAIT_SALES_STOCK_EXPRESS => ['skuId', '线上库存'],
        self::TYPE_WAIT_SALES_STOCK_EXTRACT => ['skuId', '线下库存'],
        self::TYPE_ALREADY_SALES_STOCK_EXPRESS => ['skuId', '线上库存'],
        self::TYPE_ALREADY_SALES_STOCK_EXTRACT => ['skuId', '线下库存'],
    ];

    /**
     * excel最大支持数量
     * @var int
     */
    private $excelMaxLength = 100000;

    /**
     * @var string excel文件绝对路径
     */
    private $realPath;

    /**
     * @var int
     */
    private $enterpriseId;

    /**
     * @var int
     */
    private $type;

    /**
     * @var Connection
     */
    private $redis;

    public function __construct(int $enterpriseId, int $type)
    {
        $this->enterpriseId = $enterpriseId;
        $this->type = $type;
        $this->redis = Redis::connection();
    }

    public function run(UploadedFile $file): bool
    {
        if ($this->canChange()) {
            $this->uploadExcel($file)->dispatchJob();
            return true;
        }
        return false;
    }

    /**
     * 是否可进行改价
     * @return bool
     */
    private function canChange(): bool
    {
        return !$this->redis->exists($this->getProgressSystemKey());
    }

    /**
     * 上传并保存excel
     * @param UploadedFile $file
     * @return $this
     */
    private function uploadExcel(UploadedFile $file): self
    {
        $fileName = md5($this->enterpriseId . microtime()) . '.' . $file->getClientOriginalExtension();
        $file->move(resource_path(), $fileName);
        $this->realPath = resource_path($fileName);
        $rawFileName = $file->getClientOriginalName();
        $this->parseExcel($this->realPath, $rawFileName);
        return $this;
    }

    /**
     * 解析并保存excel数据
     * @param string $path
     * @param string $rawFileName
     */
    private function parseExcel(string $path, string $rawFileName): void
    {
        $reader = app()->make(Excel::class);
        $reader->load($path, function (LaravelExcelReader $reader) use ($rawFileName) {
            // 读取选中的sheet 如果需要所有sheet在这里加层循环即可
            $reader->sheet($reader->sheetsSelected(), function (PHPExcel_Worksheet $sheet) use ($rawFileName) {
                $data = $sheet->toArray();
                foreach ($data as $key => $datum) {
                    if (!$datum[0]) {
                        unset($data[$key]);
                        continue;
                    }
                    array_walk($datum, static function(&$item) {
                        $item = trim($item);
                    });
                    $data[$key] = $datum;
                }
                $data = array_values($data);
                $title = array_shift($data);
                if (count($data) > $this->excelMaxLength) {
                    notice(2000001, [$this->excelMaxLength]);
                }
                if ($data) {
                    $this->saveExcelData($title, $data, $rawFileName);
                } else {
                    notice(2000003);
                }
            });
        });
    }

    /**
     * 保存excel数据
     * @param array $title
     * @param array $data
     * @param string $rawFileName
     */
    private function saveExcelData(array $title, array $data, string $rawFileName): void
    {
        if (isset(self::TYPE_TITLE_COLLECT[$this->type]) && $title === self::TYPE_TITLE_COLLECT[$this->type]) {
            // 存储数据
            $this->redis->set($this->getSaveListKey(), json_encode($data));
            // 初始化进度数据
            $this->redis->set($this->getProgressSystemKey(), json_encode($this->makeProgressData(count($data), 0, $rawFileName)));
            $this->redis->set($this->getProgressKey(), json_encode($this->makeProgressData(count($data), 0, $rawFileName)));
            // 初始化结果数据
            $this->redis->set($this->getResultKey(), json_encode($this->makeResultData(self::STATUS_RUNNING, [], [], count($data))));
        } else {
            notice(2000004);
        }
    }

    /**
     * 投递到队列
     */
    private function dispatchJob(): void
    {
        dispatch(new GoodsBatchChange($this->enterpriseId, $this->type, request()->header('enterprise-hash')));
    }

    /**
     * 根据类型获取存储数据的KEY/job会掉用
     * @return string
     */
    public function getSaveListKey(): string
    {
        return sprintf(Constant::GOODS_BATCH_CHANGE_SAVE_LIST, $this->enterpriseId, $this->type);
    }

    /**
     * 保存进度的key
     * @return string
     */
    public function getProgressKey(): string
    {
        return sprintf(Constant::GOODS_BATCH_CHANGE_PROGRESS, $this->enterpriseId, $this->type);
    }

    /**
     * 保存进度的key(系统使用)
     * @return string
     */
    public function getProgressSystemKey(): string
    {
        return sprintf(Constant::GOODS_BATCH_CHANGE_PROGRESS_SYSTEM, $this->enterpriseId, $this->type);
    }

    /**
     * 保存结果的key
     * @return string
     */
    public function getResultKey(): string
    {
        return sprintf(Constant::GOODS_BATCH_CHANGE_RESULT, $this->enterpriseId, $this->type);
    }

    /**
     * @param int $total
     * @param int $current
     * @param string $rawFileName
     * @return array
     */
    public function makeProgressData(int $total, int $current = 0, string $rawFileName = ''): array
    {
        $currentRawFileName = $rawFileName;
        if (!$rawFileName) {
            $processString = $this->redis->get($this->getProgressKey());
            if ($processString) {
                $jsonDecodeArray = json_decode($processString, true);
                if ($jsonDecodeArray && isset($jsonDecodeArray['file'])) {
                    $currentRawFileName = $jsonDecodeArray['file'];
                }
            }
        }
        return [
            'file'      => $currentRawFileName,
            'total'     => $total,
            'current'   => $current,
            'current_progress'  => $total === 0 ? 0 : (int)(round($current / $total, 2) * 100)
        ];
    }

    /**
     * 拼装结果数据结构
     * @param int $status
     * @param array $errorTitle
     * @param array $errorList
     * @param int $count
     * @return array
     */
    public function makeResultData(int $status, array $errorTitle = [], array $errorList = [], int $count = 0): array
    {
        $errorNum = count($errorList);
        $successNum = $count ? $count - $errorNum : 0;
        return [
            'is_complete'   => $status,
            'success_num'   => $successNum,
            'error_num'     => $errorNum,
            'error_title'   => $errorTitle,
            'error_data'    => $errorList,
        ];
    }

    public function __destruct()
    {
        @unlink($this->realPath);
    }
}