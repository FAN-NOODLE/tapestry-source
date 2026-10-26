<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2019/11/8
 * Time: 10:55
 */
namespace App\Models;

class EnterpriseGrayscaleReleaseLog extends Base
{
    /**
     * 表名称
     *
     * @var string
     */
    protected $table = 'enterprise_grayscale_release_log';

    /**
     * 撤销状态
     */
    const CANCEL_NO  = 0;   // 未撤销
    const CANCEL_YES = 1;   // 已撤销

    /**
     * 是否允许撤销
     */
    const CAN_CANCEL_YES  = 1;   // 未撤销,允许撤销
    const CAN_CANCEL_NO   = 0;   // 已撤销,不允许进行撤销

    /**
     * 黑名单
     *
     * @var array
     */
    protected $guarded = [];
}