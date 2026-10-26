<?php
/**
 * Created by PhpStorm.
 * User: liupengwei
 * Date: 2020-03-16
 * Time: 18:52
 */

namespace App\Http\Controllers\V1\GoodsLabel;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\GoodsLabel;
use Illuminate\Support\Facades\Log;

class ListController extends BaseController
{

    /**
     * 商品标签列表
     *{domain}/app/v1/goods-label/list
     *
     * @return mixed
     * @throws \Exception
     */
    public function index()
    {
        // 接收参数
        $params = request()->only(['page', 'size', 'search']);
        // 企业id
        $params['params']['enterprise_id'] = $this->enterpriseId;

        $list = [];
        try {
            return (new GoodsLabel())->goodsLabels($params);

        } catch (\Exception $exception) {
            Log::info('商品标签列表数据获取失败' . $exception->getMessage(), $params);
            notice(920202);
        }

        return success(['data' => $list]);
    }
}
