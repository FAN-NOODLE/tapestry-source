<?php
/**
 * @author  zhaopengfei<zhaopengfei@vchangyi.com>
 * @date    2020/12/31 15:51
 */

namespace App\Http\Controllers\App\V1\DistributionPartner;


use App\Http\Controllers\App\V1\BaseController;
use CrsSdk\Work\PartnerMember;

class MobileSensitiveController  extends BaseController
{
    /**
     * 手机号是否脱敏
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2020/12/31 15:51
     */
    public function index()
    {
        $params = $this->request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $partnerInfo = $this->partnerDetail($params['enterprise_id']);
        $params['partner_id'] =$partnerInfo['partner_id'];
        $params['mobile'] = $this->mobile;
        
        $result = (new PartnerMember())->mobileSwitch($params);
        return success($result ?? []);
    }

}