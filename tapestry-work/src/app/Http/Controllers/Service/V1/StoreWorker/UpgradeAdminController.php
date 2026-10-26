<?php


namespace App\Http\Controllers\Service\V1\StoreWorker;


use App\Http\Controllers\Service\V1\BaseController;
use App\Services\ScrmMemberService;

class UpgradeAdminController extends BaseController
{
    public function index()
    {
        $params = $this->request->input();

        (new ScrmMemberService())->upgradeMemberRange($params);

        return success();
    }
}