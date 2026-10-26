<?php


namespace App\Library\Util;


/**
 * Class Logger
 * @method static info(...$params)
 * @method static debug(...$params)
 * @method static error(...$params)
 * @method static notice(...$params)
 * @method static html(...$params)
 * @method static api(...$params)
 */
class Logger
{
    static $file = '';
    static $line = 0;
    public static function __callStatic($method, $params)
    {
        static $isDev = null;
        if($isDev == null){
            $isDev = config('app.env') == 'dev';
        }
        if($isDev){
            $debug = current(debug_backtrace());
            $file = $debug['file'];
            self::$line = $debug['line'];
            self::$file = trim(str_replace(base_path(),'',$file),DIRECTORY_SEPARATOR);
            foreach ($params as $msg) {
                try {
                    static::putOne($msg, $method);
                } catch (\Throwable $e) {

                }

            }
            if ($params === []) {
                static::putOne("\n\n**************华丽丽的分割线**************\n\n", $method);
            }
        }

    }

    private static function putOne($msg, $method = 'info')
    {

        $suffix = '.log';
        if ($method == 'html') {
            ob_start();
            if ($msg instanceof \Throwable) {
                $msg = [
                    'code' => $msg->getCode(),
                    'msg'  => $msg->getMessage(),
                    'file' => $msg->getFile(),
                    'line' => $msg->getLine()
                ];
            }
            dump($msg);
            $msg = ob_get_contents();
            ob_end_clean();
            $suffix = '.html';
        } else {
            if (!is_scalar($msg)) {
                $msg = var_export($msg, true);
            }
        }
        $dest = storage_path('log') . DIRECTORY_SEPARATOR . date("m-d") . DIRECTORY_SEPARATOR;
        if (!is_dir($dest)) {
            mkdir($dest, 0777, true);
        }
        $title = self::makeTitle();
        file_put_contents(sprintf($dest . "/%s" . $suffix, $method), $title.
            $msg , FILE_APPEND);
    }

    private static function makeTitle()
    {
        static $time;
        if(!$time){
            $time = date("H:i:s");
        }
        return
            PHP_EOL.$time . " ==>  " .
            self::$file.':'.self::$line.PHP_EOL;
    }

}