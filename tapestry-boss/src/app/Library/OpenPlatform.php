<?php
/**
 * Created by Phpstorm.
 * User : zoulongbo
 * Date : 2019-10-31
 * Tome : 10:10
 */

namespace App\Library;

class OpenPlatform extends OpenPlatformClient
{
    /**
     * 微信开放平台：查询服务商的当月提审限额
     *  access_token = 任意授权该第三方的小程序接口调用令牌
     */
    const WX_OPEN_SERVICE_PROVIDER_AUDIT_QUOTA_QUERY = '/queryquota';

    /**
     * 微信开放平台：加急审核申请
     * 小程序接口调用令牌
     */
    const WX_OPEN_SERVICE_PROVIDER_SPEED_UP_AUDIT = '/speedupaudit';

    /**
     * 微信开放平台：直播间列表
     */
    const WX_OPEN_BUSINESS_GET_LIVE_INFO  = '/business/getliveinfo';


    /**
     * 查询服务商的当月提审限额 次数
     * @param $params
     *
     * @return mixed|\Psr\Http\Message\StreamInterface
     * @throws \GuzzleHttp\Exception\GuzzleException
     */
    public function auditQuota($params)
    {
        $response = $this->httpGet(self::WX_OPEN_SERVICE_PROVIDER_AUDIT_QUOTA_QUERY, $params);
        return $response;
    }

    /**
     * 加急审核申请
     * @param $params
     *
     * @return mixed|\Psr\Http\Message\StreamInterface|string
     * @throws \GuzzleHttp\Exception\GuzzleException
     */
    public function speedUpAudit($params)
    {
        $response = $this->httpPost(self::WX_OPEN_SERVICE_PROVIDER_SPEED_UP_AUDIT, $params);
        return $response;
    }

    /**
     * 直播间列表
     * @param $params
     *
     * @return mixed|\Psr\Http\Message\StreamInterface|string
     * @throws \GuzzleHttp\Exception\GuzzleException
     */
    public function getLiveInfo($params)
    {
        $response = $this->httpPost(self::WX_OPEN_BUSINESS_GET_LIVE_INFO, $params);
        return $response;
    }
}
