<?php
/**
 * 统计不同状态下的订单个数
 * Created by PhpStorm.
 * User: gongguoxin
 * Date: 2019年11月7日
 * Time: 21:49:44
 */
namespace App\Http\Controllers\V1\NewOrder;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\NewOrder;

class StatisticsOrderController extends BaseController
{
    public function index()
    {
        $input = $this->request->input();
        $input['enterprise_id'] = $this->enterpriseId;
        $response = (new NewOrder())->statisticsOrder($input);
        return success($response);
    }
}