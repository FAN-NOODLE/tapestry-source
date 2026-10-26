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
use CrsSdk\Shop\PasswordCoupon;

class DownLoadPosterController extends BaseController
{
    public function index()
    {
        $input = $this->request->input();
        $input['enterprise_id'] = $this->enterpriseId;
        $input['admin_name'] = $this->enterpriseAdmin['username'];
        $input['admin_id'] = $this->enterpriseAdmin['id'];
        $res  = (new PasswordCoupon())->passwordPosterExport($input);
        return success($res);
    }
}
