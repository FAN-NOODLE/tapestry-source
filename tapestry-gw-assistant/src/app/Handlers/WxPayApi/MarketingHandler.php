<?php
/**
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/5/24
 * Time: 17:51
 */

namespace App\Handlers\WxPayApi;

class MarketingHandler extends WxPayClient
{
    /**
     * 根据核销商户查用户的券API
     * %s = openid
     */
    const WX_PAYMENT_MARKET_USERS_COUPONS = '/marketing/favor/users/%s/coupons';

    /**
     * 根据批次号下载批次核销明细API
     */
    const WX_MARKETING_FAVOR_STOCKS = '/marketing/favor/stocks/%s/use-flow';

    /**
     * 根据批次号下载批次核销明细API
     */
    const WX_MARKETING_FAVOR_CALLBACKS = '/marketing/favor/callbacks';

    /**
     * 根据批次号查询批次详情API
     */
    const WX_MARKETING_FAVOR_STOCK_INFO = '/marketing/favor/stocks/%s';

    /**
     * 条件查询批次列表API
     */
    const WX_MARKETING_FAVOR_STOCK_LIST = '/marketing/favor/stocks';

    /**
     * 根据用户open_id和优惠券id查询优惠券详情
     */
    const WX_MARKETING_USER_COUPON_DETAIL = '/marketing/favor/users/%s/coupons/%s';

    /**
     * 查询用户优惠券具体批次信息
     * @param array $data
     * @return mixed|\Psr\Http\Message\StreamInterface|string
     * @throws \GuzzleHttp\Exception\GuzzleException
     */
    public function usersCouponDetail($data = [])
    {
        $url = sprintf(self::WX_MARKETING_USER_COUPON_DETAIL, $data['open_id'], $data['coupon_id']);
        $response = $this->wxHttpGet($url, $data);
        return $response;
    }

    /**
     * 根据核销商户查用户的券
     *
     * @param array $data
     *          int $data[enterprise_id]                企业id
     *          int $data[open_id]                      openid
     *          int $data[param]                        请求参数
     *          string $data[params][appid]             商户appId
     *          string $data[params][available_mchid]   商户号
     *
     * @return mixed|\Psr\Http\Message\StreamInterface|string
     * @throws \GuzzleHttp\Exception\GuzzleException
     */
    public function usersCoupons($data = [])
    {
        $url      = sprintf(self::WX_PAYMENT_MARKET_USERS_COUPONS, $data['open_id']);
        $response = $this->wxHttpGet($url, $data);
        return $response;
    }

    /**
     * 根据批次号下载批次核销明细
     *
     * @param array $data
     *          int $data[enterprise_id]                企业id
     *          int $data[stock_id]                     批次号
     *          int $data[param]                        请求参数
     *          string $data[params][appid]             商户appId
     *          string $data[params][available_mchid]   商户号
     *
     * @return mixed|\Psr\Http\Message\StreamInterface|string
     * @throws \GuzzleHttp\Exception\GuzzleException
     */
    public function favorStocks($data = [])
    {
        $url      = sprintf(self::WX_MARKETING_FAVOR_STOCKS, $data['stock_id']);
        $response = $this->wxHttpGet($url, $data);
        return $response;
    }

    /**
     * 根据批次号查询批次详情API
     *
     * @param array $data
     *          int $data[enterprise_id]                企业id
     *          int $data[stock_id]                     批次号
     *          int $data[param]                        请求参数
     *          string $data[params][appid]             商户appId
     *          string $data[params][available_mchid]   商户号
     *
     * @return mixed|\Psr\Http\Message\StreamInterface|string
     * @throws \GuzzleHttp\Exception\GuzzleException
     */
    public function favorStockInfo($data = [])
    {
        $url      = sprintf(self::WX_MARKETING_FAVOR_STOCK_INFO, $data['stock_id']);
        $response = $this->wxHttpGet($url, $data);
        return $response;
    }

    /**
     * 根据批次号下载批次核销明细
     *
     * @param       $apiUrl
     * @param       $method
     * @param array $data
     *
     * @return mixed
     */
    public function getSign($apiUrl, $method, $data = [])
    {
        $authorization = $this->getAuthorization($apiUrl, $method, [], $data['paymentConfig']);
        return $authorization;
    }


    /**
     * 设置消息回调通知API
     *
     * @param array $data
     *          int     $data[enterprise_id]                企业id
     *          array   $data[paymentConfig]                支付配置参数
     *
     * @return mixed|\Psr\Http\Message\StreamInterface|string
     * @throws \GuzzleHttp\Exception\GuzzleException
     */
    public function setFavorCallbacksUrl($data = [])
    {
        $url      = self::WX_MARKETING_FAVOR_CALLBACKS;
        $response = $this->wxHttpPost($url, $data);
        return $response;
    }

    /**
     * 条件查询批次列表API
     * @param $data
     * @return mixed|\Psr\Http\Message\StreamInterface|string
     * @throws \GuzzleHttp\Exception\GuzzleException
     */
    public function getStockList($data)
    {
        $url = self::WX_MARKETING_FAVOR_STOCK_LIST;
        return $this->wxHttpGet($url, $data);
    }
}
