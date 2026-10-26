<?php
/**
 * 企业详情
 * User: daijun
 * Date: 2018/12/11
 * Time: 14:10
 */

namespace App\Http\Controllers\V1\Home;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Boss\Enterprise;

class AuthDetailController extends BaseController
{

    /**
     * 获取企业是否授权通讯录
     *
     * @return array
     */
    public function index()
    {

        $result = (new Enterprise())->detail(['enterprise_id' => $this->enterpriseId]);

        if (isset($result) && isset($result['data']) && isset($result['data']['domain']) && !empty($result['data']['domain'])) {

            return success(['data' => ['contact_auth_status' => 1]]);

        } else {

            return success(['data' => ['contact_auth_status' => 2]]);
        }
    }
}
