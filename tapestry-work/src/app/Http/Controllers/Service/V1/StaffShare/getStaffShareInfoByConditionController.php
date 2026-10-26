<?php


namespace App\Http\Controllers\Service\V1\StaffShare;


use App\Http\Controllers\Service\V1\BaseController;
use App\Models\StaffShare;

class getStaffShareInfoByConditionController extends BaseController
{
    public function index()
    {
        $condition = request()->input('condition', []);
        $beforeBinding = (new StaffShare())->getStaffShareInfoByCondition($condition, true);
        return success($beforeBinding);
    }

}