<?php


namespace App\Http\Controllers\Service\V1\StaffShare;


use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StaffSettingService;

class HandleBindUserController extends BaseController
{
    public function index()
    {
        $input = $this->request->only([
            'enterprise_id',
            'user_id',
            'is_staff'
        ]);

        return (new StaffSettingService())->handleBindUser($input['user_id'], $input['is_staff'], $input['enterprise_id']);
    }

}