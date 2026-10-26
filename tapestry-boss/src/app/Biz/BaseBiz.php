<?php

/*
 * This file is part of the car/chedianai_bc.
 *
 * (c) chedianai <i@chedianai.com>
 *
 * This source file is subject to the MIT license that is bundled
 * with this source code in the file LICENSE.
 */

namespace App\Biz;

use Closure;
use Illuminate\Support\Facades\DB;

class BaseBiz
{
    /**
     * 异常回滚事务的闭包.
     *
     * @param Closure $callback
     * @param Closure $errorCallBack
     *
     * @return mixed
     *
     * @throws \Exception
     *
     * @author Caikeal <caikeal@qq.com>
     */
    protected function transaction(Closure $callback, Closure $errorCallBack = null)
    {
        DB::beginTransaction();
        try {
            $result = $callback($this);
            DB::commit();

            return $result;
        } catch (\Exception $e) {
            $errorCallBack && $errorCallBack($this);
            DB::rollback();
            throw $e;
        }
    }
}
