<?php
/**
 * 提现认证 - 详情
 * Created by PhpStorm.
 * User: Yang Hui
 * Date: 2020年10月16日16:56:4
 */

namespace App\Http\Controllers\Service\V1\Withdrawal;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\Enterprise;
use App\Models\DistributionGoods;
use App\Models\StoreWorkerCertification;
use Illuminate\Http\Request;
use App\Services\WithdrawalServices;

class CertificationDetailController extends BaseController
{
    /**
     * 提现认证 - 详情
     * @param Request $request
     * @return mixed
     * @throws \Exception
     * @author   Yang Hui
     * @date     2020年10月16日17:1:28
     */
    public function index(Request $request)
    {
        $params = $request->only([
            'worker_id',
            'enterprise_id',
            'mem_type'
        ]);

        $certification = StoreWorkerCertification::query()
            ->where('worker_id', $params['worker_id'])
            ->where('enterprise_id', $params['enterprise_id'])
            ->where('mem_type', $params['mem_type'] ?? StoreWorkerCertification::MEM_TYPE_WORKER)
            ->first();
        $data = empty($certification) ? [] : $certification->toArray();
        $returnData = [
            'data' => $data
        ];
        return success($returnData);
    }
    /**
     * 提现认证 - 详情
     * @param Request $request
     * @return mixed
     * @throws \Exception
     * @author   Yang Hui
     * @date     2020年10月16日17:1:28
     */
    public function certificationFileUrl(Request $request)
    {
        $params = $request->only([
            'worker_id',
            'enterprise_id',
            'mem_type'
        ]);
        $url = (new WithdrawalServices())->getUserCertificationFileUrl($params);

        $returnData = [
            'data' => [
                'url' => $url
            ]
        ];
        return success($returnData);
    }
}
