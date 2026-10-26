<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/27
 * Time: 15:29
 */

namespace App\Http\Controllers\Service\V1\GuideManage;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\MemberGuide;
use App\Models\MemberUserRelation;
use App\Services\MemberUserRelationService;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Validator;

class DeleteController extends BaseController
{
    /**
     * 删除导购
     *
     * @param Request                   $request
     * @param MemberUserRelationService $memberUserRelationService
     *
     * @return array
     */
    public function index(Request $request, MemberUserRelationService $memberUserRelationService)
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

        $memberUserRelationService->deleteGuide($params);

        $data = ['data' => []];

        return success($data);
    }
}
