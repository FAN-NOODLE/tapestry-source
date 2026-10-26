<?php

namespace App\Providers;

use App\Library\Util\Sql;
use Illuminate\Support\ServiceProvider;

class AppServiceProvider extends ServiceProvider
{
    /**
     * Bootstrap any application services.
     *
     * @return void
     */
    public function boot()
    {
        //记录sql到日志
        if (config('app.env') !== 'prod') {
            \DB::listen(function ($query) {
                $sql = $this->getQuerySql($query);
                Sql::put($sql);
                \Log::channel('sqllog')->info("执行sql：" . $sql . ",执行时间: " . $query->time . "ms; ");
            });
        }
    }

    /**
     * Register any application services.
     *
     * @return void
     */
    public function register()
    {
        define('LOG_ID', $this->makeRquestHash());
    }

    /**
     * @return string
     */
    public function makeRquestHash()
    {
        return md5(uniqid()) . rand(10000, 99999);
    }

    /**
     * 获取执行sql
     * @return string
     */
    public function getQuerySql($query)
    {
        $sql = str_replace("?", "'%s'", $query->sql);

        return vsprintf($sql, $query->bindings);
    }
}
