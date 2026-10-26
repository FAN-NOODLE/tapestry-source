<?php

namespace App\Http\Controllers\Html\V1\MarketingPromotion;

use App\Common\Constant;
use App\Http\Controllers\Html\V1\BaseController;
use CrsSdk\MicroGoodsV1\Activity;

class GoodsListController extends BaseController
{
    /**
     * 获取商品列表
     *
     * @return array
     * @throws \Exception
     */
    public function index()
    {
        $params = request()->only(['keywords', 'store_id', 'page', 'size']);

        $search = [
            'enterprise_id' => $this->enterpriseId,
            'status' => [Constant::STATUS_UP],
            'need_sign_price' => Constant::SIGN_PRICE,
            'goods_type' => Constant::NORMAL_GOODS
        ];
        if (isset($params['keywords']) && $params['keywords'] != '') {
            $search['search_value'] = $params['keywords'];
        }

        $search['source_id'] = $this->enterpriseId;
        if (!empty($params['store_id']) && $this->enterprise['version'] == Constant::NO_CENTER_VERSION) {
            // 官方直营版
            $search['source_id'] = $params['store_id'];
        }
        $search['page'] = $params['page'] ?? 1;
        $search['size'] = $params['size'] ?? 10;

        switchAppHostByDomain($this->domain);
        
        $list = (new Activity())->getSkuList($search);
        foreach ($list['list'] as &$value) {
            switch ($value['sale_type']) {
                case Constant::GOODS_EXPRESS :
                case Constant::GOODS_SELF_EXTRACTING_AND_EXPRESS :
                    $value['price'] = $value['on_line_price'];
                    break;
                case Constant::GOODS_SELF_EXTRACTING :
                    $value['price'] = $value['off_line_price'];
                    break;
                default :
                    $value['price'] = $value['on_line_price'];
                    break;
            }

            unset($value['on_line_price'], $value['off_line_price'], $value['sale_type']);
        }
        unset($value);

        return success([
            'data' => $list['list'] ?? [],
            'total' => $list['total'] ?? 0,
            'page' => $search['page'],
            'size' => $search['size']
        ]);
    }
}
