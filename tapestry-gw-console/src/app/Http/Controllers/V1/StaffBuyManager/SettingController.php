<?php

namespace App\Http\Controllers\V1\StaffBuyManager;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\StaffRight;

class SettingController extends BaseController
{
    /**
     * 获取内购员工设置数据
     * @return mixed
     * @throws \Exception
     */
    public function index()
    {
        /**
         * | enterprise_id | int | __是__ | 企业id |
         */
        $params['enterprise_id'] = $this->enterpriseId;

        return (new StaffRight())->staffSetting($params);
    }
}
