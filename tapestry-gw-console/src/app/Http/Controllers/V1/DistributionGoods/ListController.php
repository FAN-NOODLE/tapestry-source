<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/12/3
 * Time: 3:12 PM
 */

namespace App\Http\Controllers\V1\DistributionGoods;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\DistributionGoods;
use Illuminate\Http\Request;

class ListController extends BaseController
{
    /**
     * 分销商品列表
     * @param Request $request
     * @return mixed
     * @throws \Exception
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/12/3 3:14 PM
     */
    public function index(Request $request)
    {
        $params = $request->only(['spu_name','sku_id','page','limit']);
        if(!empty($params['sku_id'])){
            $this->validate($request, [
                'sku_id' => 'integer',
            ], [
                'sku_id.*' => 'sku_id必须是整数类型',
            ]);
        }
        $params['enterprise_id'] = $this->enterpriseId;

        $data = (new DistributionGoods())->goodsList($params);

        return success($data);
    }
}
