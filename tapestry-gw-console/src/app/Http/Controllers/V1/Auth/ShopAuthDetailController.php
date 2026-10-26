<?php
/**
 * Created by PhpStorm.
 * User: daijun
 * Date: 2018/11/27
 * Time: 下午8:46
 */

namespace App\Http\Controllers\V1\Auth;

use App\Http\Controllers\V1\BaseController;
use Illuminate\Http\Request;
use CrsSdk\Boss\Auth;

class ShopAuthDetailController extends BaseController
{

    /**
     * 获取商城小程序授权详情
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $data = [
            'client'        => 'console',
            'enterprise_id' => $this->enterpriseId,
            'back_url'      => $request->input('back_url', '')
        ];

        $result = (new Auth())->shopAuthDetail($data);

        $this->writeOperatorLog("商城小程序授权");

        return success($result['data'] ?? []);
    }
}
