<?php
/*
 * This file is part of the car/chedianai_bc.
 *
 * (c) chedianai <i@chedianai.com>
 *
 * This source file is subject to the MIT license that is bundled
 * with this source code in the file LICENSE.
 */

namespace App\Services\OpenPlatform;

use App\Exceptions\Error;
use App\Models\Enterprise;
use App\Models\EnterpriseThirdPartAccount;
use App\Services\BaseService;
use CrsSdk\MicroAuthority\WX;
use EasyWeChat\OpenPlatform\Application;
use Illuminate\Support\Facades\Log;

class OpenPlatformService extends BaseService
{
    /**
     * @var Application
     */
    public $openPlatform;

    /**
     * @var Enterprise
     */
    protected $enterprise;

    const WX_FORMAL_TEST = ['wx570bc396a51b8ff8', 'wxd101a85aa106f53e'];

    public function __construct(Application $openPlatform, Enterprise $enterprise)
    {
        $this->openPlatform = $openPlatform;
        $this->enterprise = $enterprise;
    }

    /**
     * 判断该 appId 对应的是小程序还是服务号，并返回对应的授权信息.
     *
     * @param $appId
     * @return array
     */
    public function getAuthorizerConfigByAppId($appId, $enterpriseId)
    {
        // 切换企业开放平台配置
        changeOpenPlatformConfig($enterpriseId, $this->openPlatform);
        $auth = $this->openPlatform->getAuthorizer($appId);
        Log::info('auth-info：' . var_export($auth, true));

        $authInfo = $auth['authorizer_info'];
        $authConfig = $auth['authorization_info'];

        // 判断是小程序还是公众号(服务号或者订阅号)
        if (key_exists('MiniProgramInfo', $authInfo)) {
            $type = EnterpriseThirdPartAccount::MINI_PROGRAM;
        } else {
            if ($authInfo['service_type_info']['id'] == 2) {
                $type = EnterpriseThirdPartAccount::FU_WU_HAO;
            } else {
                $type = EnterpriseThirdPartAccount::DING_YUE_HAO;
            }
        }
        $authInfo['type'] = $type;

        return [
            'authorizer_info'    => $authInfo,
            'authorization_info' => $authConfig,
        ];
    }

    /**
     * 第三方平台全网发布时，应对微信平台的测试.
     * 返回bool，代表不处理，是用户真实请求
     *
     * @param $appId
     * @param $message
     * @return string|bool
     * @throws \EasyWeChat\Kernel\Exceptions\InvalidArgumentException
     * @throws \EasyWeChat\Kernel\Exceptions\RuntimeException
     */
    public function wechatWholeNetworkRelease($appId, $message, $enterpriseId)
    {
        changeOpenPlatformConfig($enterpriseId, $this->openPlatform);
        if (!in_array($appId, self::WX_FORMAL_TEST)) {
            return false;
        }

        if ('event' == $message['MsgType']) {
            return $message['Event'] . 'from_callback';
        } elseif ('text' == $message['MsgType'] && 'TESTCOMPONENT_MSG_TYPE_TEXT' == $message['Content']) {
            return $message['Content'] . '_callback';
        }

        if (str_contains($message['Content'], 'QUERY_AUTH_CODE')) {
            $authCodeArr = explode(':', $message['Content']);

            $authInfo = $this->openPlatform->handleAuthorize($authCodeArr[1]);
            $this->openPlatform->miniProgram($authInfo['authorization_info']['authorizer_appid'], $authInfo['authorization_info']['authorizer_refresh_token'])->customer_service->message($authCodeArr[1] . '_from_api')
                ->to($message['FromUserName'])
                ->send();
        }

        return '';
    }

    /**
     * 代开小程序、公众号的第三方的开放平台.
     *
     * @param $enterpriseId
     * @param $appId
     * @param $refreshToken
     * @param $type
     * @return mixed|null
     * @throws \EasyWeChat\Kernel\Exceptions\InvalidConfigException
     */
    public function setVirtualOpenPlatform($enterpriseId, $appId, $refreshToken, $type)
    {
        // 切换企业开放平台配置
        changeOpenPlatformConfig($enterpriseId, $this->openPlatform);
        // 查看该企业微信配置是否存在 open_appid
        $enterprise = $this->enterprise->with('virtualOpenPlatform')->find($enterpriseId);
        // 本系统open_app_id
        $openAppId = $enterprise->virtualOpenPlatform ? $enterprise->virtualOpenPlatform->open_app_id : null;

        // 微信平台open_app_id
        $openRes = $this->virtualOpenPlatformBinding($appId, $refreshToken, $type);

        if ($openAppId) { // 本系统存在open_app_id
            if ($openRes['errcode'] == 0) { // 已绑定过微信平台的open_app_id
                if ($openRes['open_appid'] != $openAppId) { // 和系统记录open_app_id不一致
                    // 解绑
                    $this->virtualOpenPlatformUnBind($openRes['open_appid'], $appId, $refreshToken, $type);
                    // 再次加入开放平台
                    $this->virtualOpenPlatformBindTo($openAppId, $appId, $refreshToken, $type);
                }
            } else { // 未绑定微信平台的open_appid
                // 绑定
                $this->virtualOpenPlatformBindTo($openAppId, $appId, $refreshToken, $type);
            }
        } else { // 否则判断是否要创建第三方平台
            if ($openRes['errcode'] != 0) { // 创建第三方平台
                $res = $this->virtualOpenPlatformCreate($appId, $refreshToken, $type);
                $openAppId = $res['open_appid'];
            } else { // 不用创建
                $openAppId = $openRes['open_appid'];
            }
        }

        return $openAppId;
    }

