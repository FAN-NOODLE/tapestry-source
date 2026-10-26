<?php

namespace App\Sdk\Active\CRM;

use App\Sdk\Interfaces\CdpInterface;

abstract class BaseCrm implements CdpInterface
{

    /*
    |--------------------------------------------------------------------------
    | CrmInterface方法在基类实现了，如SDK不需要则不需要重写
    |
    | 如需具体SDK实现，可重写对应方法，返回值格式参考基类
    |
    |--------------------------------------------------------------------------
    */

    public function reqResponse(array $user): array
    {
        return [
            'status'   => false,
            'union_id' => $user['union_id'] ?? '',        //String	微信unionId
            'mobile'   => $user['mobile'] ?? '',          //String	会员手机号
            'msg'      => '请求接口返回:未获取到会员信息'
        ];
    }

    public function reqMemberInfoNo(array $user): array
    {
        return [
            'status'      => false,
            'mobile'      => $user['mobile'] ?? '',          //String	会员手机号
            'open_id'     => $user['open_id'] ?? '',         //String	微信openId
            'union_id'    => $user['union_id'] ?? '',        //String	微信unionId
//            "pageNum"     => 0,
//            "pageSize"    => 10,
//            "isFirstPage" => true,
//            "isLastPage"  => true,
//            "hasNextPage" => false,
//            "total"       => 0,
//            "pages"       => 0,
            'msg'         => '请求接口返回:未获取到会员信息'
        ];
    }

}
