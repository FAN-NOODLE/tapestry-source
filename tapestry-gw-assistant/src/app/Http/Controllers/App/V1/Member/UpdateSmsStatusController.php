<?php

namespace App\Http\Controllers\App\V1\Member;

use App\Http\Controllers\App\V1\BaseController as Controller;
use CrsSdk\Work\StoreWorker;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

/**
 * 店员修改接收短信状态
 * Class UpdateSmsStatus
 * @package App\Http\Controllers\App\V1\Member
 */
class UpdateSmsStatusController extends Controller
{
    public function __invoke(Request $request)
    {
        // 验证输入参数
        $rules = [
            'worker_id'   => 'required',
            'is_send_sms' => 'required|integer'
        ];

        $message   = [
            'worker_id.required'   => '员工ID不能为空',
            'is_send_sms.required' => '是否接受短信字段不能为空',
            'is_send_sms.integer'  => '是否接受短信字段类型错误'
        ];
        $input     = $this->validate($request, $rules, $message);
        $isSendSms = $request->input('is_send_sms');
        try {
            $httpParam['where'] = [
                'id'            => $input['worker_id'],
                'enterprise_id' => $this->enterpriseId
            ];
            $httpParam['data']  = ['is_send_sms' => $isSendSms];
            $ret                = (new StoreWorker())->update($httpParam);
            return success($ret);
        } catch (\Exception $e) {
            Log::info("改变是否接收短信提醒错误,错误信息为{{$e->getMessage()}}");
            error($e->getCode(), $e->getMessage());
        }
    }
}
