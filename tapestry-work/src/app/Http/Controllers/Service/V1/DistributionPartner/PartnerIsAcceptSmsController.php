<?php
/**
 * @author  zhaopengfei<zhaopengfei@vchangyi.com>
 * @date    2020/12/31 16:07
 */

namespace App\Http\Controllers\Service\V1\DistributionPartner;


use App\Http\Controllers\Service\V1\BaseController;
use App\Models\PartnerEnterpriseMap;
use Illuminate\Http\Request;

class PartnerIsAcceptSmsController extends BaseController
{
        public function index(Request $request)
        {
            $params = $this->getParams($request);
            PartnerEnterpriseMap::query()
                ->where('enterprise_id',$params['enterprise_id'])
                ->where('partner_id',$params['partner_id'])
                ->update(['sms_switch'=>$params['sms_switch']]);
            return success();
        }

    private function getParams(Request $request): array
    {
        $rules = [
            'sms_switch' => 'required|in:1,2',
            'partner_id' => 'required',
            'enterprise_id' => 'required',
        ];

        $message = [
            'sms_switch.*' => '手机号码保护不合法',
            'enterprise_id.*' => '企业ID不合法',
            'partner_id.*'=>'合伙人ID不合法',
        ];
        $input = $this->validate($request, $rules, $message);
        return $input;
    }
}