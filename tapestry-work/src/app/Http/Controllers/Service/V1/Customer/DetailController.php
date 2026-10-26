<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/27
 * Time: 21:34
 */

namespace App\Http\Controllers\Service\V1\Customer;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Services\CustomerService;
use Illuminate\Http\Request;

class DetailController extends BaseController
{
    /**
     * 客户基础信息
     *
     * @param Request $request
     *
     * @return array|mixed
     * @throws \Exception | Error
     */
    public function index(Request $request)
    {
        // 企业id
        $enterpriseId = $request->get('enterprise_id', 0);
        if (!$enterpriseId) {
            return error(600000);
        }

        // 企业domain
        $qyDomain = $request->get('qy_domain', '');
        if (!$qyDomain) {
            return error(600003);
        }

        // 应用标识
        $appIdentifier = $request->get('app_identifier', '');
        if (!$appIdentifier) {
            return error(600004);
        }

        // 客户id
        $customerId = $request->get('user_id', 0);
        // 外部联系人记录id
        $ecId = $request->get('ec_id', '');

        // 获取客户详情数据
        $data = (new CustomerService())->getCustomerDetail($enterpriseId, $qyDomain, $appIdentifier, $customerId, $ecId);

        return success(['data' => $data]);
    }
}
