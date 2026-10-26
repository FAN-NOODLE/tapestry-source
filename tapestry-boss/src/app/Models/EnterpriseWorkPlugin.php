<?php

namespace App\Models;

class EnterpriseWorkPlugin extends Base
{
    /**
     * 表名称
     *
     * @var string
     */
    protected $table = 'enterprise_work_plugin';

    /**
     * 自建导购
     */
    const CRS_ASSISTANT = 'crs_store_assistant_self';

    /**
     * 应用商店安装导购
     */
    const CRS_APP_ASSISTANT = 'crs_store_assistant';

    const I_CLICK_CRS_APP_ASSISTANT = 'iclk_crs_store_assistant';

    /**
     * 可编辑字段
     *
     * @var string
     */
    protected $fillable = [
        'id',
        'enterprise_id',
        'corp_id',
        'agent_id',
        'secret',
        'identifier',
        'check_file',
        'file_url',
        'pl_plugin_id',
        'mini_app_id'
    ];

    protected $casts = [
        'agent_id' => 'integer',
    ];
}
