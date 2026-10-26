<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/4/18
 * Time: 15:54
 */

namespace App\Http\Controllers\V1\Warehouse;

use App\Common\EnterpriseConfigMenu;
use App\Http\Controllers\V1\BaseController;
use App\Traits\ExcelTrait;
use CrsSdk\MicroGoodsV1\Warehouse;

class ImportMaterialController extends BaseController
{
//    protected $requireLogin = false;
    use ExcelTrait;

    const SIMPLE_WAREHOUSE = 1; // 单仓
    const MULTI_WAREHOUSE = 2; // 多仓

    // 单仓导入处理方式：新增
    const TYPE_ADD = 1;
    // 单仓导入处理方式：修改
    const TYPE_UPDATE = 2;

    /**
     * 新增模板表头
     */
    const EXCEL_ADD_TITLE = [
        'SPU编码',
        '物料编码',
        '商品名称',
        '固定成本价（单位：元）',
        '标牌价（单位：元）',
        '库存总量',
        '配送方式（1.配送 2.自提 3.配送加自提）',
        '条码',
        '税目编码',
        '备注'
    ];

    /**
     * 修改模板表头
     */
    const EXCEL_UPDATE_TITLE = [
        'SPU编码',
        '物料编码',
        '商品名称',
        '固定成本价（单位：元）',
        '标牌价（单位：元）',
        '配送方式（1.配送 2.自提 3.配送加自提）',
        '条码',
        '税目编码',
        '备注'
    ];

    /**
     * 批量分派物料
     */
    public function index()
    {
        if (!$this->request->hasFile('excel')) {
            notice(100000, ['未上传excel文件']);
        }
        $file = $this->request->file('excel');
        $warehouseId = $this->request->input('warehouse_id');
        $operateType = $this->request->input('operate_type');
        $operatorId = $this->enterprise['id'] ?? 0;
        $operatorName = $this->enterprise['contact'] ?? '';

        if (empty($operateType)) {
            notice(100000, ['操作方式必填']);
        }

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
            notice(100000, ['批量处理最多支持1000条']);
        }

        // 判断excel模板是否被修改
        if ($operateType == self::TYPE_ADD) {
            $tmpTitle = self::EXCEL_ADD_TITLE;
//            if (array_slice($excelData[0], 0, count($tmpTitle)) !== $tmpTitle) {
//                return notice(100000, ['模板表头内容错误,请重新下载模板']);
//            }
        } else {
            $tmpTitle = self::EXCEL_UPDATE_TITLE;
//            if (array_slice($excelData[0], 0, count($tmpTitle)) !== $tmpTitle) {
//                return notice(100000, ['模板表头内容错误,请重新下载模板']);
//            }
        }


        // 开始处理数据
        $warehouseSdk = new Warehouse();
        $result = [
            'success_count' => 0,
            'error_count' => 0,
            'error_list' => [],
            'title' => array_merge($tmpTitle, ['失败原因'])
        ];

        foreach ($excelData as $key => $value) {
            // 过滤表头
            if ($key === 0) {
                continue;
            }

            try {
                $params = [
                    'enterprise_id' => $this->enterpriseId,
                    'operate_id' => $operatorId,
                    'operate_name' => $operatorName,
                    'warehouse_id' => $warehouseId,
                    'operate_type' => $operateType,
                    'spu_no' => $value[0],
                    'material_no' => $value[1],
                    'name' => $value[2],
                    'price' => (empty($value[3]) ? 0 : $value[3]) * 100,
                    'sign_price' => (empty($value[4]) ? 0 : $value[4]) * 100,
                ];
                // 表头节点
                $endKey = 4;
                if ($operateType == 1) {
                    $endKey += 1;
                    $params['stock'] = $value[$endKey];
                }
                $endKey += 1;
                $params['sale_type'] = empty($value[$endKey]) ? 1 : trim($value[$endKey]);  // 默认配送方式是配送
                $endKey += 1;
                $params['bar_code'] = $value[$endKey];
                $endKey += 1;
                $params['taxab_code'] = $value[$endKey];
                $endKey += 1;
                $params['desc'] = $value[$endKey];

                $params = collect($params)->map(function ($value, $key) {
                    return trim($value) == $value ? $value : trim($value);
                })->toArray();

                $warehouseSdk->saveBaseMaterial($params);
//                if ($operateType == 1) {
//                    $warehouseSdk->saveMaterialStock($params);
//                }

                $result['success_count']++;
            } catch (\Exception $e) {
                $value = array_slice($value, 0, count($tmpTitle));
                $value[] = $e->getMessage();
                $result['error_list'][] = $value;
                $result['error_count']++;
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
        $path = resource_path('static/物料导入模板.xlsx');

        return response()->download($path);
    }


    /**
     * 下载模板
     *
     * @return \Symfony\Component\HttpFoundation\BinaryFileResponse
     */
    public function downloadInsertTemplate()
    {
        $path = resource_path('static/添加物料信息模板.xlsx');

        return response()->download($path);
    }

    /**
     * 下载模板
     *
     * @return \Symfony\Component\HttpFoundation\BinaryFileResponse
     */
    public function downloadUpdateTemplate()
    {
        $path = resource_path('static/修改物料信息模板.xlsx');

        return response()->download($path);
    }

}