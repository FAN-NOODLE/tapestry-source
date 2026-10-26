<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/27
 * Time: 11:02
 */

namespace App\Http\Controllers\V1\Customer;

use App\Http\Controllers\V1\BaseController;
use App\Jobs\Export\CustomerStatisticExportJob;
use CrsSdk\Shop\User;
use CrsSdk\Work\DataCompass;
use CrsSdk\Work\Member;
use Illuminate\Http\Request;

class StatisticsController extends BaseController
{
    const STAT_COLUMN_NAME_MAP = [
        'customer_total' => 'cumulative_customer_num',
        'guide_customer_total' => 'cumulative_contact_card_way_customer_num',
        'other_customer_total' => 'cumulative_other_way_customer_num',
        'customer_new' => 'today_customer_num',
        'guide_customer_new' => 'today_contact_card_way_customer_num',
        'other_customer_new' => 'today_other_way_customer_num',
        'chat_total' => 'chat_cnt',
        'reply_percentage' => 'reply_percentage',
        'binding_customer_total' => 'cumulative_binding_customer_num',
        'binding_customer_new' => 'today_binding_customer_num',
        'member_total' => 'cumulative_member_num',
        'member_new' => 'today_member_num',
        'non_member_total' => 'cumulative_non_member_num',
        'non_member_new' => 'today_non_member_num',
    ];

    const OVERVIEW_COLUMN_NAME_MAP = [
        'data1' => 'cumulative_customer_num',
        'data2' => 'today_customer_num',
        'data3' => 'customer_access',
        'data4' => 'customer_deal',
        'data5' => 'cumulative_contact_card_way_customer_num',
        'data6' => 'cumulative_other_way_customer_num',
        'data7' => 'today_contact_card_way_customer_num',
        'data8' => 'today_other_way_customer_num',
        'data9' => 'chat_cnt',
        'data10' => 'reply_percentage',
        'data11' => 'cumulative_binding_customer_num',
        'data12' => 'today_binding_customer_num',
        'data13' => 'cumulative_member_num',
        'data14' => 'today_member_num',
        'data15' => 'cumulative_non_member_num',
        'data16' => 'today_non_member_num',
    ];

    /**
     * 客户统计数据
     *
     * @param Request $request
     *
     * @return array|\Illuminate\Http\JsonResponse
     */
    public function index(Request $request)
    {
        // 客户id
        $userId = $request->get('user_id', 0);
        if (!$userId) {
            error(320051);
        }

        // 调用sdk
        $params     = [
            'enterprise_id' => $this->enterpriseId,
            'user_id'       => $userId
        ];
        $statistics = (new User())->count($params);

        $data = [
            'user_id'       => intval($userId),
            'total_amount'  => $statistics['data']['yearSpend'],
            'avg_amount'    => $statistics['data']['yearMeanSpend'],
            'avg_period'    => $statistics['data']['yearMeanSpendDay'],
            'category_name' => $statistics['data']['yearCategoryName']
        ];

        return success(['data' => $data]);
    }

