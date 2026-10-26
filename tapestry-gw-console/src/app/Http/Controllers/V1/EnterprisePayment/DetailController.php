<?php
/**
 * 获取企业微信红包活动详情
 * User: liurui
 * Date: 2018/11/26
 * Time: 14:10
 */

namespace App\Http\Controllers\V1\EnterprisePayment;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\EnterprisePayment;
use Illuminate\Http\Request;

class DetailController extends BaseController
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

        $data = (new EnterprisePayment())->detail($params);

        return success($data['data']);
    }
}
