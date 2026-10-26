<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/23
 * Time: 17:35
 */

namespace App\Models;

class CommissionSetting extends Base
{
    /**
     * 导购类型： 1：专属导购，2：普通导购，3：服务导购
     */
    const TYPE_EXCLUSIVE = 1;
    const TYPE_CONVENTION = 2;
    const TYPE_SERVICE = 3;

    /**
     * 表名
     *
     * @var string
     */
    protected $table = 'commission_setting';

    /**
     * 主键
     *
     * @var string
     */
    public $primaryKey = 'id';

    /**
     * @var array
     */
    protected $fillable = [
        'id',
        'enterprise_id',
        'type',
        'value'
    ];
}
