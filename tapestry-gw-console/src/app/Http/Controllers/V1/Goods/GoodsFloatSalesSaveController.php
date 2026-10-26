<?php
/**
 * 商品浮动销量详情
 * GoodsSalesController.php
 * Created on 2019-08-02 15:14
 * Created by chengyujia@vchangyi.com
 */

namespace App\Http\Controllers\V1\Goods;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\GoodsFloatSales;

class GoodsFloatSalesSaveController extends BaseController
{
    /**
     * @return mixed|\Psr\Http\Message\StreamInterface
     */
    public function index()
    {
        $params = $this->request->only(['data', 'source_id']);

        $data = [
            'enterprise_id' => $this->enterpriseId,
            'source_id'     => empty($params['source_id']) ? $this->enterpriseId : $params['source_id'],
            'data'          => $params['data'] ?? []
        ];

        $info = (new GoodsFloatSales())->save($data);

        return success($info);
    }
}
