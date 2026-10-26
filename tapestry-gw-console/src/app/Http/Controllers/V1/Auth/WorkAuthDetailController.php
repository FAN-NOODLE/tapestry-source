<?php
/**
 * Created by PhpStorm.
 * User: daijun
 * Date: 2018/11/27
 * Time: 下午8:46
 */

namespace App\Http\Controllers\V1\Auth;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Boss\Auth;
use Illuminate\Http\Request;

class WorkAuthDetailController extends BaseController
{

    /**
     * 获取企业微信小程序授权详情
     *
     * @param Request $request
     * @return mixed
     */
    public function index(Request $request)
    {
        $data = [
            'enterprise_id' => $this->enterpriseId,
            'back_url'      => $request->input('back_url', '')
        ];

        $result = (new Auth())->workAuthDetail($data);

        // 日志埋点
        $this->writeOperatorLog("企业微信小程序授权");

        return success($result);
    }
}
