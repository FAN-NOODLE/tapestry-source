<?php
/**
 * Created by :PhpStorm
 * User: lvshuaishuai
 * Date: 2020/8/5
 * Time: 12:19 下午
 */

namespace App\Http\Controllers\V1\DistributionBalance;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\DistributionBalance;

class DetailsController extends BaseController
{

    /**
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index()
    {
        $params = $this->request->input();

        $params['enterprise_id'] = $this->enterpriseId;
        $params['domain'] = $this->enterprise['domain'];
        $list = (new DistributionBalance())->balanceDetails($params);

        return success($list);
    }

}