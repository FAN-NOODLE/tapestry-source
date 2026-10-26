<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/3/2 0002
 * Time: 下午 3:33
 */

namespace App\Http\Controllers\V1\StaffGoods;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\StaffGoods;

class StaffGoodsIncreaseController extends BaseController
{
    /**
     * 新增内购商品
     * @return mixed
     * @throws \Exception
     * @author: DLong
     */
    public function index()
    {
        $all                  = $this->request->all();
        $all['enterprise_id'] = $this->enterpriseId;

        $response = (new StaffGoods())->staffGoodsIncrease($all);
        $this->writeOperatorLog('内购商品保存页');

        return $response;
    }
}