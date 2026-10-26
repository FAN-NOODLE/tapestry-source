<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/23
 * Time: 17:32
 */

namespace App\Http\Controllers\Service\V1\Commission;

use App\Http\Controllers\Service\V1\BaseController;
use App\Http\Resources\CommissionSettingResource;
use App\Models\CommissionSetting;
use Illuminate\Http\Request;

class SaveController extends BaseController
{
    /**
     * 佣金设置信息
     *
     * @param Request $request
     *
     * @return mixed
     */
    public function index(Request $request)
    {
        // 企业id
        $enterpriseId = $request->post('enterprise_id', 0);
        if (!$enterpriseId) {
            error(600000);
        }

        // 佣金设置数据
        $commissionSetting = $request->post('commission_setting', []);
        // 条件
        $where = ['enterprise_id' => $enterpriseId];

        if ($commissionSetting) {
            // 更新或创建
            foreach ($commissionSetting as $value) {
                $condition = [
                    'enterprise_id' => $enterpriseId,
                    'type'          => $value['type']
                ];
                app(CommissionSetting::class)->updateOrCreate($condition, $value);
            }
        } else {
            // 删除值
            CommissionSetting::where($where)->update(['value' => 0]);
        }

        $data = CommissionSetting::where($where)->get();

        return CommissionSettingResource::collection($data);
    }
}
