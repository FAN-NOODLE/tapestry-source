<?php


namespace App\Http\Controllers\V1\Account;


use App\Http\Controllers\V1\BaseController;
use App\Library\Util\Logger;
use App\Services\EnterpriseAdminService;

class ForbidLeaveAdmin extends BaseController
{
    public function index()
    {
        EnterpriseAdminService::forbidLeaveAdmin();
        return success([]);
    }
}