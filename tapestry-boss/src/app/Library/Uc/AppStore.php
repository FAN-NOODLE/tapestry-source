<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/10/12
 * Time: 10:33 AM
 */

namespace App\Library\Uc;


class AppStore extends Abstracts
{
    /**
     * 查询订单
     *
     * @var string
     */
    const APP_STORE_ORDER = '%s/a/%s/%s/appStore/order';



    /**
     * 获取应用列表接口
     *
     * @param string $qyDomain 企业域名
     * @param array $params 筛选条件
     * @return array|\Psr\Http\Message\StreamInterface
     * @throws \Exception
     */
    public function order($qyDomain, $params = [])
    {
        $appIdentifier = 'crs_store_assistant';
        $url = vsprintf(self::APP_STORE_ORDER, [$this->getUcApiUrl(), $qyDomain, $appIdentifier]);
        return $this->post($url, $params);
    }

}