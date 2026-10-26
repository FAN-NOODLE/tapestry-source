<?php


namespace App\Http\Controllers\App\V1\Common;


use App\Common\Constant;
use App\Http\Controllers\App\V1\BaseController as Controller;
use App\Traits\CodeTrait;
use CrsSdk\MicroAuthority\Common;
use CrsSdk\Work\PartnerMember;
use CrsSdk\Boss\Enterprise;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class SendSmsController extends Controller
{
    protected $isRequireLogin = false;

    use CodeTrait;

    const MODULE_LOGIN = 'LOGIN';
    const MODULE_CHECK_SMS_MOBILE = 'CHECK_SMS_MOBILE';
    const PARTNER_MEMBER_REGISTER = 'PARTNER_REGISTER';//合伙人注册
    const PARTNER_MEMBER_REGISTER_TLL = 60;//合伙人注册验证码有效期
    //deng
    /**
     * 模拟的手机号
     *
     * @var array
     */
    private $mockMobile = [
        '13227766759',
        // 魏梦洁
        '13720452952',
    ];
    /**
     * 模拟的手机号 - 测试环境用
     * @desc 判断在非正式环境将两个数组进行了合并，添加测试手机号时添加到下面的数组中
     * @var array
     */
    private $mockMobileTest = [
        '18229054347',
        '15829001821',
        '13022977050',
        // hxm
        '17092933378',
        // hxm
        '18700465991',
        // 罗兰
        '18220843860',
        // 张海印
        '18372683027',
        // 陈耀宏
        '18706740335',
        // 高正城
        '18291425164',
        // 张哪
        '13038590099',
        '18792598960',
        // 景凯
        '13109522355',
        // 龙阳
        '18161732940',
        '13259961589',
        // 李强
        '18291875030',
        '15021955295',
        '13991003353',
        '15099922630',
        '17795851356',
        '18202940551',
        '15129184497',
        '17691046130',
        '18209229895',
        '18629365475',
        '18668019363',
        '18392949820',
        '17691246013',
        '18392415020',
        '15829239163',
        '17691082097',
        '15829239163',
        '18192503204',
        '13811541196',
    ];

    /**
     * 发送短信验证码
     *
     * @param Request $request
     * @param Common $common
     * @return \Illuminate\Http\JsonResponse|void
     */
    public function __invoke(Request $request, Common $common)
    {
        $env = config('app.env');
        if ($env != 'prod') {
            $this->mockMobile += $this->mockMobileTest;
        }
        
        $mobile = trim($request->get('mobile', ''));
//        $corpId = trim($request->input('corp_id'));
        $module = $request->get('module', self::MODULE_LOGIN);
        if (empty($mobile)) {
            error(651000);
        }

//        if(!$corpId){
//            return error(651008);
//        }

        if (in_array($mobile, Constant::WX_LOGIN_SUPER_USER)) {
            return success();
        }
        try {
            $sendCodeParams = [
                'client' => Constant::APP_CLIENT_NAME,
                'mobile' => $mobile
            ];
            // 如果为模拟账号，可以发固定的验证码
            if (in_array($mobile, $this->mockMobile)) {
                $sendCodeParams['fixed'] = 1;
                $sendCodeParams['fixed_code'] = '1234';
            }
            // 登录发送验证码
            if ($module === self::MODULE_LOGIN) {
//                try{
//                    $enterprise = (new Enterprise())->detail([
//                        'corp_id' => $corpId
//                    ]);
//
//                    $sendCodeParams['enterprise_id'] = isset($enterprise['data']['id']) ? $enterprise['data']['id'] : '';
//                } catch (\Exception $e){
//                    Log::info('获取企业相关信息失败', [
//                        'corp_id' => $corpId
//                    ]);
//                }

                $common->SendCode($sendCodeParams);
            }
            // 更变手机号发送验证码
            if ($module == self::MODULE_CHECK_SMS_MOBILE) {
                $fixed = isset($sendCodeParams['fixed']) ?? null;
                $fixed_code = isset($sendCodeParams['fixed_code']) ?? null;
                $this->sendCode($mobile, $sendCodeParams['client'], $fixed, $fixed_code);
            }
            //合伙人注册验证码
            if ($module === self::PARTNER_MEMBER_REGISTER) {
                $fixed = isset($sendCodeParams['fixed']) ?? null;
                $fixed_code = isset($sendCodeParams['fixed_code']) ?? null;
                $this->SendCode($mobile,$sendCodeParams['client'],$fixed,$fixed_code,
                    self::PARTNER_MEMBER_REGISTER_TLL,'partner_register');
            }

        } catch (\Exception $exception) {
            Log::error("给{$mobile}发送短信验证码失败,错误信息:{$exception->getMessage()}");
            error($exception->getCode(), $exception->getMessage());
        }
        return success();
    }
}
