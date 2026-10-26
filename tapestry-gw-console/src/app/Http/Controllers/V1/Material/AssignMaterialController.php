<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/3/24
 * Time: 15:54
 */
namespace App\Http\Controllers\V1\Material;

use App\Http\Controllers\V1\BaseController;
use App\Traits\ExcelTrait;
use CrsSdk\MicroGoodsV1\Material;

class AssignMaterialController extends BaseController
{
    use ExcelTrait;

    /**
     * 模板表头
     */
    const EXCEL_TITLE = [
        '物料编码',
        '成本价(单位：元)',
        '状态(启用/禁用)'
    ];

    /**
     * 批量分派物料
     *
     */
    public function index()
    {
        if (!$this->request->hasFile('excel')) {
            notice(450000, ['未上传excel文件']);
        }
        $file = $this->request->file('excel');
        $sourceId = $this->request->input('source_id', $this->enterpriseId);

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
        $materialSdk = new Material();
        $result = [
            'success_count' => 0,
            'error_count'   => 0,
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
                    'source_id'     => $sourceId,
                    'material_no'   => $value[0],
                    'price'         => $value[1] * 100,
                    'status'        => $value[2]
                ];
                $materialSdk->assignMaterial($params);
                $result['success_count']++;
            } catch (\Exception $e) {
                $value = array_slice($value, 0, count(self::EXCEL_TITLE));
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
        $path = resource_path('static/物料分派模板.xlsx');

        return response()->download($path);
    }
}