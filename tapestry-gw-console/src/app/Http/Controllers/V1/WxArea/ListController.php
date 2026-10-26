<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2019/11/18
 * Time: 14:53
 */
namespace App\Http\Controllers\V1\WxArea;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\WxArea;
use Illuminate\Http\Request;

class ListController extends BaseController
{
    /**
     * 会员等级列表
     *
     * @return array
     */
    public function index(Request $request)
    {
        // 参数列表
        $params = $request->input();

        // 调用sdk
        $data = (new WxArea())->wxAreaList($params);

        return success($data);
    }
}