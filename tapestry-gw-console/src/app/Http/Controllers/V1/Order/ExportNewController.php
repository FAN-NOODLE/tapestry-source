<?php
/**
 * 异步导出订单
 * User: 龚国新
 * Date: 2020年2月6日
 * Time: 10:54:50
 */

namespace App\Http\Controllers\V1\Order;

use App\Common\Constant;
use CrsSdk\Shop\Order;
use Illuminate\Http\Request;
use App\Http\Controllers\V1\BaseController;

class ExportNewController extends BaseController
{

    /**
     * 异步导出筛选的订单
     *
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $input                  = $request->all();
        $input['param_json']    = $request->all();
        $input['enterprise_id'] = $this->enterpriseId;
        $input['operator_id']   = $this->enterpriseAdmin['id'];
        $input['operator_name'] = $this->enterpriseAdmin['username'];

        $data = (new Order())->exportExcel($input);

        return success();
    }

}
