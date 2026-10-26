<?php
/**
 * @author  zhaopengfei<zhaopengfei@vchangyi.com>
 * @date    2020/11/19 22:23
 */

namespace App\Http\Controllers\V1\DistributionConfig;


use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\DistributionConfig;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class PartnerSaveController extends BaseController
{
    /**
     * 分销配置存储
     *
     * @param Request $request
     * @param WorkerDistribution $workerDistribution
     * @return \Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index(DistributionConfig $distribution)
    {
        $rules = [
            'partner_status'   => 'in:0,1',
            'commission_type'  => 'in:1,2',
            'commission_rate' => 'integer',
            'withdraw_status'  => 'in:0,1',
        ];

       /* $messages               = [
            'partner_status.*'   => '是否开启分销不合法',
            'commission_type.*'  => '分账方式不合法',
            'commission_rate.*' => '佣金比例不合法',
            'withdraw_status.*'          => '是否可提现不合法',
        ];
        $input                  = $this->validate($request, $rules, $messages);*/
        $input = $this->request->all();
        $input['enterprise_id'] = $this->enterpriseId;
        $input['opt_admin_id'] = $this->enterpriseAdminId;
        $input['opt_admin_name'] = $this->enterpriseAdmin['username'];
        $distribution->partnerSave($input);

        return success();
    }
}