<?php
/**
 * Created by :PhpStorm
 * User: lvshuaishuai
 * Date: 2020/8/5
 * Time: 10:25 上午
 */

namespace App\Services;


use App\Common\Constant;
use App\Models\DistributionOrderSettlement;
use App\Models\PartnerMember;
use App\Models\ShoppingGuideCashWithdrawal;
use App\Library\Uc\Member;
use Illuminate\Support\Facades\DB;

class DistributionBalanceService extends BaseService
{


    /**
     * 余额列表搜索导购姓名、电话
     */
    CONST SEARCH_NAME = 1;
    CONST SEARCH_MOBILE = 2;

    /**
     * 门店查询
     */
    CONST IS_STORE = 1;


    /**
     * 导购提现相关统计
     *
     * @param $input
     * @return array
     * @throws \Exception
     */
    public function cashWithdrawalInfo($input)
    {
        $sql = ShoppingGuideCashWithdrawal::query()
            ->select(['expect_cash_amount', 'available_cash_amount', 'cash_withdrawal', 'frozen_amount']);
        if ($input['type'] == 2) {
            $workerId = PartnerMember::query()->where([
                'partner_mobile' => $input['mobile']
            ])->value('id');
            $sql->where([
                'user_type' => ShoppingGuideCashWithdrawal::USER_TYPE_PARTNER,
                'worker_id' => $workerId,
                'enterprise_id' => $input['enterprise_id']]);
        } else {
            $workerList = StoreWorkerService::list([
                'enterprise_id' => $input['enterprise_id'],
                'mobile' => $input['mobile'],
                'store_ids' => [$input['store_id'] ?? 0]
            ])['list'];
            //获取搜索的导购id
            if (empty($workerList)) {
                return [];
            }
            $workerId = current($workerList)['id'];
            $sql->where('worker_id', '=', $workerId);
        }

        $res = $sql->first();
        return $res ? $res->toArray() : [
            'expect_cash_amount' => 0,
            'available_cash_amount' => 0,
            'cash_withdrawal' => 0,
            'frozen_amount' => 0
        ];
    }


    /**
     * 详情列表数据返回
     *
     * @param $input
     * @param bool $isPage
     * @return array
     */
    public function balanceDetailList($input, $isPage = true)
    {
        // 查询数据
        $list = DistributionOrderSettlement::query()
            ->where('enterprise_id', '=', $input['enterprise_id'])
            ->where(function ($query) use ($input) {
                $query->where('promotion_worker_id', '=', $input['worker_id'])
                    ->orWhere('service_shopping_guid_id', '=', $input['worker_id']);
            });

        $total = $list->count();
        if ($isPage) {
            $list = $list->forPage($input['page'], $input['size']);
        }

        $list = $list->orderBy('id', 'desc')
            ->get()
            ->toArray();

        //屏蔽推广佣金或者服务佣金不属于此导购的佣金金额
        foreach ($list as $key => &$item) {
            if ($item['promotion_worker_id'] != $input['worker_id']) {
                $item['promotion_commission'] = 0;
            }
            if ($item['service_shopping_guid_id'] != $input['worker_id']) {
                $item['service_commission'] = 0;
            }

        }

        return [
            'list'  => $list ?? [],
            'page'  => $input['page'],
            'size'  => $input['size'],
            'total' => $total ?? 0
        ];
    }


    /**
     * 余额管理统计
     *
     * @param $enterpriseId
     * @return array
     * User: zhaoyang@vchangyi.com
     * Date: 2020/10/16 15:44
     */
    public function balanceStatistics($enterpriseId)
    {
        // 计算总可提现金额和已经提现金额
        $sumPrice = ShoppingGuideCashWithdrawal::query()
            ->select([
                DB::raw('SUM(expect_cash_amount) as expect_cash_amount'),
                DB::raw('SUM(available_cash_amount) as available_cash_amount'),
                DB::raw('SUM(cash_withdrawal) as cash_withdrawal')
            ])
            ->where('enterprise_id', $enterpriseId)
            ->where('user_type',1)
            ->first()
            ->toArray();
        return [
            'expect_cash_amount'    => $sumPrice['expect_cash_amount'],
            'available_cash_amount' => $sumPrice['available_cash_amount'],
            'cash_withdrawal'       => $sumPrice['cash_withdrawal']
        ];
    }

