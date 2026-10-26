<?php

namespace App\Http\Controllers\Service\V1\StaffShare;

use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StaffSettingService;

class UnbindStaffShareController extends BaseController
{
    public function index()
    {
        $memUids = request()->input('mem_uids');

        (new StaffSettingService())->unBindingStaffShareBymemUid($memUids);

        return success();
    }
}
