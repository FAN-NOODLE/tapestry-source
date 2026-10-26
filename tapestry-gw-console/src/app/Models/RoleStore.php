<?php
/*
 * Created by PhpStorm.
 * User: wangce
 * Date Time: 2019/3/20 15:23
 */

namespace App\Models;

class RoleStore extends Base
{

    /**
     * 角色门店关系表
     *
     * @var string
     */
    protected $table = 'role_store';

    /**
     * @var array
     */
    protected $fillable = [
        'role_id',
        'story_id'
    ];

}
