<?php
/**
 * Created by PhpStorm.
 * User: weihua
 * Email: 1367856219@qq.com
 * Time: 9:51
 */

namespace App\Http\Controllers\V1\DistributionBalance;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\DistributionBalance;

class ListController extends BaseController
{
    /**
     * 余额列表
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index()
    {
        $params = $this->request->input();

        $params['enterprise_id'] = $this->enterpriseId;
        $params['domain'] = $this->enterprise['domain'];
        $list = (new DistributionBalance())->balanceList($params);

        return success($list);
    }
}
