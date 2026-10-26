<?php
/**
 * 提现 - 申请填写页面数据
 * Created by PhpStorm.
 * User: Yang Hui
 * Date: 2020年10月16日16:56:4
 */

namespace App\Http\Controllers\Service\V1\Withdrawal;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\PartnerMember;
use App\Models\ShoppingGuideCashWithdrawal;
use App\Models\StoreWorkerCertification;
use Illuminate\Http\Request;
use App\Services\WithdrawalServices;
use App\Models\WithdrawalSetting;

class WithdrawalHomeController extends BaseController
{
    /**
     * 提现 - 申请填写页面数据
     * @param Request $request
     * @return mixed
     * @throws \Exception
     * @author   Yang Hui
     * @date     2020年10月16日17:1:28
     */
    public function index(Request $request)
    {
        $params = $request->only([
            'enterprise_id',
            'worker_id'
        ]);
        $params['mem_type'] = $request->get('mem_type', StoreWorkerCertification::MEM_TYPE_WORKER);

        switchAppHost(['enterprise_id' => $params['enterprise_id']]);
        $data = [];
        // 查询用户是否已经认证，返回认证信息
        $memberInfo = [];
        try {
            //查询当前的导购/合伙人信息
            $memberInfo = (new WithdrawalServices())->getCurrentWorkerInfo($params);
            if (!empty($memberInfo)) {
                $params['mobile'] = $memberInfo['mobile'] ?? $memberInfo['partner_mobile'];
            }
        } catch (\Exception $e) {
            //
            $memberInfo = [];
        }
        try {
            $data['info'] = (new WithdrawalServices())->getUserCertificationInfo($params);
        } catch (\Exception $e) {
            $data['info'] = [];
        }
        //查询配置信息，验证导购是否可以提现
        $withdrawalConfig = WithdrawalSetting::query()
            ->where('enterprise_id', $params['enterprise_id'])
            ->first();
        $allowWithdrawal = false;
        if (isset($memberInfo['type']) && isset($withdrawalConfig->type) && !empty($withdrawalConfig->type)) {
            $tmpArr = explode(',', $withdrawalConfig->type);
            if (in_array($memberInfo['type'], $tmpArr)) {
                $allowWithdrawal = true;
            }
        }
        //用户是否可提现
        $data['allow_withdrawal'] = boolval($allowWithdrawal);

        $data['need_certification'] = true;
        if (!empty($data['info'])) {
            //查询用户可提现的金额
            $balanceInfo = ShoppingGuideCashWithdrawal::query()
                ->where('worker_id', $params['worker_id'])
                ->where('enterprise_id', $params['enterprise_id'])
                ->where('user_type', $params['mem_type'])
                ->select([
                    'available_cash_amount'
                ])
                ->first();
            $data['info']['available_cash_amount'] = $balanceInfo->available_cash_amount ?? 0;
            $data['need_certification'] = false;
        }

        //第三方平台的项目编号
        $data['transfer_no'] = '';
        if (isset($withdrawalConfig->transfer_no) && !empty($withdrawalConfig->transfer_no)) {
            $data['transfer_no'] = md5($withdrawalConfig->transfer_no);
        }

        $returnData = [
            'data' => $data
        ];
        return success($returnData);
    }
}
