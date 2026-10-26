<?php
/*
 * @Author: zhaojie
 * @Date: 2020-02-06 22:49:25
 * @LastEditTime: 2020-02-24 20:46:57
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \console\src\app\Http\Controllers\App\V1\AppCenter\PasswordCoupon\ListController.php
 */

namespace App\Http\Controllers\V1\ActivityCenter\PasswordCoupon;


use App\Http\Controllers\V1\BaseController;
use App\Micros\Activity\PasswordCouponMicro;

class GeneratePostersController extends BaseController
{
    public function index()
    {
        $params = $this->getParams();

        $res = (new PasswordCouponMicro())->generatePosters($params);

        return success($res);
    }

    protected function getParams()
    {
        // 验证参数
        $rule = [
            'password_coupon_id' => 'required|integer',
            'num'                => 'required|integer',
        ];

        $messages = [
            'password_coupon_id.*' => '口令券主键必填',
            'num.*'                => '下载海报数量必填',
        ];

        $input = $this->validate($this->request, $rule, $messages);

        $params = [
            'password_coupon_id' => $input['password_coupon_id'],
            'num'                => $input['num'],
            'enterprise_id'      => $this->enterpriseId,
            'admin_name'         => $this->enterpriseAdmin['username'] ?? 'username',
            'admin_id'           => $this->enterpriseAdmin['id'],
            'route_url'          => ltrim($this->request->getRequestUri(), '\/gw-console\/v1\/')
        ];

        return $params;
    }
}
