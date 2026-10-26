<?php
/*
 * @Author: zhaojie
 * @Date: 2020-02-06 22:49:25
 * @LastEditTime : 2020-02-10 11:16:27
 * @LastEditors  : Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \console\src\app\Http\Controllers\App\V1\AppCenter\PasswordCoupon\UpdateController.php
 */
/**
 * Created by PhpStorm.
 * User: admin
 * Date: 2019/6/24
 * Time: 17:36
 */

namespace App\Http\Controllers\V1\ActivityCenter\PasswordCoupon;


use App\Http\Controllers\V1\BaseController;
use App\Micros\Activity\PasswordCouponMicro;


class UpdateController extends BaseController
{

    public function index()
    {
        $input = $this->request->input();
        $input['enterprise_id'] = $this->enterpriseId;

        $response = (new PasswordCouponMicro())->update($input);

        return success($response);
    }

    
}
