<?php


namespace App\Services;


use App\Models\DistributionOrderSettlement;
use App\Models\PartnerEnterpriseMap;
use App\Models\PartnerMember;
use App\Models\PartnerStatusChange;
use App\Models\PartnerEnterpriseSettings;
use App\Models\ShoppingGuideCashWithdrawal;
use Illuminate\Support\Facades\DB;

class PartnerMemberService extends BaseService
{

    public function partnerList($input = [])
    {
        $page = $input['page'] ?? 1;
        $size = $input['size'] ?? 10;

        $orderBy = $input['orderby'] ?? 'desc';

        $sql = PartnerMember::query()
            ->join('shopping_guide_cash_withdrawal as s',function($join) use($input){
                $join->on('s.worker_id','=','partner_member.id')
                    ->where('s.enterprise_id','=',$input['enterprise_id'])
                    ->where('s.user_type', ShoppingGuideCashWithdrawal::USER_TYPE_PARTNER);
            }, null,null,'left')
            //->leftJoin('shopping_guide_cash_withdrawal as s', 's.worker_id', '=', 'partner_member.id')
            ->leftJoin('partner_enterprise_map as r', 'r.partner_id', '=', 'partner_member.id')
            ->where('r.enterprise_id', $input['enterprise_id']);

        $sql->select(DB::raw(
            "IFNULL(s.available_cash_amount,0) AS available_cash_amount,
                    IFNULL(s.frozen_amount,0) AS frozen_amount,
                    IFNULL(s.expect_cash_amount,0) AS expect_cash_amount,
                    IFNULL(s.cash_withdrawal,0) AS cash_withdrawal,
                    r.worker_id, r.partner_status, partner_member.*, r.user_id,r.created_at as m_created_at")
        );
        if (!empty($input['partner_status']) && $input['partner_status'] == 1) {
            $sql->where('r.partner_status', '=', 1);
        }
        if (!empty($input['partner_status']) && $input['partner_status'] == 2) {
            $sql->where('r.partner_status', '=', 0);
        }
        if (isset($input['worker_ids']) && !empty($input['worker_ids'])) {
            $sql->whereIn('r.worker_id', $input['worker_ids']);
        }
        $sql->when(isset($input['partner_ids']) && !empty($input['partner_ids']), function ($query) use ($input) {
            return $query->whereIn('partner_member.id', $input['partner_ids']);
        });
        if (isset($input['search_type']) && !empty($input['search_type']) && !empty($input['search_value'])) {
            switch ($input['search_type']) {
                case 1:
                    $sql->where('partner_member.partner_name', 'like', $input['search_value'] . '%');
                    break;
                case 2:
                    $sql->where('partner_member.partner_sn', 'like', $input['search_value'] . '%');
                    break;
                default:
                    break;
            }
        }

        if (isset($input['order_by_type']) && !empty($input['order_by_type'])) {
            switch ($input['order_by_type']) {
                case 1:
                    $sql->orderBy('expect_cash_amount', $orderBy);
                    break;
                case 2:
                    $sql->orderBy('available_cash_amount', $orderBy);
                    break;
                case 3:
                    $sql->orderBy('cash_withdrawal', $orderBy);
                    break;
                case 4:
                    $sql->orderBy('frozen_amount', $orderBy);
                    break;
                default:
                    $sql->orderBy('partner_member.created_at', $orderBy);
                    break;
            }
        } else {
            $sql->orderBy('partner_member.created_at', $orderBy);
        }

        if (isset($input['commission']) && is_numeric($input['commission'])) {
            $commission_min = $input['commission_min'] * 100;
            $commission_max = $input['commission_max'] * 100;
            switch ($input['commission']) {
                case 1:
                    $sql->whereRaw('IFNULL(expect_cash_amount,0) >='. $commission_min);
                    $sql->whereRaw('IFNULL(expect_cash_amount,0) <='. $commission_max);
                    break;
                case 2:
                    $sql->whereRaw('IFNULL(available_cash_amount,0) >='. $commission_min);
                    $sql->whereRaw('IFNULL(available_cash_amount,0) <='. $commission_max);
                    break;
                case 3:
                    $sql->whereRaw('IFNULL(cash_withdrawal,0) >='. $commission_min);
                    $sql->whereRaw('IFNULL(cash_withdrawal,0) <='. $commission_max);
                    break;
                default:
                    break;
            }
        }

        $total = $sql->count();

        $sql = $sql->forPage($page, $size);


        $list = $sql->get()
            ->toArray();

        if (!empty($list)) {
            $list = $this->formatData($list, $input['enterprise_id']);
        }

        return [
            'list' => $list,
            'total' => $total,
            'page' => (int)$page,
            'size' => (int)$size
        ];
    }

    /**
     * 格式化代码
     *
     * @param array $list
     * @param array $workInfo
     * @param int $enterpriseId
     * @param string $domain
     * @return array
     * @throws \Exception
     */
    public function formatData($list = [], $enterpriseId = 0)
    {

        $data = [];

        $worker_ids = array_column($list, 'worker_id');
        $search = [
            'enterprise_id' => $enterpriseId,
            'ids' => $worker_ids ?? [],
        ];
        $workInfo = StoreWorkerService::list($search,false)['list'];
        $workerList = array_column($workInfo, null, 'id');

        foreach ($list as $value) {

            $workerInfo = $workerList[$value['worker_id']] ?? [];
            $workName = $workerInfo['name'] ?? '';
            $workName .= !empty($workerInfo['worker_status']) && $workerInfo['worker_status'] == 2 ? '（已移除）' : '';
            $data[] = [
                'id' => $value['id'],
                'partner_sn' => $value['partner_sn'],
                'partner_name' => $value['partner_name'],
                'partner_mobile' => $value['partner_mobile'],
                'worker_id' => $value['worker_id'],
                'worker_avatar' => $workerInfo['avatar'] ?? '',
                'partner_avatar' => $value['avatar'],
                'worker_name' => $workName,
                'for_withdrawal' => $value['expect_cash_amount'],
                'can_withdrawal' => $value['available_cash_amount'],
                'withdrawal_of' => $value['frozen_amount'],
                'have_withdrawal' => $value['cash_withdrawal'],
                'created_at' => $value['m_created_at'],
                'partner_status' => $value['partner_status']
            ];
        }

        return $data;
    }

    public function partnerDetail($input = [])
    {
        if (empty($input['id'])) {
            return [];
        }
        $page = $input['page'] ?? 1;
        $size = $input['size'] ?? 10;

        $partner_id = $input['id'];
        $sql = PartnerMember::query()
            ->leftJoin('partner_enterprise_map as r', 'r.partner_id', '=', 'partner_member.id')
            ->where('partner_member.id', '=', $partner_id)
            ->where('r.enterprise_id','=',$input['enterprise_id'])
            ->select(['partner_member.*','r.created_at as r_created_at', 'r.worker_id', 'r.forbidden_reason', 'r.partner_status', 'r.user_id']);
        $info = $sql->first();
        if (!$info) {
            return [];
        }
        $info = $info->toArray();
        $search = [
            'enterprise_id' => $input['enterprise_id'],
            'ids' => [$info['worker_id']] ?? [],
        ];
        $workInfo = StoreWorkerService::list($search,false)['list'];
        $workInfo =  $workInfo[0] ?? [];
        $workName = $workInfo['name'] ?? '';
        $workName .= !empty($workInfo['worker_status']) && $workInfo['worker_status'] == 2 ? '（已移除）' : '';
        $partnerInfo = [
            'partner_name' => $info['partner_name'],
            'partner_sn' => $info['partner_sn'],
            'partner_id' => $info['id'],
            'partner_mobile' => $info['partner_mobile'],
            'worker_name' => $workName,
            'partner_status' => $info['partner_status'],
            'forbidden_reason' => $info['forbidden_reason'],
            'avatar' => $info['avatar'],
            'worker_avatar' => $workInfo['avatar'] ?? '',
            'created_at' => $info['r_created_at']
        ];
        $sql = DistributionOrderSettlement::query()
            ->where('partner_id', $partner_id)
            ->where('enterprise_id', $input['enterprise_id']);

        $total = $sql->count();

        $sql = $sql->forPage($page, $size);

        $settlement = $sql->get()->toArray();

        $list = [];
        foreach ($settlement as $v) {
            $list[] = [
                'id' => $v['id'],
                'order_no' => $v['order_no'],
                'order_status' => $v['order_status'],
                'settlement_status' => $v['settlement_status'],
                'order_pay_time' => $v['order_pay_time'],
                'partner_commission' => $v['partner_commission'],
            ];
        }
        return [
            'partner' => $partnerInfo,
            'list' => $list,
            'total' => $total,
            'page' => (int)$page,
            'size' => (int)$size
        ];
    }

    /**
     * 禁用合伙人
     *
     * @param array $input
     * @return bool
     */
    public function Switch($input = [])
    {
        $forbidden_reason = $input['forbidden_reason'] ?? '';
        $forbidden_reason = $input['switch'] == 1 ? '' : $forbidden_reason;
        $changInfo = [
            'partner_enterprise_id' => $input['partner_enterprise_id'],
            'partner_status' => $input['switch'],
            'forbidden_reason' => $forbidden_reason,
            'forbidden_code' => $input['switch'] == 1 ? 0 : $input['forbidden_code'],
            'created_at' => time()
        ];
        $changId = PartnerStatusChange::query()->insertGetId($changInfo);
        $model = PartnerEnterpriseMap::query()->where('partner_id', $input['id'])
            ->where('enterprise_id', $input['enterprise_id']);
        $model->update([
            'partner_status' => $input['switch'],
            'forbidden_reason' => $forbidden_reason,
            'forbidden_code' => $changId,
            'updated_at' => time()
        ]);
        return true;
    }

    /**
     * 账户余额异常 禁用合伙人
     *
     * @param $partner_id
     * @param $enterprise_id
     * @return bool
     */
    public function negativeDisable($partner_id, $enterprise_id)
    {
        $partnerInfo = PartnerEnterpriseMap::query()
            ->where(['partner_id' => $partner_id, 'enterprise_id' => $enterprise_id])
            ->first();
        if (empty($partnerInfo)) {
            return false;
        }
        $input = [
            'switch' => 0,
            'forbidden_code' => PartnerStatusChange::FORBIDDEN_CODE['negative'],
            'forbidden_reason' => '因账户余额异常已禁用，请核实后上诉管理人员申请开启',
            'enterprise_id' => $enterprise_id,
            'id' => $partner_id,
            'partner_enterprise_id' => $partnerInfo->id
        ];

        return self::Switch($input);
    }

    public function getPartnerInfo($input)
    {
        $info = PartnerMember::query()
            ->select([
                'partner_member.*',
                'partner_enterprise_map.enterprise_id',
                'partner_enterprise_map.worker_id',
                'partner_enterprise_map.worker_name',
                'partner_enterprise_map.worker_uid',
                'partner_enterprise_map.partner_status',
                'partner_enterprise_map.mobile_sensitive',
                'partner_enterprise_map.mobile_show',
                'partner_enterprise_map.worker_avatar',
                'partner_enterprise_map.user_id',
                'partner_enterprise_map.sms_switch'
            ])
            ->leftJoin('partner_enterprise_map', 'partner_enterprise_map.partner_id', '=', 'partner_member.id')
            ->where('partner_enterprise_map.enterprise_id', $input['enterprise_id'])
            ->where('partner_member.id', $input['id'])
            ->when(isset($input['status']), function ($query) use ($input) {
                return $query->where('partner_enterprise_map.partner_status', $input['status']);
            })
            ->first();
        return empty($info) ? [] : $info->toArray();
    }

    /**
     * 获取二级分销配置
     * User: Hui Yang
     * Email: 2361313833@qq.com
     * DateTime: 2020/11/26 20:39
     * @param $enterpriseId int 企业ID
     * @param int $type int 类型
     * @return array
     */
    public function getPartnerConfig($enterpriseId, $type = PartnerEnterpriseSettings::PARTNER_STATUS_ENABLE)
    {
        $config = PartnerEnterpriseSettings::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('partner_status', $type)
            ->first();
        return empty($config) ? [] : $config->toArray();
    }

    /**
     * 根据用户Id获取合伙人信息
     * @param int $enterpriseId
     * @param int $userId
     * @param int $status
     * @return array
     */
    public function getPartnerInfoByUserId($enterpriseId, $userId, $status)
    {
        if ($enterpriseId && $userId) {
            $partnerInfo = PartnerMember::query()
                ->select(['partner_member.*'])
                ->leftJoin('partner_enterprise_map', 'partner_enterprise_map.partner_id', '=', 'partner_member.id')
                ->where('partner_enterprise_map.enterprise_id', $enterpriseId)
                ->where('partner_enterprise_map.user_id', $userId)
                ->where('partner_enterprise_map.partner_status', $status)
                ->first();

            return !empty($partnerInfo) ? $partnerInfo->toArray() : [];
        }
        return [];
    }

    /**
     * Notes: 根据导购ID获取合伙人ID
     * User: Hui Yang
     * Email: yanghui@vchangyi.com
     * DateTime: 2021/1/27 15:43
     * @param $enterpriseId
     * @param array $workerIds
     * @return array
     */
    public function getPartnerIdsByWorkerIds($enterpriseId, array $workerIds)
    {
        $partnerIds = PartnerMember::query()
            ->leftJoin('partner_enterprise_map as r', 'r.partner_id', '=', 'partner_member.id')
            ->where('r.enterprise_id', $enterpriseId)
            ->whereIn('r.worker_id', $workerIds)
            ->pluck('partner_member.id')
            ->toArray();
        return $partnerIds;
    }
}
