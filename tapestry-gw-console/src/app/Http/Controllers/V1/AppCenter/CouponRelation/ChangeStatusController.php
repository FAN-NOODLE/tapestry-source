<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/26
 * Time: 10:04
 */

namespace App\Http\Controllers\V1\AppCenter\CouponRelation;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\CouponRelation;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Validation\Rule;
use Illuminate\Validation\ValidationException;

class ChangeStatusController extends BaseController
{

    /**
     * 启用或禁用关联
     * @param Request $request
     * @return JsonResponse
     * @throws ValidationException
     * @throws Exception
     */
    public function index(Request $request)
    {
        // 参数验证规则
        $rules   = [
            'id'     => 'required',
            'status' => ['required', Rule::in([0, 1])],
        ];
        $message = [
            'id.*'     => '关联id有误',
            'status.*' => '关联状态有误',
        ];
        $this->validate($request, $rules, $message);
        $input = $request->only(['id', 'status']);
        $data  = (new CouponRelation())->changeStatus($input);
        return success($data);
    }

}
