<?php

namespace App\Http\Controllers\Service\V1\StaffShare;

use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StaffSettingService;

class UnbindStaffShareOwnController extends BaseController
{
    public function index()
    {
        $userIds = request()->input('user_ids');

        (new StaffSettingService())->unBindingStaffShareByUserId($userIds);

        return success();
    }
}
