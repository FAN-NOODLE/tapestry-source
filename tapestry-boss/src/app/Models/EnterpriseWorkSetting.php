<?php
/**
 * Created by PhpStorm.
 * Date: 2018/11/15
 * Time: 4:05 PM
 */

namespace App\Models;

/**
 * todo 废弃,外部联系人secret迁移至enterprise表
 * Class EnterpriseWorkSetting
 * @package App\Models
 */
class EnterpriseWorkSetting extends Base
{
    /**
     * 表名称
     *
     * @var string
     */
    protected $table = 'enterprise_work_setting';

    /**
     * @var array
     */
    protected $fillable = [
        'id',
        'enterprise_id',
        'external_secret',
        'mini_program_secret',
        'corp_id',
        'pl_plugin_id',
        'ep_enumber'
    ];
}
