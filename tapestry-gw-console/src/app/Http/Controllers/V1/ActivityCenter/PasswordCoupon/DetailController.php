<?php
/*
 * @Author: zhaojie
 * @Date: 2020-02-06 22:49:25
 * @LastEditTime : 2020-02-07 16:15:33
 * @LastEditors  : Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \console\src\app\Http\Controllers\App\V1\AppCenter\PasswordCoupon\DetailController.php
 */

namespace App\Http\Controllers\V1\ActivityCenter\PasswordCoupon;


use App\Http\Controllers\V1\BaseController;
use App\Micros\Activity\PasswordCouponMicro;


class DetailController extends BaseController
{
    public function index()
    {
        $input = $this->request->only(['id']);

        $response = (new PasswordCouponMicro())->detail($input);

        return success($response);
    }
}