<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/10/15 0015
 * Time: 21:24
 */

namespace App\Services\Settlement;


use App\Common\Constant;
use App\Models\DistributionOrderSettlement;
use App\Services\BaseService;
use App\Services\StoreWorkerService;
use Illuminate\Contracts\Pagination\Paginator;
use Illuminate\Database\Eloquent\Builder;

class ListService extends BaseService
{
    private $queryParam;

    private $enterpriseId;

    private $page;

    private $size;

    private $orderNo;

    private $guideUser;

    private $settlementNo;

    private $orderStatus;

    private $startTime;

    private $endTime;

    private $settlementStatus;

    private $guideWorkers;

    private $promoteGuideWorkers;

    private $serviceGuideWorkers;

    private $settlementIds;

    private $storeVersion;

    private $guideService;

    public function __construct($queryParam)
    {
        $this->queryParam = $queryParam;

        $this->init();

        $this->guideService = new GuideService($this->enterpriseId, $this->storeVersion);
    }

    /**
     * 初始化查询参数
     */
    private function init()
    {
        $this->enterpriseId = $this->queryParam['enterprise_id'];
        $this->page = $this->queryParam['page'] ?? Constant::PAGE;
        $this->size = $this->queryParam['size'] ?? Constant::DEFAULT_SIZE;
        if (
            isset($this->queryParam['search_type']) &&
            ($searchType = $this->queryParam['search_type']) &&
            isset($this->queryParam['search_value']) &&
            ($searchValue = $this->queryParam['search_value']) !== ""
        ) {
            switch ($searchType) {
                case 1:
                    $this->orderNo = $searchValue;
                    break;
                case 2:
                    $this->guideUser = $searchValue;
                    break;
                case 3:
                    $this->settlementNo = $searchValue;
                    break;
            }
        }
        if(isset($this->queryParam['guide_workers_scope_type']) &&
            ($guideWorkersScopeType = $this->queryParam['guide_workers_scope_type']) &&
            isset($this->queryParam['guide_workers']) &&
            ($guideWorkersIds = $this->queryParam['guide_workers'])
        ){
            switch ($guideWorkersScopeType) {
                case 1:
                    $this->guideWorkers = $guideWorkersIds;
                    break;
                case 2:
                    $this->promoteGuideWorkers = $guideWorkersIds;
                    break;
                case 3:
                    $this->serviceGuideWorkers = $guideWorkersIds;
                    break;
            }
        }
        $this->orderStatus = $this->queryParam['order_status'] ?? null;
        $this->startTime = $this->queryParam['start_time'] ?? null;
        $this->endTime = $this->queryParam['end_time'] ?? null;
        $this->settlementStatus = $this->queryParam['settlement_status'] ?? null;
        $this->settlementIds = $this->queryParam['settlement_ids'] ?? [];
        $this->storeVersion = $this->queryParam['version'] ?? "";
    }

    /**
     * 列表接口
     *
     * @return array
     */
    public function list()
    {
        $query = DistributionOrderSettlement::query()
            ->where('enterprise_id', $this->enterpriseId)
            ->when($this->orderNo, function (Builder $query) {
                $query->where('order_no',  $this->orderNo );
            })
            ->when($this->guideUser, function (Builder $query) {
                $workIds = $this->getWorkId();
                $query->where(function (Builder $query) use ($workIds) {
                    $query->whereIn('promotion_worker_id', $workIds)
                        ->orWhereIn('service_shopping_guid_id', $workIds);
                });

            })
            ->when($this->settlementNo, function (Builder $query) {
                $query->where('id', $this->settlementNo);
            })
            ->when($this->orderStatus, function (Builder $query) {
                $query->where('order_status', $this->orderStatus);
            })
            ->when($this->startTime, function (Builder $query) {
                $query->where('created_at', '>=', $this->startTime);
            })
            ->when($this->endTime, function (Builder $query) {
                $query->where('created_at', '<=', $this->endTime);
            })
            ->when(isset($this->settlementStatus) && ($this->settlementStatus !== ''), function (Builder $query) {
                $query->where('settlement_status', $this->settlementStatus);
            })
            ->when($this->promoteGuideWorkers,function (Builder $query){
                $query->whereIn('promotion_worker_id',$this->promoteGuideWorkers);
            })
            ->when($this->serviceGuideWorkers,function (Builder $query){
                $query->whereIn('service_shopping_guid_id',$this->serviceGuideWorkers);
            })
            ->when($this->guideWorkers, function (Builder $query) {
                $query->where(function (Builder $query) {
                    $query->whereIn('promotion_worker_id', $this->guideWorkers)
                        ->orWhereIn('service_shopping_guid_id', $this->guideWorkers);
                });

            })
            ->when($this->settlementIds, function (Builder $query) {
                $query->whereIn('id', $this->settlementIds);
            })
            ->orderBy('created_at', 'desc');
//        if (config('app.env') == 'dev') { // 打印sql
//            try {
//                $query->get(['some_field']);
//            } catch (\Exception $e) {
//                file_put_contents('sql.txt', $e->getMessage() . PHP_EOL, FILE_APPEND);
//            }
//        }


        $pageData = $query->paginate($this->size, $this->getField(), '', $this->page);
        return $this->formData($pageData);
    }

