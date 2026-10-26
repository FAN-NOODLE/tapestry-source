<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/27
 * Time: 21:34
 */

namespace App\Http\Controllers\Service\V1\Customer;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\MemberUserRelation;
use Illuminate\Http\Request;

class RemarkController extends BaseController
{
    /**
     * 客户基础信息
     *
     * @param Request $request
     *
     * @return array|mixed
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index(Request $request)
    {
        $this->validate($request, [
            'union_id' => 'required',
        ], [
            'union_id.required' => '客户唯一id不能为空'
        ]);
        // 客户union_id
        $params = $request->only('union_id');

        // 专属导购mem_uid
        $data = MemberUserRelation::where($params)->first(['user_remark_name']);

        return success(['data' => $data]);
    }
}
