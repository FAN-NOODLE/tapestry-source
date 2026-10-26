<?php
/**
 * @author  zhaopengfei<zhaopengfei@vchangyi.com>
 * @date    2020/12/31 15:54
 */

namespace App\Http\Controllers\App\V1\DistributionPartner;


use App\Http\Controllers\App\V1\BaseController;
use CrsSdk\Work\PartnerMember;

class PartnerIsAcceptSmsController  extends BaseController
{
    /**
     * 合伙人是否接受短信消息
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2020/12/31 15:55
     */
    public function index()
    {
        $params = $this->request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $partnerInfo = $this->partnerDetail($params['enterprise_id']);
        $params['partner_id'] =$partnerInfo['partner_id'];
        $result = (new PartnerMember())->smsSwitch($params);
        return success($result ?? []);
    }

}