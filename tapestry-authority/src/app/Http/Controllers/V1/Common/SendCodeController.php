<?php


namespace App\Http\Controllers\V1\Common;


use App\Http\Controllers\V1\BaseController;
use App\Models\EnterpriseAdmin;
use App\Services\SendCodeService;
use App\Traits\CodeTrait;
use Illuminate\Support\Facades\Redis;
use Illuminate\Support\Facades\Log;

class SendCodeController extends BaseController
{
    use CodeTrait;

    public function index()
    {
        $enterpriseId = $this->request->input('enterprise_id');
        $mobile       = $this->request->input('mobile');
        $client       = $this->request->input('client', 'console');
        $rule         = $this->request->input('module_rule', []);
        $fixed        = $this->request->input('fixed', 2);
        $fixedCode    = $this->request->input('fixed_code', 1234);
        $ip           = $this->request->getClientIp();

        if (empty($mobile)) {
            // 手机号不能为空
            notice(700031);
        }
        if (empty($client)) {
            // 客户端标识不能为空
            notice(700013);
        }
//        if(!$enterpriseId){
//            // 企业id不能为空
//            notice(700020);
//        }

        if ($client == 'console') {
            //如果手机号token有使用过
            if ($this->redis->get('limit:skip:' . $mobile)) {
                notice(700032);
            }

            //限制同一个手机号token短时间内重复登录
            $this->redis->setex('limit:skip:' . $mobile, 50, 1);
            if (!$trueAccount = $this->redis->get($mobile)) {
                notice(700032);
            }
            $mobile = $trueAccount;
        }
        //$mobile = ltrim($mobile, '+86-');
        $mobile = str_replace('+86-', '', $mobile);
        // 根据调用端传递的验证规则验证--自己写调用规则
        if (!empty($rule)) {
            (new SendCodeService())->sendCodeRule($rule, $mobile);
        }
        $whiteMobile = [
            '18092066590', '15110188985', '15249279779', '16621128726', '13588816087'
        ];
        if (in_array($mobile, $whiteMobile)) {
            return success();
        }

        if($enterpriseId != 1){
            $enterprise = (new EnterpriseAdmin())->accountDetailByMobile($mobile, $enterpriseId, ['short_message_config.enterprise_id']);

            if(!$enterprise || !isset($enterprise['enterprise_id'])){
                Log::info("企业[{$enterpriseId}]未配置短信");
                notice(700040);
            }

            $enterpriseId = $enterprise['enterprise_id'];
        }

        $this->sendCode($mobile, $client, $fixed, $fixedCode, $ip, '', $enterpriseId);

        return success();
    }
}
