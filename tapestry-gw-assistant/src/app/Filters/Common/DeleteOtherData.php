<?php
/**
 * Created by PhpStorm.
 * Function: 区分售后类型
 * User: 龚国新
 * Date: 2019/11/21
 * Time: 12:16
 */

namespace App\Filters\Common;


use App\Filters\Filter;
use Closure;

class DeleteOtherData implements Filter
{

    public function handle($content, Closure $next)
    {
        if(!empty($content['other'])){
            unset($content['other']);
        }

        return $next($content);
    }
}