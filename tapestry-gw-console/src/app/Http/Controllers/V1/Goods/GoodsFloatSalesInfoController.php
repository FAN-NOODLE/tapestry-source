<?php
/**
 * @Description
 * @Author:         tuzhengping
 * @Date:           2020-03-30 18:23
 * @Version:        1.0
 */

namespace App\Http\Controllers\V1\Goods;

use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\GoodsFloatSales;
use Illuminate\Support\Facades\Log;

class GoodsFloatSalesInfoController extends BaseController
{
    /**
     * 商品浮动销量详情
     *
     * @return mixed|\Psr\Http\Message\StreamInterface
     *
     * @author shizhixiong@vchangyi.com
     *
     * @date 2019-08-02 14:38
     */
    public function index()
    {
        $params = $this->request->all();

        $where = [
            'enterprise_id' => $this->enterpriseId,
            'source_id'     => empty($params['source_id']) ? $this->enterpriseId : $params['source_id']
        ];

        try{
            $res = (new GoodsFloatSales())->getList($where);
        } catch (Error $error){
            Log::info('商品列表获取失败', $error->data);
            notice(100000, [$error->getMessage()], $error->data);
        } catch (\Exception $e) {
            Log::info('商品列表获取失败', [$e->getMessage()]);
            notice(100000, [$e->getMessage()]);
        }

        return success($res);
    }
}
