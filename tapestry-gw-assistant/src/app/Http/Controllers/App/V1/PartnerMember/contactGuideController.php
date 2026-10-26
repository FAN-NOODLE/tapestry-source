<?php


namespace App\Http\Controllers\App\V1\PartnerMember;


use App\Http\Controllers\App\V1\BaseController;
use CrsSdk\Work\Partner;
use Illuminate\Http\Request;

class contactGuideController extends BaseController
{
    const DISABLE = 1;//禁用联系

    const CONTACT = 2;//我的导购

    public function index(Request $request)
    {

        $enterprise_id = $this->enterpriseId;
        $partnerInFo = $this->storeWorkerDetail($enterprise_id, 0);

        $partner_id = $partnerInFo['id'];

        $msgType = $request->get('msgType', 2);

        $data = [
            'enterprise_id' => $enterprise_id,
            'partner_id' => $partner_id,
            'page' => 'pagesA/sales-partner-detail/sales-partner-detail?partner_id=' . $partner_id,
            'contactType' => $msgType
        ];
        if ($msgType == self::DISABLE) {
            $data['title'] = '分销合伙人禁用提醒';
        }
        if ($msgType == self::CONTACT) {
            $data['title'] = '分销合伙人联系您通知';
        }
        (new Partner())->partnerContactGuide($data);

        return success([]);

    }

}