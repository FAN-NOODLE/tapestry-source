<?php
/**
 * 客户画像年龄统计
 * Created by PhpStorm.
 * User: 何岳龙
 * Date: 2018/11/16
 * Time: 10:14 AM
 */

namespace App\Http\Controllers\V1\CustomerPortrait;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\CustomerPortrait;

class AgeController extends BaseController
{

    public function index()
    {
        $authData = $this->getStores($this->enterpriseAdmin['role']);

        $authType  = $this->isSuperAdmin ? 0 : $authData['type'];
        $store_ids = isset($authData['stores']) ? $authData['stores'] : [];

        return (new CustomerPortrait())->age([
            'enterprise_id' => $this->enterpriseId,
            'auth_type'     => $authType,
            'store_ids'     => $store_ids
        ]);
    }
}