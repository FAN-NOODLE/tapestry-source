<?php
/**
 * @author  zhaopengfei<zhaopengfei@vchangyi.com>
 * @date    2020/11/16 14:59
 */

namespace App\Models;


class PartnerMember extends Base
{
    /**
     * 合伙人基本信息表
     * @var string
     */
    protected $table = 'partner_member';

    /**
     * 合伙人账号状态：0 禁用   1启用
     */
    const PARTNER_STATUS_ENABLED = 1;

    const PARTNER_STATUS_DISABLED = 0;
}