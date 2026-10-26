<?php
/*
 * This file is part of the car/chedianai_bc.
 *
 * (c) chedianai <i@chedianai.com>
 *
 * This source file is subject to the MIT license that is bundled
 * with this source code in the file LICENSE.
 */

namespace App\Modules\ThirdPart;

use App\Exceptions\Error;
use App\Models\EnterpriseThirdPartAccount;
use App\Modules\BaseModule;
use EasyWeChat\Factory;
use Overtrue\LaravelWeChat\Facade as EasyWeChat;

class OpenPlatformModule extends BaseModule
{
    /**
     * @var null | \EasyWeChat\OpenPlatform\Authorizer\OfficialAccount\Application |
     *      \EasyWeChat\OfficialAccount\Application
     */
    public $official_account = null;

    /**
     * @var null | \EasyWeChat\OpenPlatform\Authorizer\MiniProgram\Application | \EasyWeChat\MiniProgram\Application
     */
    public $mini_program = null;

    /**
     * 初始化微信方法.
     *
     * @param string | int $enterpriseId
     * @param int          $type
     *
     * @return $this
     *
     * @author dj
     *
     * @throws Error
     */
    public function init($enterpriseId, $type = 0)
    {
        // 获取激活企业账号
        $params = ['status' => EnterpriseThirdPartAccount::STATUS_ACTIVE];
        if ($type) {
            $params['type'] = $type;
        }

        $params['enterprise_id'] = $enterpriseId;

        $configs = app(EnterpriseThirdPartAccount::class)->where($params)->get();
        if (!count($configs)) {
            throw new Error(100000, '未授权微信公众号');
        }

        foreach ($configs as $config) {
            if (EnterpriseThirdPartAccount::AUTHORIZE_TYPE_CUSTOMER == $config->authorize_type) { // 手动配置
                $this->initManual(
                    $config->platform,
                    $config->app_id,
                    $config->app_secret,
                    $config->server_token,
                    $config->encoding_aes_key
                );
            } elseif (EnterpriseThirdPartAccount::AUTHORIZE_TYPE_THIRD == $config->authorize_type) { // 第三方开放平台
                $this->initComponent($config->platform, $config->app_id, $config->refresh_token, $enterpriseId);
            }
        }

        return $this;
    }

    /**
     * 初始化手动配置.
     *
     * @param string $platform 平台
     * @param string $appId
     * @param string $appSecret
     * @param string $token
     * @param string $aesKey
     *
     * @author dj
     *
     * @throws Error
     */
    protected function initManual($platform, $appId, $appSecret, $token, $aesKey = '')
    {
        $config = config('wechat.defaults');

        if (!$appId || !$appSecret || !$token) {
            throw new Error(100000, '微信缺少初始化参数');
        }

        $config['app_id'] = $appId;
        $config['secret'] = $appSecret;
        $config['token'] = $token;
        $config['aes_key'] = $aesKey ? : '';

        if (EnterpriseThirdPartAccount::OFFICIAL_ACCOUNT == $platform) {
            $this->official_account = Factory::officialAccount($config);
        } elseif (EnterpriseThirdPartAccount::MINI_PROGRAM_ACCOUNT == $platform) {
            $this->mini_program = Factory::miniProgram($config);
        } else {
            throw new Error(100000, '不存在该类型账号');
        }
    }

    /**
     * 初始化第三方平台设置.
     *
     * @param string $platform
     * @param string $appId
     * @param string $refreshToken
     *
     * @author dj
     *
     * @throws Error
     */
    protected function initComponent($platform, $appId, $refreshToken, $enterpriseId)
    {
        $openPlatform = EasyWeChat::openPlatform();
        changeOpenPlatformConfig($enterpriseId, $openPlatform);
        if (!$appId || !$refreshToken) {
            throw new Error(100000, '微信缺少初始化参数');
        }

        if (EnterpriseThirdPartAccount::OFFICIAL_ACCOUNT == $platform) {
            $this->official_account = $openPlatform->officialAccount($appId, $refreshToken);
        } elseif (EnterpriseThirdPartAccount::MINI_PROGRAM_ACCOUNT == $platform) {
            $this->mini_program = $openPlatform->miniProgram($appId, $refreshToken);
        } else {
            throw new Error(100000, '不存在该类型账号');
        }
    }
}
