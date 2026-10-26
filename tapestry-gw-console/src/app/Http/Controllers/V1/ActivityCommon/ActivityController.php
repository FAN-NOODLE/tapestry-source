<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/5/12 0012
 * Time: 11:01
 */

namespace App\Http\Controllers\V1\ActivityCommon;

use App\Http\Controllers\V1\BaseController;

class ActivityController extends BaseController
{

    Const COMMON_CONTROLLER_NAMESPACE = "Common";

    public function index()
    {
        $param = $this->request->input();
        // 类型转命名空间
        $type = self::typeToNameSpace($param['type']);

        // 方法转控制器
        $method = self::methodToController($param['method']);
        // 分发到控制器
        $implementsController = __NAMESPACE__ . "\\" . $type . "\\" . $method . "Controller";

        if (@class_exists($implementsController)) { //找到对应控制器
            return (new $implementsController)->index();
        } else { // 进入公用控制器
            $implementsController = __NAMESPACE__ . "\\" . self::COMMON_CONTROLLER_NAMESPACE . "\\" . $method . "Controller";
            return (new $implementsController)->index();
        }

    }

    public static function typeToNameSpace($type)
    {
        return self::snakeToCamel($type);
    }

    public static function methodToController($method)
    {
        return self::snakeToCamel($method);
    }

    public static function snakeToCamel($snake)
    {
        $arr = explode('-', $snake);
        $res = $arr[0];
        $len = count($arr);
        if ($len > 1) {
            for ($i = 1; $i < $len; $i++) {
                $res .= ucfirst($arr[$i]);
            }
        }
        return ucfirst($res);
    }

}
