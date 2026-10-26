<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/11/27
 * Time: 15:32
 */

namespace App\Common;


class Task
{
    /**
     * 成功码
     */
    const SUCCESS_CODE = 0;

    /**
     * 不循环
     */
    const NO_IS_LOOP = 2;

    /**
     * 循环
     */
    const IS_LOOP = 1;

    /**
     * 一天
     */
    const ONE_DAY = 1;

    /**
     * 一周
     */
    const ONE_WEEKS = 2;

    /**
     * 一个月
     */
    const ONE_MONTH = 3;

    /**
     * 3个月
     */
    const THREE_MONTHS = 4;

    /**
     * 六个月
     */
    const SIX_MONTHS = 5;

    /**
     * 一年
     */
    const ONE_YEAR = 6;

    /**
     * 一小时换算分钟
     */
    const SIX_MIN = 60;

    /**
     * 任务类型，生日祝福
     */
    const TASK_TYPE_BIRTHDAY = 3;

    /**
     * 员工列表角色
     */
    const TASK_ROLE_NAME = [
        '店长',
        '导购'
    ];
}
