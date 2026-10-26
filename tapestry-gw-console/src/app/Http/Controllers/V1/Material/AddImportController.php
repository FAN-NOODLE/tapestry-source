<?php

namespace App\Http\Controllers\V1\Material;

use App\Traits\ExcelTrait;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Warehouse;

class AddImportController extends BaseController
{
    use ExcelTrait;

    /**
     * 模板表头
     */
    const EXCEL_TITLE = [
        'SPU编码',
        '物料编码',
        '物料名称',
        '物料固定成本价（单位：元）',
        '标牌价（单位：元）',
        '物料数量',
        '批次号',
        '备注',
        '权重',
        '条码',
        '税目编码',
        '配送方式（1.配送 2.自提 3.配送加自提）',
    ];

    /**
     * 物料批量导入接口
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
        $file     = $this->request->file('excel');

        // 保存临时文件
        $saveName = uniqid('excel-') . '.' . $file->getClientOriginalExtension();
        $file->move(storage_path('logs/gw-console'), $saveName);
        $savePath = storage_path('logs/gw-console/' . $saveName);

        // 解析excel
        $excelData = $this->excelRead($savePath);

        // 删除临时文件
        unlink($savePath);

        if (count($excelData) > 500) {
            notice(450000, ['批量处理最多支持500条']);
        }

        // 判断excel模板是否被修改
        if (array_slice($excelData[0], 0, count(self::EXCEL_TITLE)) !== self::EXCEL_TITLE) {
            notice(450000, ['模板表头内容错误,请重新下载模板']);
        }

        // 开始处理数据
        $warehouseSdk = new Warehouse();
        $result      = [
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
//                    'warehouse_id'  => $warehouseId,
                    'type'          => 2,
                    'operate_type'  => 2,
                    'spu_no'        => $value[0],
                    'material_no'   => $value[1],
                    'name'          => $value[2],
                    'price'         => (empty($value[3]) ? 0 : $value[3]) * 100,
                    'sign_price'    => (empty($value[4]) ? 0 : $value[4]) * 100,
                    'stock'         => $value[5],
                    'batch_no'      => $value[6],
                    'desc'          => $value[7],
                    'weight'        => $value[8],
                    'bar_code'      => $value[9],
                    'taxab_code'    => $value[10],
                    'sale_type'    => empty($value[11]) ? 1 : trim($value[11]), // 默认配送方式是配送
                    'is_add'        => 1
                ];
                $params = collect($params)->map(function ($value, $key) {
                    return trim($value) == $value ? $value : trim($value);
                })->toArray();
                $warehouseSdk->importMaterial($params);

                $result['success_count']++;
            } catch (\Exception $exception) {
                $value                  = array_slice($value, 0, count(self::EXCEL_TITLE));
                $value[]                = $exception->getMessage();
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
        $path = resource_path('static/新建物料模板.xlsx');

        return response()->download($path);
    }
}
