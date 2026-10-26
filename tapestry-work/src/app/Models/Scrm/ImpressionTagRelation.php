<?php

namespace App\Models\Scrm;

class ImpressionTagRelation extends ScrmInit
{

    /**
     * 表名
     * 外部联系人印象标签关系表
     *
     * @var string
     */
    protected $table = 'impression_tag_relation';


    /**
     * 主键
     *
     * @var string
     */
    public $primaryKey = 'id';

    /**
     * @var array
     */
    protected $guarded = [];
}