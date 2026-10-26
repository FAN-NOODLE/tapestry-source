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

class StaffGoodsDetailController extends BaseController
{
    /**
     * 内购商品详情
     * @return mixed
     * @throws \Exception
     * @author: DLong
     */
    public function index()
    {
        $all                  = $this->request->all();
        $all['enterprise_id'] = $this->enterpriseId;

        $response = (new StaffGoods())->staffGoodsDetail($all);
        $this->writeOperatorLog('内购商品详情页');

        return $response;
    }
}