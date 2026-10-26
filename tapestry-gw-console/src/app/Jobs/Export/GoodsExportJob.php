<?php


namespace App\Jobs\Export;

use App\Common\Constant;
use App\Exceptions\Error;
use CrsSdk\MicroGoodsV1\Goods;
use Illuminate\Support\Facades\Log;

class GoodsExportJob extends ExportJob
{
    // sku列表导出
    public const EXPORT_CHANNEL_SKU = 1;
    // spu列表导出
    public const EXPORT_CHANNEL_SPU = 2;
    /**
     * @var string
     */
    public $queue = 'console-goods-export';

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
    public function __construct(array $params, array $paramJson = [], string $name = '商品列表导出')
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
        // 获取导出标题
        $header = $this->getExportHeader($this->params['exportChannel']);
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
            if ($params['exportChannel'] === self::EXPORT_CHANNEL_SKU) {
                return (new Goods())->getConsoleList($params);
            }
            return (new Goods())->spuDraftList($params);
        };
        // 查询对应的商品列表
        try {
            $goodsList = [];
            $params['page'] = 1;
            $params['size'] = 50;
            startGet: $result = $getList($params);
            if ($result && isset($result['list']) && $result['list']) {
                $goodsList = array_merge($goodsList, $result['list']);
//                $goodsList += $result['list'];
                $params['page'] ++;
                goto startGet;
            }
        } catch (\Exception $exception) {
            Log::info('获取商品列表失败' . $exception->getMessage(), $params);
            throw new Error(10000, '获取商品列表失败');
        }
        $this->formatExportList($params['exportChannel'], $goodsList);
        // 插入excel
        $this->exportExcel($this->data);
        $this->exportClose();
    }

    /**
     * 获取导出标题
     * @param int $exportChannel
     * @return array
     */
    private function getExportHeader(int $exportChannel): array
    {
        $headerSku = [
            '商品唯一编码',
            'SPU',
            'SKU',
            '商品名称（中文）',
            '商品URL',
            '规格',
            '商品价格',
            '线上库存',
            '线下库存',
            '销量',
            '配送方式',
            '渠道名称',
            '商品状态',
            '创建时间',
            '更新时间'
        ];
        $headerSpu = [
            '商品URL',
            'SPU',
            '商品名称',
            '商品价格',
            '线上库存',
            '线下库存',
            '销量',
            '配送方式',
            '在架SKU数量/总SKU数量',
            '商品状态',
            '创建时间',
            '更新时间'
        ];
        if ($exportChannel === self::EXPORT_CHANNEL_SKU) {
            return $headerSku;
        }
        return $headerSpu;
    }

    /**
     * 格式化数据
     * @param int $exportChannel
     * @param array $list
     */
    private function formatExportList(int $exportChannel, array $list): void
    {
        if ($exportChannel === self::EXPORT_CHANNEL_SKU) {
            $this->formatExportSku($list);
        } else {
            $this->formatExportSpu($list);
        }
    }

    /**
     * 格式化sku导出列表
     * @param array $list
     */
    private function formatExportSku(array $list): void
    {
        foreach ($list as $item) {
            // 规格
            $attributes = array_values(array_column($item['attributes'], 'attribute_value_name'));
            $attributes = implode(',', $attributes);
            // 配送方式
            $deliveryType = [];
            $deliveryType[] = $item['delivery']['is_express'] === 1 ? '配送' : null;
            $deliveryType[] = $item['delivery']['is_extract'] === 1 ? '自提' : null;
            $deliveryType = implode('/', array_filter($deliveryType));
            $this->data[] = [
                'material_no'    => $item['material_no'],
                'spu_id'         => $item['spu_id'],
                'sku_id'         => $item['sku_id'],
                'sku_name'       => $item['sku_name'],
                'page_url'       => $item['page_url'] ?? '',
                'attributes'     => $attributes,
                'on_line_price'  => $item['price'] / 100,
                'online_stock'   => $item['stock']['online_stock'],
                'offline_stock'  => $item['stock']['offline_stock'],
                'sales'          => $item['sales'],
                'delivery_type'  => $deliveryType,
                'channel_name'    => $item['channel_name'],
                'status'         => $item['status']['valueString'],
                'created_at'     => date('Y-m-d H:i:s', $item['time']['created_at']),
                'updated_at'     => date('Y-m-d H:i:s', $item['time']['updated_at']),
            ];
        }
    }

    /**
     * 格式化spu列表
     * @param array $list
     */
    private function formatExportSpu(array $list): void
    {
        foreach ($list as $item) {
            if ($item['price']['maxPrice'] === $item['price']['minPrice']) {
                $price = $item['price']['maxPrice'] / 100;
            } else {
                $price = $item['price']['minPrice'] / 100 . '~' . $item['price']['maxPrice'] / 100;
            }
            $deliveryType = [];
            if ($item['delivery']['is_express'] === 1) {
                $deliveryType[] = '快递';
            }
            if ($item['delivery']['is_extract'] === 1){
                $deliveryType[] = '自提';
            }
            $this->data[] = [
                // page url
                'page_url' => $item['page_url'] ?? '',
                // spu id
                'spu_id' => $item['id'],
                // 商品名称
                'spu_name' => $item['spu_name'],
                // 商品价格
                'price' => $price,
                // 线上库存
                'on_line_stock' => $item['stock']['online_stock'],
                // 线下库存
                'off_line_stock' => $item['stock']['offline_stock'],
                // 销量
                'sales' => $item['sales'],
                // 配送方式
                'delivery' => implode('/', $deliveryType),
                // sku数量
                'sku_count' => $item['count']['current'] . '/' . $item['count']['all'],
                // 商品状态
                'status' => $item['status']['valueString'],
                // 创建时间
                'created_at' => date('Y-m-d H:i:s', $item['created_at']),
                // 更新时间
                'updated_at' => date('Y-m-d H:i:s', $item['updated_at']),
            ];
        }
    }
}