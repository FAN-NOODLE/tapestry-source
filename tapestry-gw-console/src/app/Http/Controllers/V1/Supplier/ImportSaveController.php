<?php


namespace App\Http\Controllers\V1\Supplier;

use App\Http\Controllers\V1\BaseController;
use App\Traits\ExcelTrait;
use CrsSdk\MicroGoodsV1\Supplier;

class ImportSaveController extends BaseController
{
    use ExcelTrait;

    /**
     * 模板表头
     */
    const EXCEL_TITLE = [
        '物料编码',
        '库存数量'
    ];

    public function index()
    {
        $this->validate($this->request, [
            'warehouse_id' => 'required|integer|min:1',
            'source_id'    => 'required|integer|min:1',
        ], [
            'warehouse_id.*' => '仓库ID不能为空',
            'source_id.*'    => '门店ID不能为空'
        ]);
        $warehouseId = $this->request->input('warehouse_id');
        $sourceId = $this->request->input('source_id');
        $isCheck = $this->request->input('is_check', 0);
        $supplierSdk = new Supplier();
        if ($isCheck == 0) {
            // 新增供货关系
            $supplierSdk->addSupplier([
                'enterprise_id' => $this->enterpriseId,
                'warehouse_id'  => $warehouseId,
                'source_id'     => $sourceId
            ]);
        }

        $result = [
            'success_count' => 0,
            'error_count'   => 0,
            'error_list'    => [],
            'title'         => array_merge(self::EXCEL_TITLE, ['失败原因'])
        ];

        if ($this->request->hasFile('excel')) {
            $file = $this->request->file('excel');
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

            foreach ($excelData as $key => $value) {
                // 过滤表头
                if ($key === 0) {
                    continue;
                }
                try {
                    $params = [
                        'enterprise_id' => $this->enterpriseId,
                        'source_id'     => $sourceId,
                        'warehouse_id'  => $warehouseId,
                        'material_no'   => $value[0],
                        'stock'         => $value[1],
                        'operate_type'  => 3, // 默认覆盖
                        'is_check'      => $isCheck
                    ];
                    $params = collect($params)->map(function ($value, $key) {
                        return trim($value) == $value ? $value : trim($value);
                    })->toArray();
                    $supplierSdk->saveSupplier($params);
                    $result['success_count']++;
                } catch (\Exception $e) {
                    $value = array_slice($value, 0, count(self::EXCEL_TITLE));
                    $value[] = $e->getMessage();
                    $result['error_list'][] = $value;
                    $result['error_count']++;
                }
            }
        }

        return success($result);
    }
}