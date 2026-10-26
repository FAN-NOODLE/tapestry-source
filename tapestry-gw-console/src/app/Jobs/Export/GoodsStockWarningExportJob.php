<?php


namespace App\Jobs\Export;

use App\Common\Constant;
use App\Exceptions\Error;
use CrsSdk\MicroGoodsV1\Goods;
use Illuminate\Support\Facades\Log;

class GoodsStockWarningExportJob extends ExportJob
{
    /**
     * @var array
     */
    protected $data = [];

    /**
     * GoodsExportJob constructor.
     * @param array $params 用于处理业务得参数
     * @param array $paramJson 前端原始请求参数-传数组-下载中心转json
     * @param string $name 任务名称或导出文件sheet名称
     * @throws \Exception
     */
    public function __construct(array $params, array $paramJson = [], string $name = '商品库存预警导出')
    {
        parent::__construct($params, $paramJson, $name);
    }

    /**
     * 导出Excel列表handle
     * @return bool
     */
    public function handle()
    {
        $params = $this->params;
        $_SERVER['enterprise_hash'] = $this->params['enterprise_hash'];
        $header = [
            'SKU',
            '商品唯一编码',
            '商品名称',
            '售卖价',
            '配送方式',
            '线上库存',
            '线下库存',
            '线上库存预警值',
            '线下库存预警值',
        ];
        // 查询获取物料库存管理列表
        try {
            $this->excelExportInit($header)->listExport($params);
        } catch (\Exception $exception) {
            dump($exception->getMessage(), $exception->getTraceAsString());
            Log::info('导出失败' . $exception->getMessage(), $params);
            // 修改状态为失败
            $this->jobStatus(false, $exception->getMessage());
            return true;
        }
        // 导出完成 修改状态为成功
        $this->jobStatus();
        return true;
    }

    /**
     * 循环获取数据并插入excel
     * @param $params
     * @throws \Box\Spout\Common\Exception\IOException
     * @throws \Box\Spout\Common\Exception\InvalidArgumentException
     * @throws \Box\Spout\Writer\Exception\WriterNotOpenedException
     */
    private function listExport($params): void
    {

        $getList = static function(array $params): array {
            return (new Goods())->warningList($params);
        };
        // 查询对应的商品列表
        try {
            $goodsList = [];
            $params['page'] = 1;
            $params['size'] = 50;
            startGet: $result = $getList($params);
            if ($result && isset($result['list']) && $result['list']) {
                $goodsList = array_merge($goodsList, $result['list']);
                $params['page'] ++;
                goto startGet;
            }
        } catch (\Exception $exception) {
            Log::info('获取商品库存预警列表失败' . $exception->getMessage(), $params);
            throw new Error(10000, '获取商品库存预警列表失败');
        }
        $this->formatExportList($goodsList);
        // 插入excel
        $this->exportExcel($this->data);
        $this->exportClose();
    }

    /**
     * 格式化数据
     * @param array $list
     */
    private function formatExportList(array $list): void
    {
        foreach ($list as $item) {
            // 配送方式
            $deliveryType = [];
            $deliveryType[] = $item['delivery']['is_express'] === 1 ? '配送' : null;
            $deliveryType[] = $item['delivery']['is_extract'] === 1 ? '自提' : null;
            $deliveryType = implode('/', array_filter($deliveryType));
            $onlineWarning = $item['stock_warning']['online_stock_warning_value'];
            $onlineWarning = $onlineWarning === '-' ? '' :$onlineWarning;
            $offlineWarning = $item['stock_warning']['offline_stock_warning_value'];
            $offlineWarning = $offlineWarning === '-' ? '' :$offlineWarning;
            $this->data[] = [
                'sku_id'         => $item['sku_id'],
                'material_no'    => $item['material_no'],
                'sku_name'       => $item['sku_name'],
                'on_line_price'  => $item['price'] / 100,
                'delivery_type'  => $deliveryType,
                'online_stock'   => $item['stock']['online_stock'],
                'offline_stock'  => $item['stock']['offline_stock'],
                'online_stock_warning_value' => $onlineWarning,
                'offline_stock_warning_value'  => $offlineWarning
            ];
        }
    }

}