<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/27
 * Time: 15:29
 */

namespace App\Http\Controllers\Service\V1\Member;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\ContactQrCode;
use App\Models\ContactQrCodeParams;
use App\Services\StaffManagerService;
use Illuminate\Http\Request;

class AddQrCodeController extends BaseController
{
    /**
     * 添加二维码(同步到本地)
     *
     * @param Request $request
     *
     * @return array
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $enterpriseId = $request->post('enterprise_id', '');
        $domain       = $request->post('domain', '');
        $ucUserId     = $request->post('uc_user_id', '');
        $ucMemId      = $request->post('uc_mem_id', '');
        $ucMemName    = $request->post('uc_mem_name', '');
        $sourceType   = $request->post('source_type', '');


        $contactQrCodeData = ContactQrCodeParams::query()
            ->join('contact_qr_code', 'contact_qr_code_params.contact_qr_code_id', '=', 'contact_qr_code.id')
            ->where([
                'contact_qr_code.enterprise_id'        => $enterpriseId,
                'contact_qr_code.source_type'          => ContactQrCode::SOURCE_TYPE_AUTO_GUIDE,
                'contact_qr_code.type'                 => ContactQrCode::CODE_TYPE_ONE,
                'contact_qr_code_params.enterprise_id' => $enterpriseId,
                'contact_qr_code_params.type'          => ContactQrCodeParams::MEMBER,
                'contact_qr_code_params.object'        => $ucMemId,
            ])
            ->first(['contact_qr_code.*']);


        if (empty($contactQrCodeData)) {
            $contactQrCodeData = (new StaffManagerService())->syncContactQrCode($enterpriseId, $domain, $ucUserId, $ucMemId, $ucMemName, $sourceType);
        }

        return success(['data' => $contactQrCodeData]);
    }
}
