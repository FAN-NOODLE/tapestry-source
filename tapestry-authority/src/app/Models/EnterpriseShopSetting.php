<?php
/**
 * Created by PhpStorm.
 * Date: 2018/11/15
 * Time: 4:05 PM
 */

namespace App\Models;

class EnterpriseShopSetting extends Base
{
    /**
     * 表名称
     *
     * @var string
     */
    protected $table = 'enterprise_shop_setting';

    /**
     * 允许写入的字段
     *
     * @var array
     */
    protected $fillable = [
        'enterprise_id',
        'mini_program_secret',
        'app_secret',
        'app_id',
        'corp_id',
        'pl_plugin_id',
        'ep_enumber',
        'created_at',
        'updated_at'
    ];
}
