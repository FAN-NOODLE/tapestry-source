<?php


namespace App\Jobs\Export;

use App\Common\Constant;
use App\Exceptions\Error;
use App\Services\MaterialService;
use CrsSdk\MicroGoodsV1\Material;
use CrsSdk\MicroGoodsV1\MaterialChannel;
use Illuminate\Support\Facades\Log;
use mysql_xdevapi\Exception;

class MaterialExportJob extends ExportJob
{

    /**
     * MaterialExportJob constructor.
     * @param array $params 用于处理业务得参数
     * @param array $paramJson 前端原始请求参数-传数组-下载中心转json
     * @param string $name 任务名称或导出文件sheet名称
     */
    public function __construct(array  $params, array $paramJson = [], string $name = '统计数据导出')
    {
        parent::__construct($params, $paramJson, $name);
    }

    /**
     * 导出Excel列表handle
     * @return bool
     * @throws \Box\Spout\Common\Exception\IOException
     * @throws \Box\Spout\Common\Exception\InvalidArgumentException
     * @throws \Box\Spout\Writer\Exception\InvalidSheetNameException
     * @throws \Box\Spout\Writer\Exception\WriterNotOpenedException
     */
    public function handle()
    {
        $params = $this->params;

        $header = ['货品编号（SPU)', '物料编码(SKU)', '品名（中文）', '固定成本价(元）', '库存总量', '可分配库存', '创建时间', '更新时间'];

        // 查询获取物料库存管理列表
        try {
            $this->excelExportInit($header)->listExport($params);
        } catch (\Exception $exception) {
            Log::info('导出失败'.$exception->getMessage(), $params);
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
    private function listExport($params)
    {
        // 查询获取物料库存管理列表
        try {
            $result = (new Material())->assignChannelList($params);
        } catch (\Exception $exception) {
            Log::info('获取物料库存列表失败'.$exception->getMessage(), $params);
            throw new Error(10000, '获取物料库存列表失败');
        }

        // 计算总库存
        $data = [];
        if (isset($result['list']) && !empty($result['list'])) {
            foreach ($result['list'] as $item) {
                $data[]      = [
                    'spu_no'      => $item['spu_no'],
                    'material_no' => $item['material_no'],
                    'name'        => $item['name'],
                    'cost_price'  => $item['price'],
                    'total_stock' => $item['total_stock'],
                    'stock'       => $item['stock'],
                    'created_at'  => date('Y-m-d H:i:s', $item['created_at']),
                    'updated_at'  => date('Y-m-d H:i:s', $item['updated_at']),
                ];

            }
        }

        // 插入excel
        $this->exportExcel($data);

        // 判断是否取完数据
        if (ceil($result['total'] / $params['size']) > $params['page']) {
            $params['page']++;
            $this->listExport($params);
        }

        $this->exportClose();
    }

}