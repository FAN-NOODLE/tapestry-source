<?php
/*
 * @Author: zhaojie
 * @Date: 2020-02-06 22:49:25
 * @LastEditTime : 2020-02-14 16:12:29
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

class DownLoadPasswordListController extends BaseController
{

    /**
     * 设置导出列宽
     * @var array
     */
    protected $colWidth = [
        'A' => 25, 'B' => 25
    ];

    public function index()
    {
        $input = $this->request->input();
        $input['enterprise_id'] = $this->enterpriseId;
        $input['admin_name'] = $this->enterpriseAdmin['username'];
        $input['admin_id'] = $this->enterpriseAdmin['id'];
        $input['route_url'] = ltrim($this->request->getRequestUri(), '\/gw-console\/v1\/');
        $res = (new PasswordCouponMicro())->passwordListExport($input);
        return success($res);
    }
}
