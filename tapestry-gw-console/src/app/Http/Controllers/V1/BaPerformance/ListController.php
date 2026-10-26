<?php

namespace App\Http\Controllers\V1\BaPerformance;

use App\Common\Constant;
use CrsSdk\GoMicroScrm\BaPerformance;
use CrsSdk\Work\BaseClient;
use CrsSdk\Work\StoreWorker;
use CrsSdk\Work\Member;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;
use Maatwebsite\Excel\Facades\Excel;

class ListController extends BaseController
{
    /**
     * 指定导购的业绩指标
     * @return JsonResponse
     * @throws \Exception
     */
    public function index()
    {
        $params                  = $this->request->only(['search', 'page', 'size']);
        $params['user_ids']      = $this->getUserIds($params);
        $params['enterprise_id'] = $this->enterpriseId;
        $params['is_page']       = 1;
        $params['name']          = empty($params['search']['name']) ? '' : $params['search']['name'];
        try {
            $result = (new BaPerformance())->List($params);
        } catch (\Exception $e) {
            error($e->getCode(), $e->getMessage());
        }
        return success($result['data'] ?? $result);
    }

    private function getUserIds($params)
    {
        $memList     = [];
        $typeMemList = [];
        //指定部门下的所有员工
        if (!empty($params['search']['dp_id'])) {
            // 员工搜索条件
            $condition                   = [
                'qy_domain'      => $this->qyDomain,
                'app_identifier' => $this->appIdentifier,
                'params'         => [
                    'pageNum'  => Constant::DEFAULT_PAGE,
                    'pageSize' => Constant::MAX_DEFAULT_SIZE,
                ],
                'enterprise_id'  => $this->enterpriseId
            ];
            $condition['params']['dpId'] = $params['search']['dp_id'];
            //获得当前及子目录的所有数据
            $condition['params']['departmentChildrenFlag'] = 1;
            $memList                                       = (new Member())->memberList($condition);
            $memList                                       = array_column($memList['list'] ?? [], 'memMobile');
        }
        //导购类型
        if (!empty($params['search']['type']) || !empty($params['search']['worker_no'])) {
            // 导购搜索条件
            $search = [
                'page'          => Constant::DEFAULT_PAGE,
                'size'          => Constant::MAX_DEFAULT_SIZE,
                'enterprise_id' => $this->enterpriseId,
                'worker_status' => Constant::IN_SERVICE
            ];
            if (!empty($params['search']['type'])) {
                $search['type'] = $params['search']['type'];
            }
            if (!empty($params['search']['worker_no'])) {
                $search['worker_no'] = $params['search']['worker_no'];
            }
            $typeMemList    = (new StoreWorker())->getList($search);
            $typeMemList    = array_column($typeMemList['list'] ?? [], 'mobile');
        }

        if (!empty($memList) && !empty($typeMemList)) {
            $userIds = array_intersect($memList, $typeMemList);
        } else {
            $userIds = array_merge($memList, $typeMemList);
        }
        return $userIds;
    }

    /**
     * 设置导出列宽
     *
     * @var array
     */
    protected $colWidth = [
        'A' => 15,
        'B' => 15,
        'C' => 15,
        'D' => 15,
        'E' => 15,
        'F' => 15,
        'G' => 15,
        'H' => 20,
        'I' => 25,
        'J' => 25,
        'K' => 25,
    ];

    public function export()
    {
        $params                  = $this->request->only(['search', 'page', 'size']);
        $params['user_ids']      = $this->getUserIds($params);
        $params['enterprise_id'] = $this->enterpriseId;
        $params['is_page']       = 0;
        $params['name']          = empty($params['search']['name']) ? '' : $params['search']['name'];

        try {
            $result = (new BaPerformance())->List($params);
        } catch (\Exception $e) {
            error($e->getCode(), $e->getMessage());
        }
        $cellData[] = [
            '区域（大区）',
            '区域',
            '城市',
            '柜台code',
            '柜台名称',
            '导购BA姓名',
            '手机号',
            '月度销售目标（元）',
            '新二购且一码三关注月度目标（人）',
            '老客购且一码三关注月度目标（人）',
            '新客购且一码三关注月度目标（人）',
        ];
        foreach ($result['list'] ?? [] as $value) {
            $cellData[] = [
                $value['regional'],
                $value['area'],
                $value['city'],
                $value['counter_code'],
                $value['counter_name'],
                $value['ba_name'],
                $value['mobile'],
                $value['monthly_sales'],
                $value['monthly_second_buy_and_three_scan'],
                $value['monthly_old_customer_and_three_scan'],
                $value['monthly_new_customer_and_three_scan'],
            ];
        }

        $name = 'BA业绩目标导入模板';
        return Excel::create($name, function ($excel) use ($cellData) {
            $excel->sheet('ba-performance', function ($sheet) use ($cellData) {
                $sheet->setWidth($this->colWidth)->rows($cellData);
            });
        })->export('xlsx');
    }


}