    /**
     * 数据罗盘.顾客数据概览
     * @param Request $request
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function guideCustomer(Request $request)
    {
        $param['start_time']    = $request->input('start_time');
        $param['end_time']      = $request->input('end_time');
        $param['worker_ids']    = $request->input('worker_ids', []);
        $param['enterprise_id'] = $this->enterpriseId;
        $param['qy_domain']     = $this->qyDomain;

        $overview = (new DataCompass())->customerStatisticOverviewForAdmin($param);
        $overviewData = $overview['overview'];
        $statisticList = $overview['detail_list'];
        $memberLevelStatistics = $overview['member_level_statistics'];
        $memberTotalCount = $overview['member_total_count'];
        $memberLevelStatisticsByTimeRange = $overview['member_level_statistics_by_time_range'];
        $memberTotalCountByTimeRange = $overview['member_total_count_by_time_range'];

        $result = [];
        foreach (self::OVERVIEW_COLUMN_NAME_MAP as $key => $value) {
            $result[$key] = $overviewData[$value] ?? 0;
        }
        $x = $y = [];
        if (!empty($statisticList)) {
            $x = array_map(function($statTime) {
                return date('Y-m-d', $statTime);
            }, array_column($statisticList, 'stat_time'));
            $defaultData = array_fill(0, count($x), 0);
            foreach (self::OVERVIEW_COLUMN_NAME_MAP as $key => $value) {
                if (in_array($key, ['data3', 'data4'])) {
                    $y[$key] = $defaultData;
                } else {
                    $y[$key] = array_column($statisticList, $value);
                }
            }
        }

        $result['detail_data'] = [
            'x' => $x,
            'y' => $y
        ];

        $result['member_level_statistics'] = $memberLevelStatistics;
        $result['member_total_count'] = $memberTotalCount;
        $result['member_level_statistics_by_time_range'] = $memberLevelStatisticsByTimeRange;
        $result['member_total_count_by_time_range'] = $memberTotalCountByTimeRange;
        return success($result);

    }

    /**
     * 数据罗盘.顾客数据按时间列表
     * @param Request $request
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function guideCustomerDate(Request $request)
    {
        $param['page'] = $request->input('page', 1);
        $param['size'] = $request->input('size', 10);

        $param['start_time']    = $request->input('start_time');
        $param['end_time']      = $request->input('end_time');
        $param['sort_column']   = self::STAT_COLUMN_NAME_MAP[$request->input('sort_type', '')] ?? '';
        $param['sort_type']     = $request->input('statistical_sort');
        $param['worker_ids']    = $request->input('worker_ids', []);
        $param['enterprise_id'] = $this->enterpriseId;
        $param['qy_domain']     = $this->qyDomain;
        $param['for_page'] = true;

        $data = (new DataCompass())->customerStatisticListByDate($param);
        $resultList = [];
        foreach ($data['list'] as $key => $value) {
            $statistic = [];
            foreach (self::STAT_COLUMN_NAME_MAP as $pageColumnName => $statColumnName) {
                $statistic[$pageColumnName] = $value[$statColumnName] ?? 0;
            }
            $statistic['time'] = date('Y-m-d', $value['stat_time']);
            $resultList[] = $statistic;
        }
        $data['list'] = $resultList;

        return success($data);
    }

    /**
     * 数据罗盘.顾客数据按导购列表
     * @param Request $request
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function guideCustomerList(Request $request)
    {
        $param['page'] = $request->input('page', 1);
        $param['size'] = $request->input('size', 10);

        $param['start_time']    = $request->input('start_time');
        $param['end_time']      = $request->input('end_time');
        $param['sort_column']   = self::STAT_COLUMN_NAME_MAP[$request->input('sort_type', '')] ?? '';
        $param['sort_type']     = $request->input('statistical_sort');
        $param['worker_ids']    = $request->input('worker_ids', []);
        $param['enterprise_id'] = $this->enterpriseId;
        $param['qy_domain']     = $this->qyDomain;
        $param['for_page'] = true;

        $data = (new DataCompass())->customerStatisticListByWorker($param);

        $resultList = [];
        foreach ($data['list'] as $key => $value) {
            $statistic = [];
            foreach (self::STAT_COLUMN_NAME_MAP as $pageColumnName => $statColumnName) {
                $statistic[$pageColumnName] = $value[$statColumnName] ?? 0;
            }
            $resultList[] = array_merge($statistic, [
                'worker_no' => $value['worker_no'],
                'name' => $value['name'],
                'worker_name' => $value['name'],
                'avatar' => $value['avatar'],
                'store_id' => $value['store_id'],
                'mem_uid' => $value['mem_uid'],
                'store_name' => $value['store_name'],
                'dp_name' => $value['dp_name'],
            ]);
        }
        $data['list'] = $resultList;

        return success($data);
    }

    /**
     * @param Request $request
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function guideCustomerListExport(Request $request)
    {
        $param['start_time']    = $request->input('start_time');
        $param['end_time']      = $request->input('end_time');
        $param['sort_column']   = $request->input('sort_type');
        $param['sort_type']     = $request->input('statistical_sort');
        $param['worker_ids']    = $request->input('worker_ids', []);
        $param['export_type']   = $request->input('export_type');
        $param['enterprise_id'] = $this->enterpriseId;
        $param['qy_domain']     = $this->qyDomain;
        $param['operate_name']  = $this->enterpriseAdmin['username'] ?? '';
        $param['operate_id']    = $this->enterpriseAdminId;
        $param['tag']           = $param['tag'] ?? 'export-customer-statistical';
        $param['route_url']     = $param['route_url'] ?? 'customer/statistic-export';
        $param['version']       = $request->input('version');
//        $param['sort_type'] = $this->setSortType($param['sort_type']);
        $param['for_page'] = false;
        $name = '导购顾客数据统计-' . date('YmdHis');
        // 导出
        (new CustomerStatisticExportJob($param, $param, $name))->handle();

        return success();
    }

    private function setSortType($sortType)
    {
        switch ($sortType) {
            case 'desc':
                return SORT_DESC;
            case 'asc':
                return SORT_ASC;
        }
    }
}
