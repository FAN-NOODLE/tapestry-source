<?php


namespace App\Http\Controllers\Service\V1\StaffOrder;

use App\Http\Controllers\Service\V1\BaseController;
use App\Http\Requests\StaffOrderRequest;
use App\Services\StaffOrderService;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\DB;

class AddController extends BaseController
{
    /**
     * 新增员工内购订单
     *
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request)
    {
        $input = $request->post();

        // 验证参数
        StaffOrderRequest::add($input);
        try {

            DB::beginTransaction();

            // 插入数据
            (new StaffOrderService())->createStaffOrder($input);

            DB::commit();

        } catch (\Exception $e) {

            Log::error("新增员工内购订单异常信息：", [
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