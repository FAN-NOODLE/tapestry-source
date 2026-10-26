<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;

class IndexController extends Controller
{
    /**
     * 首页
     *
     * 仅在开发环境才会输入日志
     *
     * @param Request $request
     * @return array|mixed
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/12/7 10:48 AM
     */
    public function index(Request $request)
    {
        $env = config('app.env');
        $delete = $request->get('delete', false);
        $password = $request->get('p', false);

        // 如果是开发和测试环境才输出日志
        if ($env === 'dev' || $env === 'test') {

            if ($password != date('md')) {
                exit('对不起，没有日志查看的权限');
            }

            $dateFile = empty($request->get('file')) ? date('Y/m/d/\h-H', time()).'.log' : $request->get('file');
            $file = storage_path('logs/'.$dateFile);

            if (!empty($delete)) {
                if (file_exists($file)) {
                    unlink($file);
                }
                return redirect('/?p='.date('md'));
            }

            if (!file_exists($file)) {
                return view('logging')
                    ->with('file', $dateFile)
                    ->with('size', 0)
                    ->with('line', 0)
                    ->with('data', '没有任何内容');
            }

            $fileSize = filesize($file);
            $fileSize = $fileSize / 1024;
            $fileSize = floor($fileSize);

            // 日志大于2M直接删除掉
            if ($fileSize > 2048) {
                unlink($file);
                return redirect('/?p='.date('md'));
            }

            $handle = fopen($file, 'rb');

            $line = 0;
            while (!feof($handle)) {
                if (fgets($handle)) {
                    $line++;
                }
            }

            $data = file_get_contents($file);
            return view('logging')
                ->with('file', $dateFile)
                ->with('size', $fileSize)
                ->with('line', $line)
                ->with('data', $data);
        }

        return 'console.'.$env;
    }

    /**
     * 查看运行日志
     *
     * @param $date string 日期:2019-01-01
     * @return string
     */
    public function runtime($date = 'today')
    {
        if (config('app.env') !== 'dev') {
            return '对不起，没有日志查看的权限';
        }

        $date = $date === 'today' ? date('Y/m/d/\h-H') : $date;

        // 日志文件路径
        $logPath = storage_path('logs/runtime/'.$date).'.log';

        // 日志文件不存在
        if (!file_exists($logPath)) {
            return '运行日志不存在';
        }

        $log = file_get_contents($logPath);

        return '<pre>'.$log.'</pre>';
    }
}
