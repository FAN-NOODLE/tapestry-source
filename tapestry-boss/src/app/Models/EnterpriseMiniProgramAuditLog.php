<?php
/**
 * Created by PhpStorm.
 * Date: 2018/11/15
 * Time: 4:05 PM
 */

namespace App\Models;

use App\Models\ConstTrait\MiniProgram;

class EnterpriseMiniProgramAuditLog extends Base
{
    /**
     * 表名称
     *
     * @var string
     */
    protected $table = 'enterprise_mini_program_audit_log';

    /**
     * 引入公共常量
     *
     */
    use MiniProgram;

    /**
     * 黑名单
     *
     * @var array
     */
    protected $guarded = [];
}
