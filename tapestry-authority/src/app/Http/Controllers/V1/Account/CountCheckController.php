<?php


namespace App\Http\Controllers\V1\Account;


use App\Http\Controllers\V1\BaseController;
use App\Services\EnterpriseAdminService;
use App\Services\EnterpriseService;
use Illuminate\Http\Request;

/**
 * 判断子账号数量
 * @api {domain}/micro-authority/v1/account/check
 * @author lijian lijian@vchangyi.com 2020/04/13
 * Class CountCheckController
 * @package App\Http\Controllers\V1\Account
 */
class CountCheckController extends BaseController
{
    /**
     * 检查企业子账号数量是否达到上限
     *
     * @api {domain}/micro-authority/v1/account/check
     * @author lijian lijian@vchangyi.com 2020/04/13
     * @return \Illuminate\Http\JsonResponse
     */
    public function index()
    {
        $input = $this->request->only([
            'enterprise_id'
        ]);

        // 验证创建账号规则
        $detail = (new EnterpriseService())->enterpriseDetail($input['enterprise_id']);
        if (empty($detail)) {
            notice(701005); // 企业不存在
        }

        $accountCount = (new EnterpriseAdminService())->getEnterpriseChildAdminCount($input['enterprise_id']);
        if ($detail['account_limit'] <= $accountCount) {
            notice(701004, [$detail['account_limit']]); // 已达到允许创建子账号上限
        }

        return success([]);
    }

}