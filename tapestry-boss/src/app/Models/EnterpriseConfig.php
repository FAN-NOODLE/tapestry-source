<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/11/15
 * Time: 4:05 PM
 */

namespace App\Models;

class EnterpriseConfig extends Base
{
    /**
     * 单仓
     */
    const SINGLE_WAREHOUSE = 1;

    /**
     * 多仓
     */
    const MORE_WAREHOUSE = 2;

    // ------------------------- type -------------------------

    // 打印小票类型
    const PRINT_QRCODE = 13;
    const PRINT_QRCODE_KEY = 'print_qrcode';

    /**
     * 表名称
     *
     * @var string
     */
    protected $table = 'enterprise_config';

    /**
     * 可编辑字段
     *
     * @var string
     */
    protected $fillable = [
        'id',
        'enterprise_id',
        'type',
        'value',
        'key',
        'desc',
        'name'
    ];
}
