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
use CrsSdk\MicroGoodsV1\Activity as GoodsActivity;
use CrsSdk\MicroGoodsV1\Goods;

class ListController extends BaseController
{
    /**
     * 库存预警列表
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index()
    {
        $param = $this->request->input();
        $param['enterprise_id'] = $this->enterpriseId;
        $result = (new Goods())->warningList($param);
        return success($result);
    }

}
