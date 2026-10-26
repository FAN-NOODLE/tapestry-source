<?php
/**
 * 物料新增
 * Created by PhpStorm.
 * User: 李坚
 * Date: 2020/3/12
 * Time: 11:31
 */

namespace App\Http\Controllers\V1\Material;

use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\WarehouseStock;
use Illuminate\Support\Facades\Log;

class ChangeSaleTypeController extends BaseController
{


    /**
     * 更新物料配送方式
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index()
    {

        $params = $this->request->only([
            'material_id', // 物料id
            'sale_type', // 配送方式
            'warehouse_id', // 仓库id
        ]);
        $params['enterprise_id'] = $this->enterpriseId;
        $params['operate_id'] = $this->enterpriseAdminId;

        try{
            (new WarehouseStock())->updateSaleType($params);
        } catch (Error $error) {
            Log::error('获取商品列表失败', $error->data);
            notice(451508, [], $error->data);
        } catch (\Exception $exception) {
            Log::error('获取商品列表失败'.$exception->getMessage(), $params);
            notice(451508);
        }

        return success();

    }

}
