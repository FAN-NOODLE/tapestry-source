<?php
/**
 * Created by :PhpStorm
 * User: lvshuaishuai
 * Date: 2020/8/5
 * Time: 12:19 下午
 */

namespace App\Http\Controllers\App\V1\DistributionBalance;

use App\Http\Controllers\App\V1\BaseController;
use CrsSdk\Work\DistributionBalance;
use Illuminate\Http\Request;

class CashWithdrawalInfoController extends BaseController
{

    /**
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $params = $this->request->input();

        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->loginMobile();
        $params['store_id']      = $request->get('store_id', 0);

        $params['type'] = $this->userType;
        $list = (new DistributionBalance())->cashWithdrawalInfo($params);

        return success($list);
    }

}