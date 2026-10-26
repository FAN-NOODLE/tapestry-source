<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/26
 * Time: 10:04
 */

namespace App\Http\Controllers\V1\UserLevel;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\UserLevel;
use Illuminate\Http\Request;

class DetailController extends BaseController
{
    /**
     * 会员等级详情
     *
     * @param Request $request
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function index(Request $request)
    {
        // 等级id
        $id = $request->post('id', 0);
        if (!$id) {
            error(320001);
        }

        // 参数列表
        $params = [
            'enterprise_id' => $this->enterpriseId,
            'id'            => $id
        ];
        // 调用sdk
        $data = (new UserLevel())->detail($params);

        return success($data);
    }
}
