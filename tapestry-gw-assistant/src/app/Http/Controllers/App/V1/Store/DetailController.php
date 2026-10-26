<?php
/**
 * Created by PHPStorm.
 * author: 刘朋波
 * Date: 2020/4/10
 * Time: 14:36
 */

namespace App\Http\Controllers\App\V1\Store;

use App\Http\Controllers\App\V1\BaseController as Controller;
use CrsSdk\Shop\Store;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;

/**
 * 门店详情
 * Class DetailController
 * @package App\Http\Controllers\App\V1\Store
 */
class DetailController extends Controller
{
    /**
     * 获取门店详情
     * @param Request $request
     * @param Store $store
     * @return array|\Illuminate\Http\JsonResponse|void
     * @throws \Illuminate\Validation\ValidationException
     */
    public function __invoke(Request $request, Store $store)
    {
        // 参数校验
        $rules = [
            'id' => 'required|integer'
        ];
        $message = [
            'id.required' => '门店ID参数不能为空',
            'id.integer'  => '门店ID参数类型有误',
        ];
        $this->validate($request, $rules, $message);
        $enterpriseId = Auth::user()['enterprise_id'];
        $id = $request->get('id');
        try {
            return success($store->detail([
                'enterprise_id' => $enterpriseId,
                'id'            => $id
            ]));
        } catch (\Exception $e) {
            error(653001);
        }
    }
}
