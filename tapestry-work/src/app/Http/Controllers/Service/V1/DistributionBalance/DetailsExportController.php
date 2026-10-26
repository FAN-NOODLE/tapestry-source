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

class DetailsExportController extends BaseController
{
    // 接受参数
    private $input;

    /**
     * 导出数据列表
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        // 参数验证
        $this->getParams();

        // 切换域名
        switchAppHost(['enterprise_id' => $this->input['enterprise_id']]);
        $data = (new DistributionBalanceService())->balanceDetailList($this->input,false);

        return success($data);
    }

    /**
     * 检验参数
     * @return bool
     * @throws \Illuminate\Validation\ValidationException
     */
    private function getParams()
    {
        $rules    = [
            'enterprise_id' => 'required',
            'domain'        => 'required',
            'worker_id'     => 'required'
        ];
        $messages = [
            'enterprise_id.*' => '请确认企业',
            'domain.*'        => '请确认domain',
            'worker_id.*'     => '导购信息不能为空'
        ];

        $this->validate($this->request, $rules, $messages);
        $input = $this->request->all();
        // 接值
        $this->input = $input;
        return true;
    }


}