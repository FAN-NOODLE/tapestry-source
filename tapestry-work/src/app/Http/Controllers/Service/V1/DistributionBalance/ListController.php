<?php
/**
 * Created by :PhpStorm
 * User: zhaoyang
 * Date: 2020/10/16
 * Time: 2:47 下午
 */

namespace App\Http\Controllers\Service\V1\DistributionBalance;

use App\Common\Constant;
use App\Http\Controllers\Service\V1\BaseController;
use App\Services\DistributionBalanceService;

class ListController extends BaseController
{

    // 接受参数
    private $input;

    /**
     * 全员分销 -余额管理
     *
     * @param DistributionBalanceService $distributionBalanceService
     * @return \Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     * @throws \Exception
     */
    public function index(DistributionBalanceService $distributionBalanceService)
    {
        // 参数验证
        $this->getParams();
        // 切换域名
        switchAppHost(['enterprise_id' => $this->input['enterprise_id']]);
        $listObj       = $distributionBalanceService->balanceList($this->input, true);
        $statisticsArr = $distributionBalanceService->balanceStatistics($this->input['enterprise_id']);
        $returnArr = array_merge($statisticsArr, $listObj);
        return success($returnArr);
    }

    /**
     * @return bool
     * @throws \Illuminate\Validation\ValidationException
     */
    private function getParams()
    {

        $rules    = [
            'enterprise_id' => 'required',
            'domain'        => 'required',
        ];
        $messages = [
            'enterprise_id.*' => '请确认企业',
            'domain.*'        => '企业Domain信息不准确',
        ];

        $this->validate($this->request, $rules, $messages);
        $this->input = $this->request->all();
        if (!empty($this->input['worker_ids'])) {
            $this->input['worker_ids'] = (array)$this->input['worker_ids'];
        } else {
            $this->input['worker_ids'] = [];
        }
        $this->input['page'] = $this->input['page'] ?? Constant::PAGE;
        $this->input['size'] = $this->input['size'] ?? 10;
        return true;
    }
}