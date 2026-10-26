<?php
/**
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/2/28
 * Time: 16:38
 */
namespace App\Models;

class TaskType extends Base
{
    /**
     * 表名
     * @var string
     */
    protected $table = 'task_type';

    /**
     * @var array
     */
    protected $fillable = [
        'id',
        'name',
        'describe',
        'serve_type'
    ];
}
