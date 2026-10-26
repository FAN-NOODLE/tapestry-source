<?php
/**
 * Created by PhpStorm.
 * User: huangxm
 * Date: 2020-04-15
 * Time: 12:55
 */

namespace App\Models;

class StoreManager extends BaseModel
{
    protected $table = 'store_manager';

    /**
     * 是否是第一次登陆此小程序 (是)
     * 
     * @var int
     */
    const IS_FIRST_LOGIN_YES = 1;

    /**
     * 是否是第一次登陆此小程序 (否)
     * 
     * @var int
     */
    const IS_FIRST_LOGIN_NO = 0;
}
