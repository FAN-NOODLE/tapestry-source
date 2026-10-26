<?php

namespace App\Library\Uc;

class Admin extends Abstracts
{
    /**
     * 微信授权登录接口 传code值即可
     *
     * @var string
     */
    const DETAIL_URL = '%s/b/%s/adminer/member/detail';

    /**
     * 管理员详情
     *
     * @param $qyDomain
     * @param $params
     * @return mixed
     * @throws \Exception
     */
    public function detail($qyDomain, $params)
    {
        $url = vsprintf(self::DETAIL_URL, [$this->getUcApiUrl(), $qyDomain]);

        return $this->post($url, $params);
    }
}
