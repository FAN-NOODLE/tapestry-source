<?php
/**
 * 客户画像区域分页统计
 * Created by PhpStorm.
 * User: 何岳龙
 * Date: 2018/11/16
 * Time: 10:14 AM
 */

namespace App\Http\Controllers\V1\CustomerPortrait;

use App\Common\Constant;
use App\Http\Controllers\V1\BaseController;
use Illuminate\Http\Request;
use CrsSdk\Shop\CustomerPortrait;

class ProvincesPageController extends BaseController
{
    public function index(Request $request)
    {
        $type = $request->post('type', '');

        // 区域类型不能为空
        if (empty($type)) {
            notice(900001);
        }

        // 区域类型错误
        if (!in_array($type, [Constant::PROVINCE_TYPE, Constant::CITY_TYPE])) {
            notice(900002);
        }

        $authData = $this->getStores($this->enterpriseAdmin['role']);

        $authType = $this->isSuperAdmin ? 0 : $authData['type'];

        $store_ids = isset($authData['stores']) ? $authData['stores'] : [];

        return (new CustomerPortrait())->provincesPage([
            'enterprise_id' => $this->enterpriseId,
            'type' => $type,
            'page' => $request->post('page', Constant::DEFAULT_PAGE),
            'size' => $request->post('size', Constant::DEFAULT_SIZE),
            'auth_type'     => $authType,
            'store_ids'     => $store_ids
        ]);
    }

}
