<?php
/**
 * 内购订单统计
 * User: liurui
 * Date: 2018/11/26
 * Time: 14:10
 */

namespace App\Http\Controllers\V1\StaffOrder;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\StaffOrder;
use Illuminate\Http\Request;

class StatisticsController extends BaseController
{

    /**
     * 订单列表接口
     * {domain}/app/v1/order/list
     *
     * @param Request $request
     *
     * @return mixed|\Psr\Http\Message\StreamInterface
     */
    public function index(Request $request)
    {
        $input = $request->only(['search', 'page', 'size']);
        $input['search']['enterprise_id'] = $this->enterpriseId;

        return  (new StaffOrder())->statistics($input);
    }
}
