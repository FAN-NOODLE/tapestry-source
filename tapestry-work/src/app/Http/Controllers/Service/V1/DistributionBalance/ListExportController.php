<?php
/**
 * Created by :PhpStorm
 * User: zhaoyang
 * Date: 2020/10/16
 * Time: 2:47 下午
 */

namespace App\Http\Controllers\Service\V1\DistributionBalance;

use App\Http\Controllers\Service\V1\BaseController;
use App\Services\DistributionBalanceService;

class ListExportController extends BaseController
{

    private $input;

    /**
     * @return \Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     * @throws \Exception
     */
    public function index()
    {
        // 参数验证
        $this->getParams();
        // 切换域名
        switchAppHost(['enterprise_id' => $this->input['enterprise_id']]);
        $data = (new DistributionBalanceService())->balanceList($this->input, false)['list'];

        return success($data);

    }

    /**
     * 检验参数
     * @return bool
     * @throws \Illuminate\Validation\ValidationException
     */
    private function getParams()
    {
        // 接值
        $this->input = $this->request->input();

        $rules    = [
            'enterprise_id' => 'required',
        ];
        $messages = [
            'enterprise_id.*' => '请确认企业',
        ];

        $this->validate($this->request, $rules, $messages);

        return true;
    }

}