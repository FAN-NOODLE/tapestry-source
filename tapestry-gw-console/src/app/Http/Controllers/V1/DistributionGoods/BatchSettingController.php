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

class BatchSettingController extends BaseController
{
    /**
     * 分销商品批量设置
     * @param Request $request
     * @return mixed
     * @throws \Exception
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/12/3 3:14 PM
     */
    public function index(Request $request)
    {
        $params = $request->only(['distribution_goods_ids', 'promotion_commission_amount', 'service_commission_amount']);
        $params['enterprise_id'] = $this->enterpriseId;
        (new DistributionGoods())->batchSetting($params);

        return success();
    }
}
