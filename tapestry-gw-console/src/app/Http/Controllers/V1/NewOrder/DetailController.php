<?php

namespace App\Http\Controllers\V1\NewOrder;

use App\Common\Constant;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\Open\Enterprise;
use CrsSdk\Shop\Order;
use CrsSdk\Work\Partner;
use Illuminate\Http\JsonResponse;
use Illuminate\Support\Facades\Validator;
use Exception;

class DetailController extends BaseController
{
    /**
     * 订单详情
     *
     * @return JsonResponse
     * @throws Exception
     */
    public function index()
    {
        // 接收参数
        $params = request()->only(['order_id', 'refund_apply_id']);
        // 验证规则
        $rules = [
            'order_id'        => 'required|integer',
            'refund_apply_id' => 'integer'
        ];
        // 提示语
        $messages = [
            'order_id.required'       => '[参数]订单ID[order_id]不能为空',
            'order_id.integer'        => '[参数]订单ID[order_id]格式不正确',
            'refund_apply_id.integer' => '[参数]售后单ID[refund_apply_id]格式不正确',
        ];
        // 参数验证
        Validator::make($params, $rules, $messages)->validate();

        $params['enterprise_id'] = $this->enterpriseId;
        // 调用Shop SDK获取订单详情
        $data = (new Order())->orderDetail($params);
        // 企业详情
        $enterpriseDetail = (new Enterprise())->detail(['enterprise_hash' => $this->enterprise['hash']])['data'];

        // 订单详情是否展示"补推Erp"按钮状态
        $showPushErp = Constant::SHOW_PUSH_ERP_FALSE;
        if (!empty($enterpriseDetail['erp_package']) && $data['third_order_id'] === '0') {
            $showPushErp = Constant::SHOW_PUSH_ERP_TRUE;
        }

        $data['show_push_erp'] = $showPushErp;

        return success($this->joinPartnerInfo($data));
    }

    /**
     * 处理加入合伙人信息
     * User: Hui Yang
     * Email: 2361313833@qq.com
     * DateTime: 2020/12/10 19:54
     * @param $data
     * @return array
     */
    public function joinPartnerInfo($data): array
    {
        if(!isset($data['store_info']) || !empty($data['store_info']['store_staff_mobile'])){
            return $data;
        }
        $requestParam = [
            'enterprise_id' => $data['enterprise_id'],
            'id' => $data['store_info']['store_staff_id']
        ];

        $partnerInfo = (new Partner())->detail($requestParam);
        if(isset($partnerInfo['partner']) && isset($partnerInfo['partner']['partner_id'])){
            $data['store_info']['store_staff_mobile'] = $partnerInfo['partner']['partner_mobile'];
            $data['store_info']['store_staff_name'] = $partnerInfo['partner']['partner_name'];
        }
        return $data;
    }
}
