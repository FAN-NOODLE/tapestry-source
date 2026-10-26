<?php

namespace App\Sdk\Interfaces;

interface CdpInterface
{

    /**
     * 获取第三方会员信息
     *
     * @param $user array 用户信息
     *
     * @return mixed
     */
    public function getMemberInfo(array $user);



}
