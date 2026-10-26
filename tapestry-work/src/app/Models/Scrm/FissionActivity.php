<?php
/**
 * Created By PhpStorm
 * @Date: 2021/6/25
 * @Time: 9:15
 * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
 * @File: FissionActivity.php
 */

namespace App\Models\Scrm;


class FissionActivity extends ScrmInit
{
    /**
     * 活动状态：1-上架中，2-已下架
     */
    const FISSION_ACTIVITY_ON = 1;
    const FISSION_ACTIVITY_OFF = 2;

    /**
     * 表名
     *
     * @var string
     */
    protected $table = 'fission_activity';


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