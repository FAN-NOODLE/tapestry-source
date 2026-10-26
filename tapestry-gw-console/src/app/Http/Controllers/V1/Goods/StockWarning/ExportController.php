<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/5/12 0012
 * Time: 14:55
 */

namespace App\Http\Controllers\V1\Goods\StockWarning;


use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use App\Jobs\Export\GoodsExportJob;
use App\Jobs\Export\GoodsStockWarningExportJob;
use CrsSdk\MicroGoodsV1\Activity as GoodsActivity;
use CrsSdk\MicroGoodsV1\Goods;

class ExportController extends BaseController
{
    /**
     * 库存预警列表
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index()
    {
        $input = $this->request->input();
        $input['enterprise_id']   = $this->enterpriseId;
        $input['enterprise_hash'] = $this->request->header('enterprise-hash');
        $input['route_url']       = $input['route_url'] ?? 'goods/stock-warning/export';
        $input['tag']             = $input['tag'] ?? 'stock-warning-export';
        $input['operate_name']    = $this->enterpriseAdmin['username'] ?? '未知';
        $input['operate_id']      = $this->enterpriseAdminId;
//        (new GoodsStockWarningExportJob($input, $input, '库存预警' . date('YmdHis', time())))->handle();die;
        $this->dispatch(new GoodsStockWarningExportJob($input, $input, '库存预警' . date('YmdHis', time())));
        return success([]);
    }

}
