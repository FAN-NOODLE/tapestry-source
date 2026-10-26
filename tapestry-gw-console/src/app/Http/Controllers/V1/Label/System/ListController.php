<?php
/**
 * Created by sublime.
 * User: zhaoqi
 * Date: 2019/01/17
 * Time: 10:00
 */
namespace App\Http\Controllers\V1\Label\System;

use CrsSdk\Shop\SystemLabel;
use App\Http\Controllers\App\V1\Label\LabelBaseController;

class ListController extends LabelBaseController
{
    /**
     * 系统标签列表
     *
     * @return array
     */
    public function index()
    {
        // 默认值
        $params['enterprise_id'] = $this->enterpriseId;
        return (new SystemLabel())->list($params);
    }
}
