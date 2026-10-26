<?php
/**
 * Created by PhpStorm.
 * Date: 2018/11/15
 * Time: 4:05 PM
 */

namespace App\Models;

class EnterpriseThirdPartTester extends Base
{
    /**
     * 表名称
     *
     * @var string
     */
    protected $table = 'enterprise_third_part_tester';

    /**
     * @var array
     */
    protected $fillable = [
        'id',
        'account_id',
        'wx_account',
        'name',
        'remark'
    ];
}
