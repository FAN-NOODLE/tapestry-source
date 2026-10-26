<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/10/13
 * Time: 2:40 PM
 */

namespace App\Library\Uc;

class Common extends Abstracts
{
    /**
     * 微信授权登录接口 传code值即可
     *
     * @var string
     */
    const AUTH_URL = '%s/a/%s/%s/qy/commondapi/wxuserlogin';

    /**
     * 小程序登录
     *
     * @var string
     */
    const MINI_PRO_URL = '%s/a/%s/%s/qy/miniprogram/jscode2session';

    /**
     * 平台小程序(天鲟CRS)登录
     *
     * @var string
     */
    const HUB_MINI_PRO_URL = '%s/s/qy/suite-miniprogram/jscode2session';

    /**
     * 企业微信小程序消息发送
     *
     * @var string
     */
    const SEND_MESSAGE_URL = '%s/a/%s/crs-mini-program/qy/send_mini';

    /**
     * 根据auth_code获取授权信息
     */
    const WX_CODE_INFO_URL = '%s/s/qy/login/wx-qrcode-login';

    /**
     * @param       $qyDomain
     * @param       $appIdentifier
     * @param array $params
     *
     * @return mixed
     * @throws \App\Exceptions\Error
     * @throws \Exception
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/10/14 3:33 PM
     */
    public function authLogin($qyDomain, $appIdentifier, $params = [])
    {
        $url = vsprintf(self::AUTH_URL, [$this->getUcApiUrl(), $qyDomain, $appIdentifier]);
        return $this->post($url, $params);
    }

    /**
     * 小程序登录
     *
     * @param        $qyDomain
     * @param string $appIdentifier 应用标识
     * @param array  $params
     *
     * @return mixed
     * @throws \Exception
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/12/3 9:43 AM
     */
    public function miniProgramLogin($qyDomain, $appIdentifier, $params = [])
    {
        $url = vsprintf(self::MINI_PRO_URL, [$this->getUcApiUrl(), $qyDomain, $appIdentifier]);
        return $this->post($url, $params);
    }

    /**
     * 平台小程序(天鲟CRS)登录
     *
     * @param array  $params
     *
     * @return mixed
     * @throws \Exception
     */
    public function hubMiniProgramLogin($params = [])
    {
        $url = vsprintf(self::HUB_MINI_PRO_URL, [$this->getUcApiUrl()]);
        return $this->post($url, $params);
    }

    /**
     * 企业微信侧小程序消息发送
     *
     * @param $qyDomain
     * @param $params
     *
     * @return mixed
     * @throws \Exception
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/12/6 12:07 PM
     */
    public function miniProgramSendMessage($qyDomain, $params)
    {
        $url = vsprintf(self::SEND_MESSAGE_URL, [$this->getUcApiUrl(), $qyDomain]);
        return $this->post($url, $params);
    }

    /**
     * auth_code换取登录用户的令牌接口
     *
     * @param $authCode
     * @return mixed
     * @throws \Exception
     */
    public function authCodeLogin($authCode)
    {
        $url = vsprintf(self::WX_CODE_INFO_URL, [$this->getUcApiUrl()]);

        return $this->post($url, [
            'wxAuthCode' => $authCode
        ]);
    }
}
