<?php

namespace App\Http\Controllers\Service\V1\Enterprise;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\EnterpriseThirdPartAccount;
use App\Modules\ThirdPart\MiniProgramModule;

class MiniProgramCodeController extends BaseController
{
    public function index()
    {
        $params = $this->request->all();

        if (empty($params['enterprise_id'])) {
            error(201101);
        }

        $account = EnterpriseThirdPartAccount::getMiniProgram($params['enterprise_id']);

        if (empty($account)) {
            error(201108);
        }

        if (empty($account->main_qrcode)) {
            $mainQrcode = MiniProgramModule::mainMiniProgramQrCode($account, "pages/index/index", 2);
        } else {
            $account->refresh();
            $mainQrcode = $account->main_qrcode;
        }

        return success([
            'data' => [
                'main_qrcode'       => $mainQrcode,
                'mini_program_name' => $account->name,
            ]
        ]);
    }
}