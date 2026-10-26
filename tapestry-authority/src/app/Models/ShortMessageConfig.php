<?php
/**
 * Created By PhpStorm
 * @Date: 2021/4/14
 * @Time: 20:06
 * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
 * @File: ShortMessageConfig.php
 */

namespace App\Models;

class ShortMessageConfig extends Base
{
    protected $table = 'short_message_config';

    protected $fillable = [
        'id',
        'enterprise_id',
        'name',
        'brand_code',
        'sms_type',
        'config',
        'created_at',
        'updated_at',
        'deleted_at',
    ];
}