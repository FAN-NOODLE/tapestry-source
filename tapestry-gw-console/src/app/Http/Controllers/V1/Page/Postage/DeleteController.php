<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2019/4/4
 * Time: 14:11
 */
namespace App\Http\Controllers\V1\Page\Postage;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\GlobalPostage;
use Illuminate\Http\Request;
use \Exception;

class DeleteController extends BaseController
{
    /**
     * 删除
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     * @throws Exception
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

        return success((new GlobalPostage())->delete($params));
    }
}
