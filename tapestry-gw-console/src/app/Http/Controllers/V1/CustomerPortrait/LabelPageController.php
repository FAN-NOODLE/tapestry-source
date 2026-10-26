<?php
/**
 * 客户画像分页标签统计
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

class LabelPageController extends BaseController
{
    public function index(Request $request)
    {
        $authData = $this->getStores($this->enterpriseAdmin['role']);

        $authType = $this->isSuperAdmin ? 0 : $authData['type'];

        $store_ids = isset($authData['stores']) ? $authData['stores'] : [];

        return (new CustomerPortrait())->labelPage([
            'enterprise_id' => $this->enterpriseId,
            'type' => $request->post('type', ''),
            'page' => $request->post('page', Constant::DEFAULT_PAGE),
            'size' => $request->post('size', Constant::DEFAULT_SIZE),
            'auth_type'     => $authType,
            'store_ids'     => $store_ids
        ]);
    }

}
