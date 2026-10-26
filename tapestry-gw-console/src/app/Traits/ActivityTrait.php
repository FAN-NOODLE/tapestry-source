<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/5/22 0022
 * Time: 16:20
 */

namespace App\Traits;


Trait ActivityTrait
{
    function ActivityArrayToJson(&$arr, $key, $subKey, $delete = false)
    {
        $subValues = array_unique(array_column($arr[$key], $subKey));
        $arr[$subKey] = count($subValues) == 1 ? $subValues[0] : implode(',', $subValues);
        if ($delete) {
            unset($arr[$key]);
        }
    }
}