<?php


namespace App\Http\Controllers\Service\V1\StaffShare;


use App\Services\StaffSettingService;
use Illuminate\Http\Request;
use App\Http\Controllers\Service\V1\BaseController;

class ChangeOwnerController extends BaseController
{
    public function index(Request $request)
    {

        $input = $this->request->only([
            'mem_uid',
            'parent_user_id',
            'userId',
            'share_user_info',
            'share_parent_user_info',
            'is_staff'
        ]);

        return (new StaffSettingService())->changeOwner($input['parent_user_id'], $input['mem_uid'], $input['share_user_info'], $input['share_parent_user_info'], $input['is_staff']);
    }

}