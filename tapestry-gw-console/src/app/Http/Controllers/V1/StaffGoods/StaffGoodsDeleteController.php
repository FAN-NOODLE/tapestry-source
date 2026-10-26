<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/3/2 0002
 * Time: 下午 4:05
 */

namespace App\Http\Controllers\V1\StaffGoods;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\StaffGoods;

class StaffGoodsDeleteController extends BaseController
{
    /**
     * 删除内购商品
     * @return mixed
     * @throws \Exception
     * @author: DLong
     */
    public function index()
    {
        $all                  = $this->request->all();
        $all['enterprise_id'] = $this->enterpriseId;

        $response = (new StaffGoods())->staffGoodsDelete($all);
        $this->writeOperatorLog('内购商品列表页');

        return $response;
    }
}