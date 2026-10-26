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
        \DB::listen(function($query){
            $sql=$this->getQuerySql($query);
            Sql::put($sql);
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
        // define('LOG_ID', $this->makeRquestHash());
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
     * @param $query
     * @return mixed|string
     * @author Zhujiaqi&huanWang
     * @time 2020-10-30 18:02:14
     */
    public function getQuerySql($query)
    {
        if (empty($query->bindings)) {
            $sql = str_replace("?", "'%s'", $query->sql);
            return $sql;
        }

        array_walk($query->bindings, function (&$value) {
            if (is_string($value)) {
                $value = "'" . $value . "'";
            }
        });

        $sql = str_replace("?", "%s", $query->sql);
        return vsprintf($sql, $query->bindings);
    }
}
