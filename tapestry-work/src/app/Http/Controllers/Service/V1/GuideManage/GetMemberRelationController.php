<?php

namespace App\Http\Controllers\Service\V1\GuideManage;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\MemberUserRelation;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Validator;

class GetMemberRelationController extends BaseController
{
    public function index(Request $request)
    {

        // 企业id
        $params = $request->input();

        $rules = [
            'enterprise_id' => 'required',
            'mem_uid'       => 'required',
        ];

        $message = [
            'enterprise_id.required' => '企业[enterprise_id]不能为空',
            'mem_uid.required'       => '导购[mem_uid]不能为空',
        ];

        Validator::make($params, $rules, $message)->validate();

        $data = MemberUserRelation::query()
            ->where(['mem_uid' => $params['mem_uid'],'enterprise_id' => $params['enterprise_id']])
//            ->where('user_id','<>',0)
            ->get()->toArray();

        return success($data);
    }
}
