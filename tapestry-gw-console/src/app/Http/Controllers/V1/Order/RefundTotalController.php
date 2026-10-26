<?php
/**
 * Created by PhpStorm.
 * User: zq
 * Date: 19-4-8
 * Time: 下午6:20
 */

namespace App\Http\Controllers\V1\Order;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\Order;

class RefundTotalController extends BaseController
{
    /**
     * 售后头部统计
     * @return mixed
     * @throws \Exception
     */
    public function index()
    {
        $input['enterprise_id'] = $this->enterpriseId;
        $response = (new Order())->refundTotal($input);
        return success($response);
    }
}
