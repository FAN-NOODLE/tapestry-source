<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/12/3
 * Time: 7:15 PM
 */

namespace App\Http\Controllers\Service\V1\Enterprise;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\Enterprise;
use App\Models\EnterpriseWorkSetting;
use Illuminate\Http\Request;

class WorkSettingController extends BaseController
{
    /**
     * 保存企业填写的外部联系人的接口调用凭证信息
     *
     * @param Request $request
     *
     * @return array
     * @author 陈朔 chenshuo@vchangyi.com
     * @date   2018/12/3 7:21 PM
     */
    public function index(Request $request)
    {
        $enterpriseId = $request->post('enterprise_id');
        $externalSecret = $request->post('external_secret');

        Enterprise::query()
            ->where('id', $enterpriseId)
            ->update(['external_secret' => $externalSecret]);

        return success(['data' => []]);
    }
}
