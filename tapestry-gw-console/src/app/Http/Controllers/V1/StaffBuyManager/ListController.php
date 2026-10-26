<?php


namespace App\Http\Controllers\V1\StaffBuyManager;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\StaffRight;

class ListController extends BaseController
{
    /**
     * 获取内购员工列表
     * @return mixed
     * @throws \Exception
     */
    public function index()
    {
        /**
         * | qy_domain | string | __是__ | 企业标识 |
         * | app_identifier | string | __是__ | 应用标识 |
         * | enterprise_id | int | __是__ | 企业id |
         */
        $params['enterprise_id'] = $this->enterpriseId;
        $params['qy_domain'] = $this->qyDomain;

        return (new StaffRight())->staffList($params);
    }
}
