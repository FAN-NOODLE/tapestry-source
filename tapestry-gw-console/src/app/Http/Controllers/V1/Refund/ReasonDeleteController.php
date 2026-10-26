<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2019/4/4
 * Time: 14:11
 */
namespace App\Http\Controllers\V1\Refund;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\Refund;
use Illuminate\Http\Request;

class ReasonDeleteController extends BaseController
{
    /**
     * 退款原因保存
     *
     * @param Request $request
     * @return array
     */
    public function index(Request $request)
    {
        $params = $request->only(['id']);

        $rules = [
            'id' => 'required|integer|min:1'
        ];

        $message = [
            'id.required' => '原因ID不存在',
            'id.integer' => '原因ID数据类型有误',
            'id.min' => '原因ID数据类型有误'
        ];

        $this->validate(request(), $rules, $message);

        return success((new Refund())->reasonDelete($params));
    }
}