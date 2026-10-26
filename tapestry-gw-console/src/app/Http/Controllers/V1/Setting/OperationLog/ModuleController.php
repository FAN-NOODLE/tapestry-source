<?php
/*
 * Created by PhpStorm.
 * User: wangce
 * Date Time: 2019/2/25 11:17
 */

namespace App\Http\Controllers\V1\Setting\OperationLog;

use App\Http\Controllers\V1\BaseController;
use App\Models\ApiPower;

class ModuleController extends BaseController
{

    /**
     * 获取全部模块列表
     *
     * @return array
     */
    public function index()
    {
        $modules = [
//            '首页',
//            '店铺',
//            '商品',
//            '订单',
//            '任务',
//            '客户',
//            '组织',
//            '应用',
//            '设置'
            '管理后台配置',
            '智慧导购',
            '企业微信'
        ];

        return success(['list' => $modules]);
    }

}