    /**
     * 名称查询导购id
     *
     * @return array
     */
    private function getWorkId()
    {
        $worker = StoreWorkerService::list(
            [
                'enterprise_id' => $this->enterpriseId,
                'name'          => $this->guideUser
            ]);
        return array_column($worker['list'], 'id');
    }

    /**
     * 筛选所需字段
     *
     * @return array
     */
    private function getField()
    {
        return [
            'distribution_order_settlement.id',
            'distribution_order_settlement.order_no',
            'distribution_order_settlement.order_status',
            'distribution_order_settlement.settlement_status',
            'distribution_order_settlement.order_pay_time',
            'distribution_order_settlement.promotion_worker_id',
            'distribution_order_settlement.promotion_commission',
            'distribution_order_settlement.service_shopping_guid_id',
            'distribution_order_settlement.service_commission',
            'distribution_order_settlement.service_shopping_guid_id'
        ];

    }

    /**
     * 组装前端数据
     *
     * @param Paginator $pageData
     * @return array
     */
    private function formData(Paginator $pageData)
    {

        return [
            'list'       => $this->compat($pageData->items()),
            'statistics' => [
                'wait_settle'   => $this->totalStatistics(DistributionOrderSettlement::WAIT_SETTLE),
                'adopt_settle'  => $this->totalStatistics(DistributionOrderSettlement::ADOPT_SETTLE),
                'cancel_settle' => $this->totalStatistics(DistributionOrderSettlement::CANCEL_SETTLE)
            ],
            'total'      => $pageData->total() ?? 0,
            'page'       => (int)$this->page,
            'size'       => (int)$this->size
        ];
    }

    /**
     * 组装列表格式
     *
     * @param array $settlementList
     * @return array
     */
    private function compat(array $settlementList)
    {

        $formatUserArr = [];
        $this->guideService->setGuideIds(
            array_unique(array_column($settlementList, 'promotion_worker_id')),
            array_unique(array_column($settlementList, 'service_shopping_guid_id'))
        );

        foreach ($settlementList as $settlement) {
            $this->guideService->setUser($settlement->toArray());
            $formatUserArr[] = $this->guideService->formatSingle();
        }

        return $formatUserArr;
    }

    /**
     * 统计佣金
     *
     * @param $status
     * @return int
     */
    private function totalStatistics($status)
    {
        if(in_array($status,[DistributionOrderSettlement::WAIT_SETTLE,DistributionOrderSettlement::ADOPT_SETTLE])){
            $promotionCommission = DistributionOrderSettlement::query()
                ->where('enterprise_id', $this->enterpriseId)
                ->where('settlement_status', $status)
                ->where('promotion_worker_id', '!=', 0)
                ->sum('promotion_commission');

            $serviceCommission = DistributionOrderSettlement::query()
                ->where('enterprise_id', $this->enterpriseId)
                ->where('settlement_status', $status)
                ->where('service_shopping_guid_id', '!=', 0)
                ->sum('service_commission');
        }else{
            $promotionCommission = DistributionOrderSettlement::query()
                ->where('enterprise_id', $this->enterpriseId)
                ->where('settlement_status', $status)
                ->where('promotion_worker_id', '!=', 0)
                ->sum('promotion_commission_snapshot');

            $serviceCommission = DistributionOrderSettlement::query()
                ->where('enterprise_id', $this->enterpriseId)
                ->where('settlement_status', $status)
                ->where('service_shopping_guid_id', '!=', 0)
                ->sum('service_commission_snapshot');
        }
        
        return (int)$promotionCommission + (int)$serviceCommission;
    }

}