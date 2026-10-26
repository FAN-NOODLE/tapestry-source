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

class DetailController extends BaseController
{
    /**
     * 查看佣金设置接口
     *
     * @param Request $request
     *
     * @return \Illuminate\Http\Resources\Json\AnonymousResourceCollection|mixed
     */
    public function index(Request $request)
    {
        // 企业id
        $enterpriseId = $request->get('enterprise_id', 0);

        if (!$enterpriseId) {
            error(600000);
        }

        $data = CommissionSetting::where(['enterprise_id' => $enterpriseId])->get();

        return CommissionSettingResource::collection($data);
    }
}
