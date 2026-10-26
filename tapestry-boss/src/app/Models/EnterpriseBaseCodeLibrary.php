<?php
/**
 * Created by PhpStorm.
 * Date: 2018/11/15
 * Time: 4:05 PM
 */

namespace App\Models;

class EnterpriseBaseCodeLibrary extends Base
{
    /**
     * 表名称
     *
     * @var string
     */
    protected $table = 'enterprise_base_code_library';

    /**
     * @var array
     */
    protected $fillable = [
        'id',
        'type',
        'version',
        'description',
        'template_id'
    ];
}
