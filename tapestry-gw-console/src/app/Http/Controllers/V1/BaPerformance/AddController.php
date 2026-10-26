<?php

namespace App\Http\Controllers\V1\BaPerformance;

use App\Common\Constant;
use CrsSdk\Work\StoreWorker;
use Exception;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;

class AddController extends BaseController
{
    /**
     * 添加导购
     * @return JsonResponse
     * @throws Exception
     */
    public function index()
    {
        $params                  = $this->request->only(['type', 'member', 'range', 'range_type']);
        $params['enterprise_id'] = $this->enterpriseId;
        $params['operator_id']   = $this->enterpriseAdminId;
        $params['operator_name']   = $this->enterpriseAdmin['username'];

        (new StoreWorker())->addGuide($params);

        try {
            // 更新用户分享记录
            $memUids = array_column($params['member'], 'mem_uid');
            if ($memUids) {
                (new StoreWorker())->staffShareUnbindingChange([
                    'enterprise_id' => $params['enterprise_id'],
                    'mem_uids'      => $memUids,
                    'type'          => Constant::GUIDE_TYPE
                ]);
            }
        } catch (Exception $e) {

        }

        $this->writeOperatorLog('添加导购');

        return success([]);
    }
}
