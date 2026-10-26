<?php

namespace App\Providers;

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
     * 创建不重复的HASH
     * @return string
     */
    public function makeRquestHash()
    {
        $hash = md5(uniqid()) . rand(10000, 99999);
        return $hash;
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
