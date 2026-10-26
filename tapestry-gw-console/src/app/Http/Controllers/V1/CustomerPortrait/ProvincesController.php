<?php
/**
 * 客户画像区域统计
 * Created by PhpStorm.
 * User: 何岳龙
 * Date: 2018/11/16
 * Time: 10:14 AM
 */

namespace App\Http\Controllers\V1\CustomerPortrait;

use App\Common\Constant;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\CustomerPortrait;
use Illuminate\Http\Request;

class ProvincesController extends BaseController
{
    public function index(Request $request)
    {
        $params = $request->only(['type']);

        // 区域类型不能为空
        if (empty($params['type'])) {
            notice(900001);
        }

        // 区域类型错误
        if (!in_array($params['type'], [Constant::PROVINCE_TYPE, Constant::CITY_TYPE])) {
            notice(900002);
        }

        $authData = $this->getStores($this->enterpriseAdmin['role']);

        $authType = $this->isSuperAdmin ? 0 : $authData['type'];

        $store_ids = isset($authData['stores']) ? $authData['stores'] : [];

        return (new CustomerPortrait())->provincesNotPage([
            'enterprise_id' => $this->enterpriseId,
            'type' => $params['type'],
            'auth_type'     => $authType,
            'store_ids'     => $store_ids
        ]);
    }

}
