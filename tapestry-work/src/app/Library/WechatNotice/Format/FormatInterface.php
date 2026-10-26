<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2021/4/9 0009
 * Time: 10:16
 */

namespace App\Library\WechatNotice\Format;


interface FormatInterface
{
    public function format(string $content):array;
}