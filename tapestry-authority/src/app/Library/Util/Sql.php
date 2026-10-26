<?php


namespace App\Library\Util;


class Sql
{

    public static function put($sql)
    {
        static $isDEV = null;
        static $index = null;
        static $basePath = null;
        if($isDEV == null){
            $isDEV = env('APP_ENV') == 'dev';
        }
        if($isDEV == false){
            return ;
        }
        if (!$index) {
            $debug = debug_backtrace();
            $index = self::findPossibleCallIndex($debug);
        }
        $destination = storage_path('sql');
        if (!is_dir($destination)) {
            mkdir($destination, 0777, true);
        }
        $cate     = substr($sql, 0, 6);
        $callFile = debug_backtrace()[$index] ?? [];
        $file     = $destination . "/{$cate}.sql";
        $content  = static::makeTitle($sql, $callFile);
        file_put_contents($file, $content, FILE_APPEND);
    }

    private static function makeTitle($sql, array $file)
    {
        static $time;
        if (!$time) {
            $time = date("H:i:s");
        }

        $fileLine = sprintf("%s:%s",$file['file']??"",$file['line']??"");
        return sprintf('##调用时间:  %s', $time) . PHP_EOL .
            sprintf('##调用文件:  %s', trim(str_replace(base_path(), '', $fileLine), DIRECTORY_SEPARATOR)) . PHP_EOL .
            $sql . PHP_EOL.PHP_EOL;
    }

    private static function findPossibleCallIndex(array $debug)
    {
        unset($debug[0]);
        foreach ($debug as $key => $stack) {
            if (isset($stack['file'])
                &&
                strpos($stack['file'], base_path() . DIRECTORY_SEPARATOR . "app") !== false
            ) {
                return $key;
            }
        }
    }

}