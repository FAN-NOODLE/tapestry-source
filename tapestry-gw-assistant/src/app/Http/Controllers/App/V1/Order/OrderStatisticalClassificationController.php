<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/8/18 0008
 * Time: 17:37
 */

namespace App\Http\Controllers\App\V1\Order;

use CrsSdk\Work\OrderStatistical;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class OrderStatisticalClassificationController extends OrderBaseController
{
    /**
     * 订单分导购和合伙人统计
     *
     * @param Request $request
     *
     * @return \Illuminate\Http\JsonResponse
     * @author  HuiYang
     * @date    2020年11月18日20:29:3
     */
    public function index(Request $request)
    {
        $param['start_time']     = $request->input('start_time');
        $param['end_time']       = $request->input('end_time');
        $param['store_id']       = $request->input('store_id', 0);
        //$param['worker_ids']     = $request->input('worker_ids');
        //$param['store_ids']      = $request->input('store_ids');
        //$param['department_ids'] = $request->input('department_ids');
        //$param['enterprise_id']  = $this->enterpriseId;
        //$param['qy_domain']      = $this->loginAdminEnterprise['domain'];
        $param['admin_id']       =  $this->loginAdminEnterprise['admin_id'];

        $result = (new OrderStatistical())->orderStatistical($param);

        return success($result);


    }
}
