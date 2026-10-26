<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/12/3
 * Time: 3:12 PM
 */

namespace App\Http\Controllers\App\V1\EnterprisePayment;

use App\Http\Controllers\App\V1\BaseController;
use CrsSdk\Work\EnterprisePayment;
use Illuminate\Http\Request;

class RechargeRecordController extends BaseController
{
    /**
     * 充值记录接口
     *
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request)
    {
        $params = $request->only(['enterprise_id', 'page', 'size']);

        $params['page'] = $params['page'] ?? 1;
        $params['size'] = $params['size'] ?? 10;

        $params['search'] = [];
        $params['search']['enterprise_id'] = $params['enterprise_id'];
        $params['search']['type'] = 1; // 充值记录
        unset($params['enterprise_id']);

        $list = (new EnterprisePayment())->rechargeRecord($params);

        return $list;
    }
}
