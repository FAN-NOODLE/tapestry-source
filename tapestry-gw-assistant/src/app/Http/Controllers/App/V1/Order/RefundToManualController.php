<?php
/**
 * 一键转人工
 *
 * Created by PhpStorm.
 * User: huangxm
 * Date: 2020-04-28
 * Time: 12:42
 */

namespace App\Http\Controllers\App\V1\Order;

use App\Http\Controllers\App\V1\BaseController as Controller;
use CrsSdk\Shop\Refund;
use Illuminate\Http\Request;

class RefundToManualController extends Controller
{
    /**
     * @param Request $request
     * @param Refund $refund
     * @return \Illuminate\Http\JsonResponse|void
     */
    public function index(Request $request, Refund $refund)
    {
        $orderNo = $request->get('order_no', null);

        if (is_null($orderNo)) {
            error(654015);
        }

        $refund->RefundToManual(['order_no' => $orderNo]);
        return success();
    }
}
