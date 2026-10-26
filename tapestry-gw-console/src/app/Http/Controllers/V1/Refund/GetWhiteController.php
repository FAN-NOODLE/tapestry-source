<?php
/**
 * 售后列表获取可以显示'退款'按钮的企业ID
 * Created by PhpStorm.
 * User: cjp
 * Date: 2019/6/13
 * Time: 14:32
 */
namespace App\Http\Controllers\V1\Refund;

use App\Http\Controllers\V1\BaseController;
use App\Models\Setting;

class GetWhiteController extends BaseController
{
    public function index(Setting $setting)
    {
        $enterprise_id = $setting->query()
            ->where(['package' => 'refund_enterprise_white'])
            ->first(['config']);
        $enter_arr = json_decode($enterprise_id->config, true);

        return success([
            'data' => ['enterprise_id' => $enter_arr]
        ]);
    }
}
