<?php
/*
 * @Author: zhaojie
 * @Date: 2020-02-06 22:49:25
 * @LastEditTime: 2020-02-19 15:32:15
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \console\src\app\Http\Controllers\App\V1\AppCenter\PasswordCoupon\ListController.php
 */

namespace App\Http\Controllers\V1\ActivityCenter\PasswordCoupon;


use App\Http\Controllers\V1\BaseController;
use App\Micros\Activity\PasswordCouponMicro;

class DownLoadPosterMoreStatusController extends BaseController
{
    public function index()
    {
        $input = $this->request->input();
        $input['enterprise_id'] = $this->enterpriseId;
        $res = (new PasswordCouponMicro())->DownloadPosterMoreStatus($input);
        return success($res);
    }
}
