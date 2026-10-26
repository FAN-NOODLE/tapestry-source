<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/4/18
 * Time: 15:54
 */

namespace App\Http\Controllers\V1\Warehouse;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Warehouse;

class AddMaterialController extends BaseController
{
    /**
     * 添加物料
     */
    public function index()
    {
        $warehouseId = $this->request->input('warehouse_id') ?? $this->enterpriseId;
        $operateType = $this->request->input('operate_type');
        $operatorId = $this->enterprise['id'] ?? 0;
        $operatorName = $this->enterprise['contact'] ?? '';
        $spuNo = $this->request->input('spu_no');
        $materialNo = $this->request->input('material_no');
        $name = $this->request->input('name') ?? '';
        $price = $this->request->input('price') ?? 0;
        $signPrice = $this->request->input('sign_price') ?? 0;
        $stock = $this->request->input('stock') ?? 0;
        $saleType = $this->request->input('sale_type') ?? 1;
        $barCode = $this->request->input('bar_code') ?? '';
        $taxabCode = $this->request->input('taxab_code') ?? '';
        $desc = $this->request->input('desc') ?? '';

        $params = [
            'enterprise_id' => $this->enterpriseId,
            'operate_id' => $operatorId,
            'operate_name' => $operatorName,
            'warehouse_id' => $warehouseId,
            'operate_type' => $operateType,
            'spu_no' => $spuNo,
            'material_no' => $materialNo,
            'name' => $name,
            'price' => $price * 100,
            'sign_price' => $signPrice * 100,
            'stock' => $stock,
            'sale_type' => $saleType,
            'bar_code' => $barCode,
            'taxab_code' => $taxabCode,
            'desc' => $desc,
        ];

        $params = collect($params)->map(function ($value, $key) {
            return trim($value) == $value ? $value : trim($value);
        })->toArray();

        // 开始处理数据
        $warehouseSdk = new Warehouse();

        $result = $warehouseSdk->saveBaseMaterial($params);
        return success($result);
    }
}