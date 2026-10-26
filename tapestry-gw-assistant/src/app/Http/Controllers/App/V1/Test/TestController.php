<?php


namespace App\Http\Controllers\App\V1\Test;


use App\Http\Controllers\Controller;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Cache;

class TestController extends Controller
{
    public function __invoke(Request $request)
    {
        $key = $request->get('key');
        $type = $request->get('type');
        if ($type == 1) {
            // 删除
            Cache::forget($key);
        } else {
            // 查询
            dd(Cache::get($key));
        }
    }
}