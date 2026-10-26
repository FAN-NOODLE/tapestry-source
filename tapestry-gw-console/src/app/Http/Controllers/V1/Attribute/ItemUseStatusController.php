<?php

namespace App\Http\Controllers\V1\Attribute;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Attribute;

class ItemUseStatusController extends BaseController
{
    /**
     * @return mixed
     * @throws \Exception
     */
    public function index()
    {
        // 接受参数
        $params = $this->request->only(['attribute_item_id']);

        $params['enterprise_id'] = $this->enterpriseId;

        return success((new Attribute())->itemUseStatus($params));
    }

}
