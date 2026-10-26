<?php


namespace App\Http\Controllers\V1\ActivityCommon\DiscountCard;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroActivity\DiscountCard;

class CardStatusController extends BaseController
{
    /**
     * 修改先享卡卡片状态 暂停开卡和恢复开卡
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index()
    {
        $params = $this->request->only([
            'rule_id',
            'id',
            'status',
        ]);

        $params['enterprise_id'] = $this->enterpriseId;

        (new DiscountCard())->cardStatus($params);

        return success([]);
    }

}