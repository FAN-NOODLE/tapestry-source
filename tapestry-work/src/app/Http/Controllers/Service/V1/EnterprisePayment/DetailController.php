<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/12/3
 * Time: 3:12 PM
 */

namespace App\Http\Controllers\Service\V1\EnterprisePayment;

use App\Http\Controllers\Service\V1\BaseController;
use App\Http\Requests\EnterprisePaymentRequest;
use App\Exceptions\Error;
use App\Models\EnterprisePayment;
use App\Services\EnterprisePaymentService;

class DetailController extends BaseController
{
    /**
     * 获取企业微信红包活动详情
     *
     * @return array
     */
    public function index()
    {
        $params = $this->request->input();

        // 验证数据
        (new EnterprisePaymentRequest())->commonId($params);

        $detail = EnterprisePayment::query()
            ->where('enterprise_id', $params['enterprise_id'])
            ->where('id', $params['id'])
            ->first();
        if (empty($detail)) {
            throw new Error(110009, '活动信息不存在');
        }

        $detail         = $detail->toArray();
        $statisticalNum = (new EnterprisePaymentService())->getStatisticalNum($detail);

        $detail = array_merge($detail, $statisticalNum);

        return success(['data' => $detail]);
    }
}
