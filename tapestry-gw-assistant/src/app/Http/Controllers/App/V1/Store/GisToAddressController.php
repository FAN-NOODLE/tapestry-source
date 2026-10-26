<?php
/**
 * Created by PHPStorm.
 * author: 刘朋波
 * Date: 2020/4/10
 * Time: 14:36
 */

namespace App\Http\Controllers\App\V1\Store;

use App\Http\Controllers\App\V1\BaseController as Controller;
use App\Library\TencentLbs;
use CrsSdk\Shop\Store;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;
use Illuminate\Support\Facades\Log;

/**
 * 根据经纬度返回省市区
 * Class DetailController
 * @package App\Http\Controllers\App\V1\Store
 */
class GisToAddressController extends Controller
{
    /**
     * 根据经纬度返回省市区
     * @param Request $request
     * @param Store $store
     * @return array|\Illuminate\Http\JsonResponse|void
     * @throws \Illuminate\Validation\ValidationException
     */
    public function __invoke(Request $request, Store $store)
    {
        // 参数校验
        $rules = [
            'lon_no' => 'required',
            'lat_no' => 'required'
        ];
        $message = [
            'lon_no.required' => '经度不能为空',
            'lat_no.required'  => '纬度不能为空',
        ];
        $input = $this->validate($request, $rules, $message);

        try {

            $result = (new TencentLbs())->getAddressByLocation($input['lon_no'],$input['lat_no']);
            if (!isset($result['address_component'])) {
                error(653006);
            }
            $result = $result['address_component'];
            if (!empty($result)) {
                return success([
                    'province' => $result['province'],
                    'city'     => $result['city'],
                    'district' => $result['district'],
                ]);
            }
        } catch (\Exception $e) {
            Log::info('逆地址解析失败：' . $e);
            error(653006);
        }
        return success([]);
    }
}
