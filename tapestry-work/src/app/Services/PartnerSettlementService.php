<?php


namespace App\Services;


use App\Models\DistributionOrderSettlement;
use App\Services\Settlement\GoodsService;

class PartnerSettlementService extends BaseService
{
    public function SettlementList($input = [])
    {
        $page = $input['page'] ?? 1;
        $size = $input['size'] ?? 10;

        $orderBy = $input['orderby'] ?? 'desc';

        $sql = DistributionOrderSettlement::query()
            ->leftJoin('partner_member', 'distribution_order_settlement.partner_id', '=', 'partner_member.id')
            ->where('distribution_order_settlement.enterprise_id', $input['enterprise_id'])
            ->where('partner_member.id', '<>', '');

        if (isset($input['worker_ids']) && !empty($input['worker_ids'])) {
            $sql->whereIn('distribution_order_settlement.promotion_worker_id', $input['worker_ids']);
        }

        if (isset($input['search_type']) && !empty($input['search_type']) && !empty($input['search_value'])) {
            switch ($input['search_type']) {
                case 1:
                    $sql->where('distribution_order_settlement.order_no', '=', $input['search_value']);
                    break;
                case 2:
                    $sql->where('distribution_order_settlement.id', '=', $input['search_value']);
                    break;
                default:
                    break;
            }
        }

        if (isset($input['status_type']) && !empty($input['status_type']) && is_numeric($input['status_value'])) {
            switch ($input['status_type']) {
                case 1:
                    $sql->where('distribution_order_settlement.settlement_status', '=', $input['status_value']);
                    break;
                case 2:
                    $sql->where('distribution_order_settlement.order_status', '=', $input['status_value']);
                    break;
                default:
                    break;
            }
        }

        if (isset($input['order_by_type']) && !empty($input['order_by_type'])) {
            switch ($input['order_by_type']) {
                case 1:
                    $sql->orderBy('distribution_order_settlement.partner_commission', $orderBy);
                    break;
                case 2:
                    $sql->orderBy('distribution_order_settlement.promotion_commission', $orderBy);
                    break;
                default:
                    $sql->orderBy('distribution_order_settlement.order_pay_time', $orderBy);
                    break;
            }
        }else{
            $sql->orderBy('distribution_order_settlement.order_pay_time', $orderBy);
        }

        $total = $sql->count();

        $sql = $sql->forPage($page, $size);


        $list = $sql->get([
            'distribution_order_settlement.id',
            'distribution_order_settlement.order_no',
            'distribution_order_settlement.order_status',
            'distribution_order_settlement.settlement_status',
            'distribution_order_settlement.order_pay_time',
            'distribution_order_settlement.partner_id',
            'distribution_order_settlement.partner_commission',
            'distribution_order_settlement.promotion_commission',
            'distribution_order_settlement.promotion_worker_id',
            'partner_member.partner_name',
            'partner_member.partner_name',
            'partner_member.avatar'
        ])
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

        $worker_ids = array_column($list, 'promotion_worker_id');
        $search = [
            'enterprise_id' => $enterpriseId,
            'ids' => $worker_ids ?? [],
        ];;
        $workInfo = StoreWorkerService::list($search,false)['list'];
        $workerList = array_column($workInfo, null, 'id');

        foreach ($list as $value) {

            $workName = $workerList[$value['promotion_worker_id']]['name'] ?? '';
            $workName .= !empty($workerList[$value['promotion_worker_id']]['worker_status']) && $workerList[$value['promotion_worker_id']]['worker_status'] == 2 ? '（已移除）' : '';
            $value['worker_name'] = $workName;
            $value['worker_avatar'] = $workerList[$value['promotion_worker_id']]['avatar'] ?? '';
            $data[] = $value;
        }

        return $data;
    }

    public function Detail($input = [])
    {
        $id = $input['id'];
        $enterpriseId = $input['enterprise_id'];
        $detail = DistributionOrderSettlement::query()
            ->leftJoin('partner_member', 'distribution_order_settlement.partner_id', '=', 'partner_member.id')
            ->where('distribution_order_settlement.enterprise_id', $enterpriseId)
            ->where('distribution_order_settlement.id', $id)
            ->first([
                'distribution_order_settlement.id',
                'distribution_order_settlement.order_no',
                'distribution_order_settlement.order_status',
                'distribution_order_settlement.settlement_status',
                'distribution_order_settlement.order_pay_time',
                'distribution_order_settlement.partner_id',
                'distribution_order_settlement.partner_commission',
                'distribution_order_settlement.promotion_commission',
                'distribution_order_settlement.promotion_worker_id',
                'partner_member.partner_name',
                'partner_member.avatar as partner_avatar',
                'distribution_order_settlement.service_shopping_guid_id',
                'distribution_order_settlement.service_commission',
            ])
            ->toArray();
        $search = [
            'enterprise_id' => $enterpriseId,
            'ids' => [$detail['promotion_worker_id']] ?? [],
        ];
        $workInfo = StoreWorkerService::list($search,false)['list'];
        $workInfo =  $workInfo[0] ?? [];
        $workName = $workInfo['name'] ?? '';
        $workName .= !empty($workInfo['worker_status']) && $workInfo['worker_status'] == 2 ? '（已移除）' : '';
        $detail['worker_name'] = $workName;
        $detail['worker_avatar'] = $workInfo['avatar'] ?? '';
        $goodsService = new GoodsService($enterpriseId);
        $goodsService->setSettlement($detail);

        return [
            'detail' => $detail ?? [],
            'settlement_goods' => $goodsService->getGoodsInfo()
        ];
    }

}