<?php


namespace App\Http\Controllers\Service\V1\StaffShare;


use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StaffShareService;
use Illuminate\Http\Request;

class GetShareParentInfoController extends BaseController
{

    public function index(Request $request)
    {

        $input = $this->request->only([
            'mem_uid',
            'parent_user_id',
            'enterprise_id',
            'user_id',
            'share_user_info',
            'is_staff'
        ]);
        $rules = [
            'enterprise_id' => 'required',
        ];
        $messages = [
            'enterprise_id.required' => '企业id不能为空',
        ];

        $this->validate($this->request, $rules, $messages);


        return (new StaffShareService())->getShareParentInfo($input['mem_uid'], $input['parent_user_id'], $input['user_id'], $input['enterprise_id'], $input['share_user_info'], $input['is_staff']);
    }

}