<?php
/**
 * Created by PhpStorm.
 * User: qihongchuan
 * Date: 2020/02/12
 * Time: 14:35 PM
 */

namespace App\Models;

class EnterpriseSdkConfig extends Base
{
    /**
     * 小程序类型：1.商城小程序，2.企业微信小程序
     */
    const TYPE_SHOP  = 1;
    const TYPE_QY_WX = 2;

    /**
     * 表名称
     *
     * @var string
     */
    protected $table = 'enterprise_sdk_config';

    /**
     * 可编辑字段
     *
     * @var string
     */
    protected $fillable = [
        'id',
        'enterprise_id',
        'app_type',
        'sdk_name',
        'request_domain',
        'pack_config',
        'ext_config'
    ];

}
