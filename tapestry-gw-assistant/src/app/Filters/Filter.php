<?php
/**
 * Created by PhpStorm.
 * User: 龚国新
 * Date: 2019/08/14
 * Time: 14:34
 */

namespace App\Filters;


use Closure;

interface Filter
{
    public function handle($content, Closure $next);
}