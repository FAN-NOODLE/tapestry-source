<?php
/**
 * Created By PhpStorm
 * @Date: 2021/6/25
 * @Time: 9:15
 * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
 * @File: FissionActivity.php
 */

namespace App\Models\Scrm;


class DrainageContacts extends ScrmInit
{
    /**
     * 表名
     *
     * @var string
     */
    protected $table = 'drainage_contacts';


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