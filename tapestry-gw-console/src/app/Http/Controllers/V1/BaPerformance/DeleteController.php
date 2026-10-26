<?php

namespace App\Http\Controllers\V1\BaPerformance;

use App\Common\Constant;
use CrsSdk\Boss\EnterpriseAdmin;
use CrsSdk\Work\StoreWorker;
use CrsSdk\Work\GuideManage;
use Exception;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;
use Overtrue\Pinyin\Pinyin;

class DeleteController extends BaseController
{
    /**
     * 删除导购
     * @return JsonResponse
     * @throws Exception
     */
    public function index()
    {
        $params                  = $this->request->only(['id', 'mem_uid']);
        $params['enterprise_id'] = $this->enterpriseId;
        $params['domain']        = $this->qyDomain;
        $params['operator_id']   = $this->enterpriseAdminId;
        //同步删除后台子账户管理员信息
        $guideParams                  = $this->request->only(['account', 'username']);
        $account                      = $guideParams['account'] ?? '';
        $username                     = $guideParams['username'] ?? '';
        $guideParams['enterprise_id'] = $this->enterpriseId;
        (new GuideManage())->delete($params);
        (new StoreWorker())->deleteGuide($params);
        if (!empty($account)) {
            $createDefaultAccount = function ($username, $mobile) {
                $account = join('', (new Pinyin())->name($username));
                $account = strtolower(str_replace(' ', '', $account));

                $account = $account . substr($mobile, -4);
                return $account;
            };
            $guideParams['account'] = $createDefaultAccount($username, $account);
            (new EnterpriseAdmin())->delete($guideParams);
        }

        try {
            // 更新用户分享记录
            if (!empty($params['mem_uid'])) {
                (new StoreWorker())->staffShareUnbindingLeave([
                    'enterprise_id' => $params['enterprise_id'],
                    'mem_uids'      => [$params['mem_uid']],
                    'type'          => Constant::GUIDE_TYPE
                ]);
            }
        } catch (Exception $e) {

        }
        $this->writeOperatorLog('删除导购');

        return success();
    }
}
