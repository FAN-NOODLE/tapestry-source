<?php


namespace App\Jobs\Export;

use App\Common\Constant;
use App\Exceptions\Error;
use App\Services\MaterialService;
use CrsSdk\MicroGoodsV1\Material;
use CrsSdk\MicroGoodsV1\MaterialChannel;
use CrsSdk\Shop\WxMall;
use Illuminate\Support\Facades\Log;
use mysql_xdevapi\Exception;

class WxMallGoodsListExportJob extends ExportJob
{

    /**
     * WxMallGoodsListExportJob constructor.
     * @param array $params
     * @param array $paramJson
     * @param string $name
     * @throws \Box\Spout\Common\Exception\UnsupportedTypeException
     */
    public function __construct(array  $params, array $paramJson = [], string $name = '好物圈商品数据导出')
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

        $header = ['SPU', '商品名称', '所属分类', '好物圈分类'];

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
            $result = (new WxMall())->goodsList([
                'enterprise_id' => $params['enterprise_id'],
                'page'          => $params['page'] ?? 1,
                'size'          => $params['size'] ?? 1000,
            ]);
        } catch (\Exception $exception) {
            Log::info('获取物料库存列表失败'.$exception->getMessage(), $params);
            throw new Error(10000, '获取物料库存列表失败');
        }

        // 计算总库存
        $data = [];
        if (isset($result['list']) && !empty($result['list'])) {
            foreach ($result['list'] as $item) {
                $data[] = [
                    'id'               =>  "\t" . $item['id'] . "\t" ,
                    'spu_name'         => $item['spu_name'],
                    'category_name'    => implode(',', array_filter($item['category_name'])),
                    'wx_mall_category' => $item['wx_mall_category'],
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