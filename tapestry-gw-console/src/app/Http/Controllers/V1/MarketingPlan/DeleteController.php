<?php
/**
 * @author  zhaopengfei<zhaopengfei@vchangyi.com>
 * @date    2021/1/11 11:32
 */

namespace App\Http\Controllers\V1\MarketingPlan;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\MarketingPlan;

class DeleteController extends BaseController
{
    public function index(MarketingPlan $marketingPlan)
    {
        $input = $this->request->all();
        $input['enterprise_id'] = $this->enterpriseId;

        $result = $marketingPlan->detele($input);

        return success($result);
    }
}