    /**
     * 创建虚拟开放平台.
     *
     * @param $appId
     * @param $refreshToken
     * @param $type
     * @return mixed|null
     * @throws \EasyWeChat\Kernel\Exceptions\InvalidConfigException
     */
    protected function virtualOpenPlatformCreate($appId, $refreshToken, $type)
    {
        $res = null;

        switch ($type) {
            case EnterpriseThirdPartAccount::MINI_PROGRAM:
                $res = $this->openPlatform->miniProgram($appId, $refreshToken)->account->create();
                break;
            case EnterpriseThirdPartAccount::DING_YUE_HAO:
            case EnterpriseThirdPartAccount::FU_WU_HAO:
                $res = $this->openPlatform->officialAccount($appId, $refreshToken)->account->create();
                break;
            default:
                throw new Error(100000, '不存在该种类型的微信公众号、小程序');
        }

        if ($res['errcode'] != 0) {
            throw new Error(100000, '创建开放平台失败，微信返回: ' . json_encode($res));
        }

        return $res;
    }

    /**
     * 绑定虚拟开放平台.
     *
     * @param $openAppId
     * @param $appId
     * @param $refreshToken
     * @param $type
     * @return mixed|null
     * @throws Error
     * @throws \EasyWeChat\Kernel\Exceptions\InvalidConfigException
     */
    protected function virtualOpenPlatformBindTo($openAppId, $appId, $refreshToken, $type)
    {
        $res = null;
        switch ($type) {
            case EnterpriseThirdPartAccount::MINI_PROGRAM:
                $res = $this->openPlatform->miniProgram($appId, $refreshToken)->account->bindTo($openAppId);
                break;
            case EnterpriseThirdPartAccount::DING_YUE_HAO:
            case EnterpriseThirdPartAccount::FU_WU_HAO:
                $res = $this->openPlatform->officialAccount($appId, $refreshToken)->account->bindTo($openAppId);
                break;
            default:
                throw new Error(100000, '不存在该种类型的微信公众号、小程序');
        }

        if ($res['errcode'] != 0) {
            throw new Error(100000, '绑定第三方公众号失败, 绑定的open_app_id: ' . $openAppId . ', 微信返回：' . json_encode($res));
        }

        return $res;
    }

    /**
     * 解绑虚拟开放平台.
     *
     * @param $openAppId
     * @param $appId
     * @param $refreshToken
     * @param $type
     * @return mixed|null
     * @throws Error
     * @throws \EasyWeChat\Kernel\Exceptions\InvalidConfigException
     */
    protected function virtualOpenPlatformUnBind($openAppId, $appId, $refreshToken, $type)
    {
        $res = null;
        switch ($type) {
            case EnterpriseThirdPartAccount::MINI_PROGRAM:
                $res = $this->openPlatform->miniProgram($appId, $refreshToken)->account->unbindFrom($openAppId);
                break;
            case EnterpriseThirdPartAccount::DING_YUE_HAO:
            case EnterpriseThirdPartAccount::FU_WU_HAO:
                $res = $this->openPlatform->officialAccount($appId, $refreshToken)->account->unbindFrom($openAppId);
                break;
            default:
                throw new Error(100000, '不存在该种类型的微信公众号、小程序');
        }

        if ($res['errcode'] != 0) {
            throw new Error(100000, '解绑第三方公众号失败, 解绑的open_app_id: ' . $openAppId . ',  微信返回：' . json_encode($res));
        }

        return $res;
    }

    /**
     * 查询虚拟开放平台.
     *
     * @param $appId
     * @param $refreshToken
     * @param $type
     * @return mixed|null
     * @throws Error
     * @throws \EasyWeChat\Kernel\Exceptions\InvalidConfigException
     */
    protected function virtualOpenPlatformBinding($appId, $refreshToken, $type)
    {
        $res = null;
        switch ($type) {
            case EnterpriseThirdPartAccount::MINI_PROGRAM:
                $res = $this->openPlatform->miniProgram($appId, $refreshToken)->account->getBinding();
                break;
            case EnterpriseThirdPartAccount::DING_YUE_HAO:
            case EnterpriseThirdPartAccount::FU_WU_HAO:
                $res = $this->openPlatform->officialAccount($appId, $refreshToken)->account->getBinding();
                break;
            default:
                throw new Error(100000, '不存在该种类型的微信公众号、小程序');
        }

        return $res;
    }

    /**
     * @param $appId
     * @return array
     * @throws \Exception
     */
    public function getAuthorizeTokenByAppId($appId)
    {
        $refreshToken = EnterpriseThirdPartAccount::query()
            ->select('enterprise_id')
            ->where('app_id', '=', $appId)
            ->where('type', '=', EnterpriseThirdPartAccount::TYPE_SHOP)
            ->where('status', '=', EnterpriseThirdPartAccount::STATUS_ACTIVE)
            ->where('platform', '=', EnterpriseThirdPartAccount::MINI_PROGRAM_ACCOUNT)
            ->first();

        if (empty($refreshToken)) {
            throw new Error(100000, '小程序refresh_token不存在');
        }

        $token = (new WX())->accessToken(['enterprise_id' => $refreshToken->enterprise_id]);

        $auth['authorizer_access_token'] = $token['access_token'];
        return $auth;
    }
}
