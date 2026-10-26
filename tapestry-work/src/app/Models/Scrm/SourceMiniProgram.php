<?php
/**
 * Created by PhpStorm.
 * User: wangle@vchangyi.com
 * Date: 2021/04/26
 */

namespace App\Models\Scrm;

use App\Models\Channel;

class SourceMiniProgram extends ScrmInit
{

    /**
     * 表名
     * 外部联系人标签表
     *
     * @var string
     */
    protected $table = 'source_mini_program';


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

