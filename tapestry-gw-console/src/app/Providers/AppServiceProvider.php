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
        //记录sql到日志
        \DB::listen(function($query){
            $sql=$this->getQuerySql($query);
            \Log::channel('sqllog')->info($sql.'---'.$query->time);
        });
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
