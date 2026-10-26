<?php
/*
 * @Author: zhaojie
 * @Date: 2020-02-06 22:49:25
 * @LastEditTime : 2020-02-12 16:34:15
 * @LastEditors  : Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \console\src\app\Http\Controllers\App\V1\AppCenter\PasswordCoupon\AddController.php
 */

namespace App\Http\Controllers\V1\ActivityCenter\PasswordCoupon;

use App\Http\Controllers\V1\BaseController;
use App\Micros\Activity\PasswordCouponMicro;

class AddController extends BaseController
{
    public function index()
    {
        $input = $this->request->input();
        $input['enterprise_id'] = $this->enterpriseId;
        $input['coupon_creator'] = $this->enterpriseAdmin['username'] ?? '';
        $response = (new PasswordCouponMicro())->add($input);

        // 日志记录
        $this->writeOperatorLog('新增口令券');

        return success($response);
    }
}
