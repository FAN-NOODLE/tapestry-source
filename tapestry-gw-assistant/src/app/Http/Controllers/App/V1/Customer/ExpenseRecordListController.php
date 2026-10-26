<?php
/**
 * Created by PhpStorm.
 * User: baizhongyang
 * Date: 2020/8/3
 * Time: 10:09
 */

namespace App\Http\Controllers\App\V1\Customer;

use CrsSdk\Work\Customer;
use Illuminate\Http\Request;
use App\Http\Controllers\App\V1\BaseController;

class ExpenseRecordListController extends BaseController
{

    /**
     * 客户消费记录列表
     *
     * @param Request $request
     *
     * @return \Illuminate\Http\JsonResponse
     * @author  baizhongyang <baizhongyang@vchangyi.com>
     * @date    2020/8/3 16:09
     */
    public function index(Request $request)
    {
        $params                  = $request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['qy_domain']     = $this->loginAdminEnterprise['domain'];
        $params['qy_hash']       = $this->loginAdminEnterprise['hash'];

        $data = (new Customer())->customerExpenseRecordList($params);

        return success($data);
    }
}