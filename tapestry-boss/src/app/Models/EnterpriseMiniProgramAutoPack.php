<?php
/**
 * 小程序自动打包
 *
 */

namespace App\Models;

class EnterpriseMiniProgramAutoPack extends Base
{
    /**
     * 打包任务状态 ： 0 = 排队中
     * 0=排队中,1=进行中,2=成功,3=失败,4=需要登录
     */
    const TASK_STATUS_QUEUE = 0;

    /**
     * 打包任务状态 ： 1 = 进行中
     * 0=排队中,1=进行中,2=成功,3=失败,4=需要登录
     */
    const TASK_STATUS_ING = 1;

    /**
     * 打包任务状态 ： 2 = 成功
     * 0=排队中,1=进行中,2=成功,3=失败,4=需要登录
     */
    const TASK_STATUS_SUCCESS = 2;

    /**
     * 打包任务状态 ： 3 = 失败
     * 0=排队中,1=进行中,2=成功,3=失败,4=需要登录
     */
    const TASK_STATUS_FAILED = 3;

    /**
     * 打包任务状态 ： 4 = 需要登录
     * 0=排队中,1=进行中,2=成功,3=失败,4=需要登录
     */
    const TASK_STATUS_NEED_LOGIN = 4;

    /**
     * shop端 小程序版本状态值
     * 0：未发布  1：审核中  2：发布  3：审核失败 4:自动打包中 5:预发布 6:已失效 7: 审核成功
     */
    const SHOP_VERSION_AUDIT_ING = 1;
    const SHOP_VERSION_PUBLISH = 2;
    const SHOP_VERSION_AUDIT_ERROR = 3;
    const SHOP_VERSION_PACK_ING = 4;
    const SHOP_VERSION_TIME_PUBLISH = 5;
    const SHOP_VERSION_HAS_EXPIRED = 6;
    const SHOP_VERSION_AUDIT_SUCCESS = 7;

    /**
     * 表名称
     *
     * @var string
     */
    protected $table = 'enterprise_mini_program_auto_pack';


    /**
     * @var array
     */
    protected $fillable = [
        'id',
        'enterprise_id',
        'type',
        'task_id',
        'status',
        'desc',
        'version',
        'task_config',
        'login_url',
        'version_id'
    ];
}
