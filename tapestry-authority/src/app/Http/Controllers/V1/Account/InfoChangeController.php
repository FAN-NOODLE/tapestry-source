<?php

namespace App\Http\Controllers\V1\Account;

use App\Models\EnterpriseAdmin;
use App\Http\Controllers\V1\BaseController;

class InfoChangeController extends BaseController
{
    public function infoChange()
    {
        $mobile = $this->request->post('mobile');
        $newMobile = $this->request->post('new_mobile');
        $enterpriseId = $this->request->post('enterprise_id');

        if (empty($mobile)) {
            error(701012);
        }

        if (empty($newMobile)) {
            error(701013);
        }

        if (empty($enterpriseId)) {
            error(700014);
        }

        if ($mobile == $newMobile) {
            return success();
        }

        $admin = EnterpriseAdmin::query()
            ->where('enterprise_id', '=', $enterpriseId)
            ->where('mobile', '=', $newMobile)
            ->first();

        if ($admin) {
            error(700026, [], [], 200);
        }

        EnterpriseAdmin::query()
            ->where('enterprise_id', '=', $enterpriseId)
            ->where('mobile', '=', $mobile)
            ->update([
                'mobile'  => $newMobile,
                'account' => '+86-' . $newMobile,
            ]);

        return success();
    }
}