    /**
     * 余额管理列表
     *
     * @param $input
     * @param bool $isPage
     * @return array
     * @throws \Exception
     */
    public function balanceList($input, $isPage = true)
    {
        $input['page'] = $input['page'] ?? Constant::PAGE;
        $input['size'] = $input['size'] ?? 10;

        if (!empty($input['store_ids']) && !is_array($input['store_ids'])) {
            $input['store_ids'] = [$input['store_ids']];
        }

        //查询的导购
        $workerIds = [];
        // 去shop查导购的条件
        $workerCondition = [];
        //导购手机、姓名
        if (!empty($input['search_text']) && !empty($input['search_type'])) {
            if ($input['search_type'] == self::SEARCH_NAME) {
                $workerCondition['name'] = trim($input['search_text']);
            } elseif ($input['search_type'] == self::SEARCH_MOBILE) {
                $workerCondition['mobile'] = trim($input['search_text']);
            }
        }
        if (!empty($input['worker_ids'])) {
            $workerCondition['ids'] = $input['worker_ids'];
        }
        if (!empty($workerCondition)) {
            $workerCondition['enterprise_id'] = $input['enterprise_id'];
            $workerCondition['page']          = 0;
            //获取搜索的导购id
            $workerData = StoreWorkerService::list($workerCondition);
            if (empty($workerData['list'])) {
                return [
                    'list'  => [],
                    'page'  => $input['page'],
                    'size'  => $input['size'],
                    'total' => 0
                ];
            }
            $workerIds = array_column($workerData['list'], 'id');
        }
        // 基础数据查询
        $listObj = ShoppingGuideCashWithdrawal::query()
            ->where('enterprise_id', $input['enterprise_id'])
            ->where('user_type',1)
            ->when(!empty($workerIds), function ($query) use ($workerIds) {
                $query->whereIn('worker_id', $workerIds);
            })
            ->when(!empty($input['start_time']), function ($query) use ($input) {
                return $query->where('created_at', '>=', $input['start_time']);
            })
            ->when(!empty($input['end_time']), function ($query) use ($input) {
                return $query->where('created_at', '<=', $input['end_time']);
            });
        $total = $listObj->count();
        if ($isPage) {
            $listObj = $listObj->forPage($input['page'], $input['size']);
        }

        $listObj = $listObj->orderBy('created_at', 'desc')
            ->get()
            ->toArray();

        if (!empty($listObj)) {
            $memUidMember = array_column($listObj, 'mem_uid');
            // 获取uc员工信息
            $memList      = (new Member())->list(
                $input['domain'],
                'contact',
                [
                    'memUids'  => $memUidMember,
                    'pageSize' => count($memUidMember),
                    'pageNum'  => 1
                ]
            );
            $members      = [];
            $workerIdArr = array_column($listObj, 'worker_id');
            if (!empty($memList)) {
                $members      = array_column($memList['list'], null, 'memUid');
            }
            // 取出workerId
            $workerList = StoreWorkerService::list([
                'enterprise_id' => $input['enterprise_id'],
                'ids' => $workerIdArr
            ], false)['list'];
            $workerList = array_column($workerList, null, 'id');
            //获取部门名，门店名、门店short名
            foreach ($listObj as &$item) {
                $item['worker_no']     = $workerList[$item['worker_id']]['worker_no'] ?? '';
                $item['name']          = $workerList[$item['worker_id']]['name'] ?? '';
                $item['nickname']      = $workerList[$item['worker_id']]['nickname'] ?? '';
                $item['mobile']        = $workerList[$item['worker_id']]['mobile'] ?? '';
                $item['worker_status'] = $workerList[$item['worker_id']]['worker_status'] ?? '';
                $item['type']          = $workerList[$item['worker_id']]['type'] ?? '';
                $item['avatar']        = $workerList[$item['worker_id']]['avatar'] ?? '';
                $item['dpName']        = isset($members[$item['mem_uid']]) ?
                    array_column($members[$item['mem_uid']]['dpName'], 'dpName') : [];
                $item['store_name']    = $item['short_name'] = $workerList[$item['worker_id']]['store_name'] ?? '';
                $item['name'] .= $item['worker_status'] == 2 ? '（已移除）' : '';
                $item['nickname'] .= $item['worker_status'] == 2 ? '（已移除）' : '';
            }
        }

        return [
            'list'  => $listObj ?? [],
            'page'  => $input['page'],
            'size'  => $input['size'],
            'total' => $total ?? 0
        ];

    }

}