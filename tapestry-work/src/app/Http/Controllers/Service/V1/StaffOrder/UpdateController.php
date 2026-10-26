<?php


namespace App\Http\Controllers\Service\V1\StaffOrder;

use App\Http\Controllers\Service\V1\BaseController;
use App\Http\Requests\StaffOrderRequest;
use App\Services\StaffManagerService;
use App\Services\StaffOrderService;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;

class UpdateController extends BaseController
{
    /**
     * 更新内购订单信息
     *
     * @param Request $request
     * @return array
     */
    public function index(Request $request)
    {
        $input = $request->post();

        // 更新内购订单信息
        StaffOrderRequest::update($input);

        try {

            DB::beginTransaction();

            if(!empty($input['out_refund_no'])){
                (new StaffOrderService())->updateStaffOrderRefund($input);
            }else{
                (new StaffOrderService())->updateStaffOrderPay($input);
            }

            DB::commit();

        } catch (\Exception $e) {
            Log::error("更新员工内购订单异常信息：", [
                'msg'  => $e->getMessage(),
                'code' => $e->getCode(),
                'file' => $e->getFile(),
                'line' => $e->getLine(),
                'request_data' => json_encode($input,JSON_UNESCAPED_UNICODE),
            ]);
            DB::rollBack();
        }

        return success(['data' => []]);
    }

}