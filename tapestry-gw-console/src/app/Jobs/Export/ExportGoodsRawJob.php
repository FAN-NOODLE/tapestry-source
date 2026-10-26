<?php


namespace App\Jobs\Export;


use CrsSdk\MicroGoodsV1\Goods;
use Exception;
use Illuminate\Support\Facades\Log;
use Throwable;

class ExportGoodsRawJob extends ExportJob
{

    /**
     * @var string
     */
    public $queue = 'console-goods-export';

    /**
     * @var int 重试1次
     */
    public $tries = 1;

    public function __construct(array $params, array $paramJson = [], string $name = '商品信息')
    {
        parent::__construct($params, $paramJson, $name);
    }

    /**
     * @throws Exception
     */
    public function handle(): void
    {
        try {
            $_SERVER['enterprise_hash'] = $this->params['enterprise_hash'];
            $title = ['SPU_ID', 'SKU_ID', '物料编码', '商品名称', '标牌价', '售价', '线上库存', '线下库存'];
            $this->excelExportInit($title);
            $data = [];
            $enterpriseId = $this->params['enterprise_id'];
            $page = 1;
            $size = 500;
            startGet: $rawList = (new Goods())->getRawList([
                'enterprise_id' => $enterpriseId,
                'status'        => $this->params['status'] ?? 1,
                'page'          => $page,
                'size'          => $size
            ]);
            if ($rawList && $rawList['list']) {
                $data += $rawList['list'];
                ++ $page;
                goto startGet;
            }
            $this->exportExcel($data);
            $this->exportClose();
            $this->jobStatus(true);
        } catch (Throwable $exception) {
            Log::info('商品导出失败:' . $exception->getMessage() . ',失败trace:' . $exception->getTraceAsString());
            $this->jobStatus(false);
        }
    }
}