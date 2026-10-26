<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/12/3
 * Time: 9:44 AM
 */

namespace App\Http\Controllers\App\V1\Member;

use App\Exceptions\Error;
use App\Http\Controllers\App\V1\BaseController;
use App\Library\Uc\Common;
use App\Models\MemberGuide;
use CrsSdk\Boss\Enterprise;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Crypt;
use Illuminate\Support\Facades\Log;
use App\Common\Constant;

class HubLoginController extends BaseController
{
    /**
     * 该接口非必须登录
     *
     * @var bool
     */
    protected $requireLogin = false;

    /**
     * 企业微信侧小程序-平台小程序(天鲟CRS)授权登录
     *
     * @param Request $request
     *
     * @return array
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $code = $request->post('code');

        // 商城小程序 = crs-shop
        // crs小程序 = crs-mini-program
        // 平台小程序(天鲟CRS) = crs-guide
        // $appIdentifier = $request->post('app_identifier', 'crs-guide');

        if (empty($code)) {
            throw new Error(110021, '参数错误');
        }

        $env        = env('APP_ENV', 'dev');
        $workPlugin = config('work_plugin.crs-guide');
        $workPlugin = $workPlugin[$env];
        $suitId     = $workPlugin['suite_id'];

        // 调用UC code2session 接口
        $result = (new Common())->hubMiniProgramLogin([
            'code'    => $code,
            'suiteid' => $suitId
        ]);

        if (!is_array($result)) {
            Log::error('110022 UC 返回错误' . $result, true);
            throw new Error(110022, '获取员工身份失败' . var_export($result, true));
        }

        // 获取企业corpId
        $corpId = $result['corpid'];

        if (empty($corpId)) {
            throw new Error(110023, '获取企业corpId失败');
        }

        // 获取该企业的授权信息
        $enterpriseData = (new Enterprise())->detail(['corp_id' => $corpId]);

        if (empty($enterpriseData['data'])) {
            throw new Error(110024, '获取企业信息有误');
        }

        $enterpriseData = $enterpriseData['data'];
        $crsStatus      = $enterpriseData['crs_status'];

        // 是否开启
        if ($crsStatus == Constant::CRS_CLOSE) {
            throw new Error(110025, 'CRS小程序未开启');
        }

        // 拼接token信息
        $token = [
            'memUid'     => $result['memUid'],
            'corpId'     => $result['corpid'],
            'sessionKey' => $result['session_key'],
            'qyDomain'   => $result['epEnumber'],
            'userId'     => $result['userid']
        ];

        $token = Crypt::encrypt(json_encode($token));

        $memberExtend = MemberGuide::query()->where('mem_uid', $result['memUid'])->first();

        $authority = MemberGuide::AUTHORITY_NEGATIVE;

        if (!empty($memberExtend->mem_uid)) {
            $authority = MemberGuide::AUTHORITY_POSITIVE;
        }

        return success(
            [
                'data' => [
                    'token'  => $token,
                    'member' => [
                        'mem_uid'   => $result['memUid'],
                        'authority' => $authority,
                    ]
                ]
            ]
        );
    }
}
