<?php
/**
 * Created by PhpStorm.
 * User: wangle@vchangyi.com
 * Date: 2021/04/26
 */

namespace App\Models\Scrm;

class CustomerInherit extends ScrmInit
{
    /**
     * 表名
     * 成员和外部联系人关系表
     *
     * @var string
     */
    protected $table = 'customer_inherit';

    /**
     * 主键
     *
     * @var string
     */
    public $primaryKey = 'id';

}
