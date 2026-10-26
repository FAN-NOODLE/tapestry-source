<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2019/7/2
 * Time: 10:51
 */
namespace App\Models;

class Setting extends Base
{
    /**
     * 角色表
     *
     * @var string
     */
    protected $table = 'setting';

    /**
     * @var array
     */
    protected $fillable = [
        'package',
        'config',
        'describe'
    ];
}