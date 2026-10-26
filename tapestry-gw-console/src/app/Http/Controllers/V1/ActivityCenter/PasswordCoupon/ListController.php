<?php
/*
 * @Author: zhaojie
 * @Date: 2020-02-06 22:49:25
 * @LastEditTime : 2020-02-10 14:18:25
 * @LastEditors  : Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \console\src\app\Http\Controllers\App\V1\AppCenter\PasswordCoupon\ListController.php
 */
/**
 * Created by PhpStorm.
 * User: admin
 * Date: 2019/6/24
 * Time: 17:35
 */

namespace App\Http\Controllers\V1\ActivityCenter\PasswordCoupon;


use App\Http\Controllers\V1\BaseController;
use App\Micros\Activity\PasswordCouponMicro;

class ListController extends BaseController
{
    public function index()
    {
        $input = $this->request->input();
        $input['enterprise_id'] = $this->enterpriseId;
        $response = (new PasswordCouponMicro())->list($input);

        return success($response);
    }
}
