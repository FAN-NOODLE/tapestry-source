<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/11/26
 * Time: 16:24
 */

namespace App\Models\ScrmEnterprise;


class MemberExternalBehaviorData extends ScrmEnterpriseInit
{


    /**
     * 主键
     *
     * @var string
     */
    public $primaryKey = 'id';
    /**
     * 表名
     *
     * @var string
     */
    protected $table = 'member_external_behavior_data';
    /**
     * @var array
     */
    protected $guarded = [];




}

