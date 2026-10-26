<?php


namespace App\Http\Controllers\Service\V1\StaffShare;


use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StaffShareService;

class CreateStaffShareDataController extends BaseController
{
    public function index()
    {
        $input = $this->request->only([
            'mem_uid',
            'parent_user_id',
            'user_id',
            'enterprise_id',
            'share_user_info',
            'share_parent_user_info',
            'is_staff'
        ]);

        return (new StaffShareService())->createStaffShareData($input['mem_uid'], $input['parent_user_id'], $input['user_id'], $input['enterprise_id'], $input['share_parent_user_info'], $input['share_user_info'], $input['is_staff']);
    }

}