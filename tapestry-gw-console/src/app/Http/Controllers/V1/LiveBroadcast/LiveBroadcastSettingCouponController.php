<?php
/**
 * Project  console
 * FileName LiveBroadcastSettingCouponController.php
 * User     吉王伟
 * Date     2020/4/9 15:43
 */

namespace App\Http\Controllers\V1\LiveBroadcast;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\LiveBroadcast;
use Illuminate\Support\Facades\Validator;

class LiveBroadcastSettingCouponController extends BaseController
{
    /**
     * @Api            ：app/v1/live-broadcast/get_coupon_list
     * @ApiDescription ：企业后台直播设置发放礼券
     * @ApiGroup       ：企业后台直播
     * @User           ：吉王伟
     * @Date           ：2020/4/9
     *
     * @return mixed
     * @throws \Exception
     */
    public function getSettingCouponList()
    {
        // 查询参数处理
        $params = request()->only(['keywords', 'page', 'size']);

        $rules = [
            'page' => 'sometimes|required|integer',
            'size' => 'sometimes|required|integer',
        ];

        $message = [
            'page.required' => '页码不能为空',
            'page.integer'  => '页码类型错误',
            'size.required' => '每页数据条数不能为空',
            'size.integer'  => '每页数据条数类型错误',
        ];

        // 参数验证
        Validator::make($params, $rules, $message)->validate();

        $enterpriseId = $this->enterpriseId;

        $where = [
            'keywords'              => isset($params['keywords']) ? trim($params['keywords']) : '',
            'enterprise_id'         => $enterpriseId,
            'channel'               => 2,
            'coupon_setting_status' => 1,
            'page'                  => $params['page'] ?? 1,
            'size'                  => $params['size'] ?? 10
        ];

        $return = (new LiveBroadcast())->getLiveBroadcastCouponSettingList($where);
        return success($return);
    }
}