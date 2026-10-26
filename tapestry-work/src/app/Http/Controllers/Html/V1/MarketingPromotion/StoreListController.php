<?php

namespace App\Http\Controllers\Html\V1\MarketingPromotion;

use App\Common\Constant;
use App\Http\Controllers\Html\V1\BaseController;
use CrsSdk\Shop\Store;
use Illuminate\Support\Facades\Log;

class StoreListController extends BaseController
{
    /**
     * 获取门店列表
     *
     * @return array
     * @throws \Exception
     */
    public function index()
    {
        $search = [
            'enterprise_id' => $this->enterpriseId,
            'mobile' => $this->member['memMobile'],
            'version' => $this->enterprise['version']
        ];

        switchAppHostByDomain($this->domain);

        $list = [];
        if (in_array($this->enterprise['version'], [Constant::CENTER_STORE_VERSION, Constant::NO_CENTER_VERSION])) {
            $list = (new Store())->listByMobile($search);
        }

        return success([
            'data' => [
                'version' => $this->enterprise['version'],
                'list' => $list
            ]
        ]);
    }
}
