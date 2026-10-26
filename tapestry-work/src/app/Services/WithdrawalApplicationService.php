<?php
/**
 * 提现申请
 * Created by PhpStorm.
 * User: cjp
 * Date: 2019/5/6
 * Time: 17:46
 */

namespace App\Services;


use App\Models\PartnerMember;
use App\Models\WithdrawalApplication;
use CrsSdk\Boss\Enterprise;
use CrsSdk\Work\Member;

class WithdrawalApplicationService extends BaseService
{

    /**
     * 获取申请列表
     *
     * @param array $input
     * @param bool  $isPage
     *
     * @return array
     * @throws \Exception
     */
    public function getList($input = [], $isPage = true)
    {

        $page    = $input['page'] ?? 1;
        $size    = $input['size'] ?? 10;
        // 获取企业详情
        $enterprise = (new Enterprise())->detail(['enterprise_id' => $input['enterprise_id']])['data'];
        switchAppHost(['enterprise_id' => $input['enterprise_id']]);
        $domain     = $enterprise['domain'];
        $search = [];
        $workIds = $input['worker_ids'] ?? [];
        if(isset($input['search_type']) && !empty($input['search_type'])){
            switch ($input['search_type']) {
                case 1:
                    $input['id'] = $input['search_value'] ?? '';
                    break;
                case 2:
                    $search['name'] = $input['search_value'] ?? '';
                    break;
                case 3:
                    $search['mobile'] = $input['search_value'] ?? '';
                    break;
                default:
                    break;
            }
        }

        if(!empty($search)){
            $search['enterprise_id'] = $input['enterprise_id'];
            $workInfo = StoreWorkerService::list($search)['list'];
            $workIds = array_unique(array_merge(array_column($workInfo,'id'),$workIds));
        }

        //如果是合伙人身份,通过导购ID数组查询出合下面的合伙人
        if ($input['mem_type'] == WithdrawalApplication::APP_TYPE_PARTNER && !empty($workIds)) {
            $input['worker_ids'] = $workIds;
            $partnerList = (new PartnerMemberService())->partnerList($input);
            $partnerList = $partnerList['list'] ?? [];
            $workIds = empty($partnerList) ? [] : array_column($partnerList, 'id');
        }

        $fields = ['withdrawal_application.*'];

        $sql = WithdrawalApplication::query()
            ->when(!empty($workIds), function ($query) use ($workIds) {
                return $query->whereIn('withdrawal_application.worker_id',$workIds);
            })
            ->where('withdrawal_application.enterprise_id', $input['enterprise_id'])
            ->where('app_type', $input['mem_type']) //按类型查询
            ->when(isset($input['id']) && !empty($input['id']), function ($query) use ($input) {
                return $query->where('withdrawal_application.id', 'like', '%' . $input['id'] . '%');
            })
            ->when(!empty($input['ids']), function ($query) use ($input) {
                return $query->whereIn('withdrawal_application.id', $input['ids']);
            })
            ->when(!empty($input['start_time']), function ($query) use ($input) {
                return $query->where('withdrawal_application.request_initiated_at', '>=', $input['start_time']);
            })
            ->when(!empty($input['status']), function ($query) use ($input) {
                return $query->where('withdrawal_application.status', $input['status']);
            })
            ->when(!empty($input['end_time']), function ($query) use ($input) {
                return $query->where('withdrawal_application.request_initiated_at', '<', $input['end_time']);
            });
        //如果是合伙人
        if ($input['mem_type'] == WithdrawalApplication::APP_TYPE_PARTNER) {
            $sql = $sql->leftJoin('partner_member', 'withdrawal_application.worker_id', '=', 'partner_member.id');
            $fields = array_merge($fields, [
                'partner_member.id as partner_id',
                'partner_member.id as partner_sn',
                'partner_member.id as partner_name',
                'partner_member.id as partner_mobile',
                'partner_member.id as partner_avatar'
            ]);
        }

        $total = $sql
            ->count();
        if ($isPage) {
            $sql = $sql->forPage($page, $size);
        }

        $list = $sql
            ->orderBy('withdrawal_application.created_at', 'desc')
            ->get($fields)
            ->toArray();

        if (!empty($list)) {
            $list = $this->formatData($list, $input['enterprise_id'],$domain);
        }

        return [
            'list'  => $list,
            'total' => $total,
            'page'  => (int)$page,
            'size'  => (int)$size
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
     */
    public function formatData($list = [], $enterpriseId = 0,$domain='')
    {

        $data          = [];

        $workerIds = [];
        $partnerIds = [];
        foreach ($list as $info) {
            if ($info['app_type'] == WithdrawalApplication::APP_TYPE_PARTNER) {
                $partnerIds[] = $info['worker_id'];
            } else {
                $workerIds[] = $info['worker_id'];
            }
        }
        //查询导购的信息
        $partnerList = [];
        if (!empty($partnerIds)) {
            $partnerInfo = (new PartnerMemberService())->partnerList([
                'enterprise_id' => $enterpriseId,
                'partner_ids' => array_unique($partnerIds)
            ]);
            $partnerList = $partnerInfo['list'] ?? [];
            $workerIds = array_column($partnerList, 'worker_id');
            $partnerList = array_column($partnerList, null, 'id');
        }

        $search = [
            'enterprise_id' => $enterpriseId,
            'ids' => $workerIds ?? [],
        ];
        $workInfo = StoreWorkerService::list($search)['list'];
        $workerList = array_column($workInfo, null, 'id');

        // 员工搜索条件
        $condition = [
            'qy_domain'      => $domain,
            'app_identifier' => '',
            'params'         => [
                'memUids' => array_column($workInfo, 'mem_uid'),
            ],
            'enterprise_id'  => $enterpriseId
        ];
        $memList = $this->memberList($condition);
        $memList = array_column($memList, null, 'memUid');

        foreach ($list as $value) {
            $workerInfo = $workerList[ $value['worker_id'] ] ?? [];
            $currentPartnerInfo = $partnerList[ $value['worker_id'] ] ?? [];
            if (!empty($currentPartnerInfo) && isset($currentPartnerInfo['worker_id'])) {
                $workerInfo = $workerList[$currentPartnerInfo['worker_id']] ?? [];
            }

            $dpName     = array_column($memList[ $workerInfo['mem_uid']??0 ]['dpName'] ?? [], 'dpName');
            $data[] = [
                'id'                   => $value['id'],
                'application_amount'   => $value['application_amount'],
                'avatar'               => $workerInfo['avatar'] ?? '',
                'worker_id'            => $workerInfo['worker_no'] ?? '',
                'name'                 => $workerInfo['name'] ?? '',
                'mobile'               => $workerInfo['mobile'] ?? '',
                'dp_name'              => implode('；', $dpName) ?? '',
                'store_name'           => $workerInfo['store_name'] ?? '',
                'dp_store'             => empty($workerInfo['store_name']) ? (empty(implode('；', $dpName)) ? '' : implode('；', $dpName)) : $workerInfo['store_name'],
                'request_initiated_at' => $value['request_initiated_at'],
                'request_handle_at'    => $value['request_handle_at'],
                'reason'               => $value['reason'],
                'status'               => $value['status'],
                'partner_id'           => $partnerList[$value['worker_id']]['id'] ?? 0,
                'partner_name'         => $partnerList[$value['worker_id']]['partner_name'] ?? '',
                'partner_mobile'         => $partnerList[$value['worker_id']]['partner_mobile'] ?? '',
                'partner_avatar'         => $partnerList[$value['worker_id']]['partner_avatar'] ?? '',
            ];
        }
        return $data;
    }


    /**
     * 员工列表
     *
     * @param $condition
     *
     * @return array
     */
    public function memberList($condition)
    {
        $page       = 1;
        $maxPages   = 0;
        $memberList = [];
        do {

            $condition['params']['pageNum']  = $page;
            $condition['params']['pageSize'] = 500;

            // 结果处理
            $result     = (new Member())->memberList($condition);

            $memberList = array_merge($memberList, $result['list']);
            $maxPages   = $maxPages ?: ceil($result['total'] / 500);

            // 步程自增
            $page++;
        } while ($page <= $maxPages);

        return $memberList;
    }
}
