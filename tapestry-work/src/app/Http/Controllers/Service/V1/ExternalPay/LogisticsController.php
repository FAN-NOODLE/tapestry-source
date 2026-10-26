<?php

namespace App\Http\Controllers\Service\V1\ExternalPay;

use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Express\Express;
use App\Models\ScrmEnterprise\ExternalPayOrder;
use Illuminate\Http\JsonResponse;

class LogisticsController extends BaseController
{


    /**
     * 物流公司
     *
     * @return void
     */
    public function logisticsCompany(): JsonResponse
    {
        return successJson(Express::LOGISTICS_KEY);
    }

    /**
     * 查询物流信息
     *
     * @return JsonResponse
     */
    public function logisticsTrack(): JsonResponse
    {

        $params  = $this->request->all();
        $rules   = [
            'enterprise_id' => 'required',
            'id'            => 'required',
        ];
        $message = [
            'enterprise_id.required' => '企业ID不能为空',
            'mobile.required'        => '手机号不能为空',
            'id.required'            => '订单ID不能为空',
        ];
        $this->validate($this->request, $rules, $message);

        $res = ExternalPayOrder::query()->setConn($params['enterprise_id'])->where('id', $params['id'])->first();

        $arr = [
            'logistics_name'  => $res->logistics_name,
            'logistics_no'    => $res->logistics_no,
            'logistics_track' => $res->logistics_track,
        ];

        if ($res->sign_at) {
            return successJson($arr);
        } else {

            $receiverPhone = $res->contact_info['phone'] ?? '';
            $receiverPhone = $receiverPhone?substr($receiverPhone,-4):'';
            $data          = (new Express())->queryExpress($res->logistics_com, $res->logistics_no, '', $receiverPhone);
            if ($data['data']) {
                $arr['logistics_track'] = $data['data']['list'];
                $res->logistics_track = json_encode($arr['logistics_track']);
                if ($data['data']['status_detail']) {
                    $res->sign_at = strtotime($data['data']['list'][0]['datetime']);
                    $res->order_state = ExternalPayOrder::ORDER_STATE_SUCCESS;
                }
                $res->save();
            }
            return successJson($arr);
        }
    }

}
