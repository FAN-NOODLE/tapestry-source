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

class RepeatInspectController extends BaseController
{
    /**
     * 分销商品重复校验
     * @param Request $request
     * @return mixed
     * @throws \Exception
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/12/3 3:14 PM
     */
    public function index(Request $request)
    {
        $params = $request->only(['sku_ids']);
        $params['enterprise_id'] = $this->enterpriseId;

        $data = (new DistributionGoods())->repeatInspect($params);

        return success($data);
    }
}
