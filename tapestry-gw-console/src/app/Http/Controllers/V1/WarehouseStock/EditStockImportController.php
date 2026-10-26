<?php

namespace App\Http\Controllers\V1\WarehouseStock;

use App\Common\EnterpriseConfigMenu;
use App\Traits\ExcelTrait;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Warehouse;

class EditStockImportController extends BaseController
{
    use ExcelTrait;

    const SIMPLE_WAREHOUSE = 1; // 单仓
    const MULTI_WAREHOUSE  = 2; // 多仓

    /**
     * 模板表头
     */
    const EXCEL_TITLE = [
        '物料编码',
        '物料数量',
        '批次号',
        '备注'
    ];

    /**
     * 物料库存导入编辑
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Box\Spout\Common\Exception\IOException
     * @throws \Box\Spout\Common\Exception\UnsupportedTypeException
     * @throws \Box\Spout\Reader\Exception\ReaderNotOpenedException
     */
    public function index()
    {
        if (!$this->request->hasFile('excel')) {
            notice(450000, ['未上传excel文件']);
        }
        $type = $this->request->input('type');
        $file = $this->request->file('excel');
        $warehouseId  = $this->request->input('warehouse_id');

        // 企业仓库类型
        $warehouseType = $this->enterpriseConfig(EnterpriseConfigMenu::WAREHOUSE_TYPE);
        if ($warehouseType == self::MULTI_WAREHOUSE) {
            if (empty($warehouseId)) {
                notice(100000, ['多仓版的仓库Id必填']);
            }
        }
        // 保存临时文件
        $saveName = uniqid('excel-') . '.' . $file->getClientOriginalExtension();
        $file->move(storage_path('logs/gw-console'), $saveName);
        $savePath = storage_path('logs/gw-console/' . $saveName);

        // 解析excel
        $excelData = $this->excelRead($savePath);

        // 删除临时文件
        unlink($savePath);

        if (count($excelData) > 1000) {
            notice(450000, ['批量处理最多支持1000条']);
        }

        // 判断excel模板是否被修改
        if (array_slice($excelData[0], 0, count(self::EXCEL_TITLE)) !== self::EXCEL_TITLE) {
            notice(450000, ['模板表头内容错误,请重新下载模板']);
        }

        // 开始处理数据
        $warehouseSdk = new Warehouse();
        $result = [
            'success_count' => 0,
            'error_list'    => [],
            'title'         => array_merge(self::EXCEL_TITLE, ['失败原因'])
        ];

        foreach ($excelData as $key => $value) {
            // 过滤表头
            if ($key === 0) {
                continue;
            }
            try {
                $params = [
                    'enterprise_id' => $this->enterpriseId,
                    'operate_id'    => $this->enterpriseAdminId,
                    'operate_name'  => $this->enterpriseAdmin['username'],
                    'way'           => 1,        // 1.内部调用  2. 外部系统导入
                    'operate_type'  => $type,
                    'material_no'   => $value[0],
                    'stock'         => $value[1],
                    'batch_no'      => $value[2],
                ];
                $params = collect($params)->map(function ($value, $key) {
                    return trim($value) == $value ? $value : trim($value);
                })->toArray();
                $warehouseSdk->saveMaterialStock($params);

                $result['success_count']++;
            } catch (\Exception $exception) {
                $value = array_slice($value, 0, count(self::EXCEL_TITLE));
                $value[] = $exception->getMessage();
                $result['error_list'][] = $value;
            }
        }

        return success($result);
    }

    /**
     * 下载模板
     *
     * @return \Symfony\Component\HttpFoundation\BinaryFileResponse
     */
    public function downloadAddTemplate()
    {
        $path = resource_path('static/物料库存编辑模板.xlsx');

        return response()->download($path);
    }

    /**
     * 下载模板
     *
     * @return \Symfony\Component\HttpFoundation\BinaryFileResponse
     */
    public function downloadUpdateTemplate()
    {
        $path = resource_path('static/修改物料库存模板.xlsx');

        return response()->download($path);
    }
}
