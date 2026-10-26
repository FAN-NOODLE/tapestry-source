<?php
/**
 * Created by PhpStorm.
 * Date: 2018/11/15
 * Time: 4:05 PM
 */

namespace App\Models;

/**
 * todo 废弃,app_secret迁移至enterprise_third_part_accout表
 * todo 商城小程序应用的agent_id和secret迁移至enterprise_work_plugin表
 * Class EnterpriseShopSetting
 * @package App\Models
 */
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
        'id',
        'enterprise_id',
        'mini_program_secret',
        'app_secret',
        'app_id',
        'corp_id',
        'pl_plugin_id',
        'agent_id',
        'ep_enumber',
        'created_at',
        'updated_at'
    ];
}
