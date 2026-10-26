<?php

namespace App\Http\Controllers\V1\Material;

use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use App\Http\Requests\MaterialRequest;
use CrsSdk\MicroGoodsV1\Material;
use CrsSdk\MicroGoodsV1\Warehouse;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Validator;

class UpdateSingleStockController extends BaseController
{
    /**
     * 更新单仓物料库存
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $params = $this->request->only([
            'material_id',
            'material_no',
            'num', // 库存操作数量
            'operate_type', // 1: 增加 2：减少 3：覆盖
            'warehouse_id', // 来源id // 门店 商城
        ]);

        MaterialRequest::editStock($params);
        $warehouseSdk = new Warehouse();
        try {
            $params = [
                'enterprise_id' => $this->enterpriseId,
                'operate_id'    => $this->enterpriseAdminId,
                'operate_name'  => $this->enterpriseAdmin['username'],
                'warehouse_id'  => $params['warehouse_id'] ?? 0,
                'type'          => 2,
                'operate_type'  => $params["operate_type"],
                'material_no'   => $params["material_no"],
                'stock'         => $params["num"],
            ];
            $params = collect($params)->map(function ($value, $key) {
                return trim($value) == $value ? $value : trim($value);
            })->toArray();
            $warehouseSdk->saveMaterialStock($params);
        }catch (\Exception $exception) {
            Log::info('更新物料库存失败'.$exception->getMessage(), $params);
            notice(451516);
        }

        $this->writeOperatorLog("管理员:".$params['operate_name'].", ID：".$params['operate_id'].",更新物料库存");

        return success();
    }

}