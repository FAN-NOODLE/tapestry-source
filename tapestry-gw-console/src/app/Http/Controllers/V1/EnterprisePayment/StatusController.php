<?php
/**
 * 更新企业微信红包活动状态
 * User: liurui
 * Date: 2018/11/26
 * Time: 14:10
 */

namespace App\Http\Controllers\V1\EnterprisePayment;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\EnterprisePayment;
use Illuminate\Http\Request;

class StatusController extends BaseController
{

    /**
     * @param Request $request
     *
     * @return mixed
     * @throws \Exception
     */
    public function index(Request $request)
    {

        // 接值
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        // 日志记录
        $this->writeOperatorLog('更新企业微信红包活动状态');

        (new EnterprisePayment())->status($params);

        return success();
    }
}
