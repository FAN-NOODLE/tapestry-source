<?php
/**
 * Created by PHPStorm.
 * author: zhaojie
 * Date: 2020/4/10
 * Time: 14:36
 */

namespace App\Http\Controllers\App\V1\Store;

use App\Http\Controllers\App\V1\BaseController as Controller;
use CrsSdk\Shop\Store;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;
use Illuminate\Support\Facades\Log;

/**
 * 门店编辑
 * Class EditController
 * @package App\Http\Controllers\App\V1\Store
 */
class EditController extends Controller
{
    const  ERROR_MESSAGE_ARR = [
        'short_name' => '门店简称不能为空',
        'province'   => '省市区信息不能为空',
        'city'       => '省市区信息不能为空',
        'district'   => '省市区信息不能为空',
        'address'    => '详细地址不能为空',
        'lon_no'     => '经度格式错误',
        'lat_no'     => '纬度格式错误'
    ];
    /**
     * 编辑门店
     * @param Request $request
     * @param Store $store
     * @return array|\Illuminate\Http\JsonResponse|void
     * @throws \Illuminate\Validation\ValidationException
     */
    public function __invoke(Request $request, Store $store)
    {
        // 参数校验
        $rules   = [
            'id' => 'required|integer'
        ];
        $message = [
            'id.required' => '门店ID参数不能为空'
        ];
        $this->validate($request, $rules, $message);
        $input                  = $this->request->all();
        if (isset($input['store_id'])) {
            unset($input['store_id']);
        }
        // 验证输入参数不能为空
        foreach ($input as $key => $val) {
            //门店座机号码单独检测
            if ($key == 'phone') {
                if ($val) {
                    if (!preg_match('/^\d{1,}$/', $val)) {
                        error(652118, '门店座机号码错误');
                    }
                }
                continue;
            }
            if (!$val) {
                error(652118, self::ERROR_MESSAGE_ARR[$key]);
            }
        }
        $input['enterprise_id'] = Auth::user()['enterprise_id'];
        //$input['enterprise_id'] = 2;
        try {
            return success($store->asssiantEdit($input));
        } catch (\Exception $exception) {
            Log::error("修改门店信息出错：{$exception->getMessage()}");
            error($exception->getCode(), $exception->getMessage());
        }
    }
}
