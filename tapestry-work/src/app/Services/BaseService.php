<?php

namespace App\Services;

use Closure;
use Illuminate\Support\Facades\DB;

/**
 *  Basic Service
 * Class BaseService.
 */
abstract class BaseService
{
    /**
     * 异常回滚事务的闭包.
     *
     * @param Closure        $callback
     * @param null | Closure $errorCallback
     *
     * @return mixed
     *
     * @throws \Exception
     */
    public function transaction(Closure $callback, Closure $errorCallback = null)
    {
        DB::beginTransaction();
        try {
            $result = $callback($this);
            DB::commit();

            return $result;
        } catch (\Exception $e) {
            DB::rollback();
            if (is_callable($errorCallback)) {
                $errorCallback($this);
            }

            throw $e;
        }
    }
}
