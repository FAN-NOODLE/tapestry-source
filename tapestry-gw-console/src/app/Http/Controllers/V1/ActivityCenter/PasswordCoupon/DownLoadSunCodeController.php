<?php
/*
 * @Author: zhaojie
 * @Date: 2020-02-06 22:49:25
 * @LastEditTime: 2020-03-03 10:10:56
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \console\src\app\Http\Controllers\App\V1\AppCenter\PasswordCoupon\ListController.php
 */

namespace App\Http\Controllers\V1\ActivityCenter\PasswordCoupon;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\PasswordCoupon;

/**
 * 生成太阳码[已弃用]
 * @Author    DuanTuo
 * @DateTime  2020-03-03
 */
class DownLoadSunCodeController extends BaseController
{
    public function index()
    {
        $input = $this->request->input();
        $input['enterprise_id'] = $this->enterpriseId;

        $response = (new PasswordCoupon())->generateSunCode($input);
        
        return success();
    }
}
