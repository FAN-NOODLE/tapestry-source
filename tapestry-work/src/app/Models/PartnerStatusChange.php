<?php


namespace App\Models;


class PartnerStatusChange extends Base
{
    /**
     * 合伙人状态变更流水表
     * @var string
     */

    //禁用原因编码
    const FORBIDDEN_CODE = [
        'custom'  => 0,    //自定义原因
        'negative' => 1    //余额为负自动禁用
    ];

    protected $table = 'partner_status_change';

}