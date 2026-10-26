<?php

namespace App\Models;

class ReleaseEnterpriseWork extends Base
{
    /**
     * 表名称
     *
     * @var string
     */
    protected $table = 'release_enterprise_work';


    /**
     * 可编辑字段
     *
     * @var string
     */
    protected $fillable = [
        'id',
        'enterprise_id',
        'identifier'
    ];
}
