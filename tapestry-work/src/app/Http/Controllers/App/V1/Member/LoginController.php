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
use App\Services\EnterpriseService;
use CrsSdk\Shop\StaffRight;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Crypt;
use Illuminate\Support\Facades\Log;
use App\Common\Constant;

class LoginController extends BaseController
{
    /**
     * 该接口非必须登录
     *
     * @var bool
     */
    protected $requireLogin = false;

    /**
     * 企业微信侧小程序授权登录
     *
     * @param Request $request
     *
     * @return array
     * @throws \Exception
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/12/3 9:46 AM
     */
    public function index(Request $request)
    {
        $code   = $request->post('code');
        $userId = $request->post('user_id');

        // 商城小程序 = crs-shop
        // crs小程序 = crs-mini-program
        $appIdentifier = $request->post('app_identifier', 'crs-mini-program');

        if (empty($code)) {
            throw new Error(110001, '参数错误');
        }

        $enterpriseHash = $request->header('Enterprise-Hash');

        if (empty($enterpriseHash)) {
            throw new Error(110002, '企业HASH不能为空');
        }

        // 获取企业详情
        $enterpriseData= (new EnterpriseService())->getHashInfo($enterpriseHash);

        $qyDomain = $enterpriseData['data']['domain'] ?? '';
        if (!$qyDomain) {
            throw new Error(110006, '请检查企业和员工信息是否正确');
        }

        $crsStatus    = $enterpriseData['data']['crs_status'];
        // 是否开启
        if ($crsStatus == Constant::CRS_CLOSE) {
            notice(110005);
        }

        // 调用UC code2session 接口
        $result = (new Common())->miniProgramLogin($qyDomain, $appIdentifier, ['code' => $code]);

        if (!is_array($result)) {
            Log::error('110003 UC 返回错误' . $result);
            throw new Error(110003, '获取员工身份失败' . var_export($result, true));
        }

        // 员工绑定用户
        if ($userId) {
            (new StaffRight())->staffBindUser([
                'user_id' => $userId,
                'mem_uid' => $result['memUid']
            ]);
        }

        // 拼接token信息
        $token = [
            'memUid'     => $result['memUid'],
            'corpId'     => $result['corpid'],
            'sessionKey' => $result['session_key'],
            'qyDomain'   => $result['epid'],
            'userId'     => $result['userid'],
            'crs_status' => $crsStatus
        ];

        $token = Crypt::encrypt(json_encode($token));

        return success(
            [
                'data' => [
                    'token'  => $token,
                    'member' => [
                        'mem_uid'      => $result['memUid'],
                    ]
                ]
            ]
        );
    }
}
