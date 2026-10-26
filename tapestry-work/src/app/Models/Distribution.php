<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/11/26
 * Time: 16:24
 */

namespace App\Models;

class Distribution extends Base
{
    protected $table = 'distribution';

    protected $primaryKey = 'id';

    /**
     * 更新的字段
     *
     * @var string[]
     */
    protected $fillable = [
        'enterprise_id',
        'distribution_text',
        'distribution_title',
        'settlement_standard',
        'poster_img',
        'id'
    ];

    /**
     * 结算状态
     *
     * @var integer
     */
    const SETTLEMENT_STANDARD_COMPLETE = 100;
}
