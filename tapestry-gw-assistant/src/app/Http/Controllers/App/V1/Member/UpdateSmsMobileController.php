<?php
/**
 * Created by PHPStorm.
 * author: 刘朋波
 * Date: 2020/4/10
 * Time: 17:06
 */

namespace App\Http\Controllers\App\V1\Member;

use App\Common\Constant;
use App\Http\Controllers\App\V1\BaseController as Controller;
use App\Traits\CodeTrait;
use CrsSdk\Work\StoreWorker;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;
use Illuminate\Support\Facades\Log;

/**
 * 店员修改接收短信手机号
 * Class UpdateSmsStatus
 * @package App\Http\Controllers\App\V1\Member
 */
class  UpdateSmsMobileController extends Controller
{
    use CodeTrait;
    public function __invoke(Request $request, StoreWorker $storeWorker)
    {
        // 验证输入参数
        $rules = [
            'store_id'    => 'required',
            'sms_mobile' => 'required',
            'sms_mobile' => 'regex:/^1[0-9]{10}$/',
            'code'       => 'required'
        ];

        $message = [
            'store_id.required'    => '门店ID不能为空',
            'sms_mobile.required' => '手机号码不能为空',
            'sms_mobile.regex'    => '手机号码格式错误',
            'code.required'       => '验证码不能为空'
        ];
        $this->validate($request, $rules, $message);

        // 接收参数
        $code      = $request->input('code');
        $smsMobile = $request->input('sms_mobile');
        $storeId   = $request->input('store_id');

        // 判断验证码是否正确
        $resultCode = $this->getCodeCache($smsMobile, Constant::APP_CLIENT_NAME);

        if ($code != $resultCode) {
            error(652001);
        }

        // 更改店员接收短信的手机号码
        $enterpriseId = Auth::user()['enterprise_id'];
        //$enterpriseId = 2;
        $workerDetail = $this->storeWorkerDetail($enterpriseId, $storeId);
        try{
            $httpParam['where'] = [
                'id'            => $workerDetail['id'],
                'enterprise_id' => $enterpriseId,
                'store_id'      => $storeId
            ];
            $httpParam['data']  = ['sms_mobile' => $smsMobile];
            $ret                = (new StoreWorker())->update($httpParam);
            return success($ret);
        }catch (\Exception $e) {
            Log::info("修改接收短信提醒的手机错误,错误信息为{{$e->getMessage()}}");
            error($e->getCode(),$e->getMessage());
        }
    }
}
