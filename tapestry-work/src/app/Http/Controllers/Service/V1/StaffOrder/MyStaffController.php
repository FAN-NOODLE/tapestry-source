<?php
/**
 * Created by PhpStorm.
 * Function: 企业员工业绩排行
 * User: 龚国新
 * Date: 2020/2/14
 * Time: 11:43
 */

namespace App\Http\Controllers\Service\V1\StaffOrder;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StaffOrderService;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Validator;

class MyStaffController extends BaseController
{

    //   protected $requireLogin = false;
    /**
     * 我的业绩数字
     *
     * @param Request $request
     *
     * @return array
     * @throws \Exception
     */
    public function index(Request $request)
    {

        $isAll  = $request->input('is_all', 0);
        $search = $request->input('search', []);

        $validator = $this->validatorData($request->all());
        if ($validator->fails()) {
            throw new Error(400000, $validator->errors()->first());
        }
        $data = (new StaffOrderService())->getMyStaff($search, $isAll);

        return success(['data' => $data]);
    }

    /**
     * @param $params
     *
     * @return \Illuminate\Contracts\Validation\Validator
     */
    public function validatorData($params)
    {

        return Validator::make($params, [
            'search.enterprise_id'  => 'required',
            'search.mem_uid'        => 'required',
            'search.qy_domain'      => 'required',
            'search.app_identifier' => 'required',
        ], [
            'search.enterprise_id.required' => '企业ID不能为空',
            'search.mem_uid.required'       => '用户ID不能为空',
            'search.qy_domain'              => '企业domain不能为空',
            'search.app_identifier'         => '企业app_identifier不能为空',
        ]);

    }
}