<?php


namespace App\Http\Controllers\V1\ActivityCommon\DiscountCard;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroActivity\DiscountCard;

class OpenCardListController extends BaseController
{
    public function index()
    {
        $params = $this->request->only([
            'rule_id',
            'search',
            'page',
            'size',
        ]);

        $params['enterprise_id'] = $this->enterpriseId;
        $params['export'] = 2; // 1-导出 2-列表

        $openCardList = (new DiscountCard())->openCardList($params);

        return success($openCardList);
    }

}