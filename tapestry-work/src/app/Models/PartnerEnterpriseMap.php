<?php
/**
 * @author  zhaopengfei<zhaopengfei@vchangyi.com>
 * @date    2020/11/16 16:32
 */

namespace App\Models;


class PartnerEnterpriseMap extends Base
{
    /**
     * 合伙人和企业导购关联表
     * @var string
     */
    protected $table = 'partner_enterprise_map';

    /**
     * 是否给合伙人发送短信
     * 1 : 开启
     * 2 : 关闭
     * @var int
     */
    const sms_switch_true = 1;
    const sms_switch_false = 2;


}