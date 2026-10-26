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

class CashWithdrawalInfoController extends BaseController
{
    // 接受参数
    private $input;

    /**
     * 导购收支相关信息
     *
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
        $data = (new DistributionBalanceService())->cashWithdrawalInfo($this->input);

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
            'mobile'     => 'required'
        ];
        $messages = [
            'enterprise_id.*' => '请确认企业',
            'mobile.*'     => '导购电话不能为空'
        ];

        $this->validate($this->request, $rules, $messages);
        $input = $this->request->all();
        // 接值
        $this->input = $input;
        return true;
    }


}