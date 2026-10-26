<?php
/**
 * 首页仪表盘模块
 *
 * @author  纪鑫鑫 <jixinxin@vchangyi.com>
 * @date    2020-08-03 10:07
 */

namespace App\Http\Controllers\Service\V1\Dashboard;

use App\Common\Constant;
use App\Common\DateUtils;
use App\Http\Controllers\Service\V1\BaseController;
use App\Models\CrsOrder;
use App\Models\DailyWorkerCustomerStatistic;
use App\Models\MemberUserRelation;
use App\Models\Scrm\Members;
use App\Models\Scrm\WorkExternalUserRelation;
use App\Models\ScrmEnterprise\WxUserCorpExternal;
use App\Models\StoreWorker;
use App\Services\DataStatistic\CustomerDataCompassService;
use App\Services\PartnerMemberService;
use CrsSdk\Shop\Order;
use Exception;
use Illuminate\Http\Request;
use Illuminate\Validation\ValidationException;

/**
 * Class IndexController
 *
 * @package App\Http\Controllers\Service\V1\Dashboard
 * @author  纪鑫鑫 <jixinxin@vchangyi.com>
 * @date    2020-08-03 10:07
 */
class IndexController extends BaseController
{
    // 导购订单类型
    const  WORKER_ORDER_TYPE = 2;

    // 推广订单查询状态
    const ORDER_PROMOTER_STATUS = [40, 60, 41, 80, 81, 83, 100, 200, 205];

    /**
     * 获取首页仪表盘数据
     *
     * @param Request $request
     *
     * @return array
     * @throws ValidationException
     * @author  纪鑫鑫 <jixinxin@vchangyi.com>
     * @date    2020-08-20 16:03
     */
    public function index(Request $request)
    {
        $params = $this->validate(
            $request,
            [
                'enterprise_id' => "required",
                'mem_uid' => "required",
                //'contact_qr_code_id' => "required",
                //'store_id'           => "required",
                'worker_id' => "required",
                //'role'               => "required",
                'tab_time' => "required|integer",
                //'service_tab_time'   => "required|integer",
                //'extension_tab_time' => "required|integer",
            ],
            [
                'enterprise_id.required' => '企业参数不可为空',
                //'store_id.required'      => '门店参数不可为空',
                //'mem_uid.required'       => '员工参数不可为空',
                'worker_id.required' => '导购参数不可为空',
                'contact_qr_code_id.required' => '导购二维码ID不能为空',
            ]
        );

        // 切换域名
        switchAppHost(['enterprise_id' => $params['enterprise_id']]);

        //根据导购二维码id统计外部联系人数
        $memberUserInfo = MemberUserRelation::query()
//            ->where('source_type', '=', $params['contact_qr_code_id'])
            ->where('delete_type', '=', 0)
            ->where('mem_uid', $params['mem_uid'])
            ->orderBy('created_at', 'desc')
            ->get(['union_id', 'created_at'])
            ->toArray();

        $updatedTime = 0;
        if (!empty($memberUserInfo)) {
            $updatedTime = $memberUserInfo[0]['created_at'];
        }

        $memberUserNumber = count($memberUserInfo) ?? 0;
        $memberRedDot = false;
        //$payedUserCount   = 0;
        // 判定顾客数据小红点，未浏览列表的数据计入
        if ($params['tab_time'] < $updatedTime) {
            $memberRedDot = true;
        }

        /*if (!empty($memberUserInfo)) {
            $unionIds = array_values(array_filter(array_column($memberUserInfo, 'union_id')));
            if (!empty($unionIds)) {
                $payedUserCount = MemberOrderStatistics::query()
                    ->where('enterprise_id', $params['enterprise_id'])
                    ->whereIn('union_id', $unionIds)
                    ->count();
            }
        }*/

        //查询出导购关联的合伙人
        $params['promoter_ids'] = [];
        if ($params['worker_id'] > 0) {
            $partnerIds = (new PartnerMemberService())->getPartnerIdsByWorkerIds($params['enterprise_id'], [$params['worker_id']]);
            if (!empty($partnerIds)) {
                $partnerIds[] = $params['worker_id'];
                $params['promoter_ids'] = $partnerIds;
            }
        }

//        $haoWuJiEnterpriseId = env('HAO_WU_JI_ENTERPRISE_ID');
//        if ($haoWuJiEnterpriseId == $params['enterprise_id']) {
//            list($extensionOrder, $serviceOrder, $workerOrder) = $this->getHaoWuJiOrderData($params);
//        } else {
//            list($extensionOrder, $serviceOrder, $workerOrder) = $this->getOrderData($params);
//        }

        // 绑定顾客数
        $bindingCustomerTotal = 0;
        $statistic = DailyWorkerCustomerStatistic::getWorkerLatestStatistic($params['enterprise_id'], $params['worker_id'], DateUtils::todayStatTime());
        if (!empty($statistic)) {
            $bindingCustomerTotal = $statistic['cumulative_binding_customer_num'];
        }

        // 今日推广实收金额
        $extensionAmountParams = [
            'enterprise_id' => $params['enterprise_id'],
            'promoter_ids' => array_unique(array_merge([$params['worker_id']], $params['promoter_ids'])),
            'start_time' => DateUtils::todayStatTime(),
            'end_time' => time()
        ];
//        $extensionAmount = (new Order())->extensionActuallyAmount($extensionAmountParams);
//        $todayExtensionActuallyAmount = $extensionAmount['extension_amount'];

        // 会员数量分布
        $memId = $this->getMemId($params['enterprise_id'], $params['mem_uid']);
        $memberDistribution = $this->getMemberDistribution($params['enterprise_id'], $memId);

        // 会员等级分布
        $memberLevelStatistics = (new CustomerDataCompassService())->getMemberLevelStatisticsNew(
            $params['enterprise_id'],
            $memId
        );

        return success(
            [
                'member_user_number' => $memberUserNumber, // 我的客户数
                'need_order' => 0, // 需处理订单
                //'payed_user_count'      => $payedUserCount, // 成交客户数
                'extension_order_count' => 0, // 推广订单数
                'member_red_dot' => $memberRedDot, // 顾客红点
                'extension_order' => false, // 推广订单红点
                'service_order' => false,  // 待接单红点
                /*'available_cash_amount' => ShoppingGuideCashWithdrawal::query()
                    ->where('worker_id', '=', $params['worker_id'])
                    ->value('available_cash_amount') ?? 0,*/ // 可提现佣金红点
                'binding_customer_total' => $bindingCustomerTotal, // 绑定顾客数
                'today_extension_actually_amount' => 0, // 今日推广实收金额
                'member_count_statistics' => $memberDistribution, // 会员数量分布
                'member_level_statistics' => $memberLevelStatistics['memberGrade'], // 会员数量分布
                'member_total_count' => $memberLevelStatistics['memberTotal']// 会员总数

            ]
        );
    }

    protected function getMemId(int $enterpriseId, string $ucMemId)
    {
        $storeWorker = StoreWorker::query()
            ->select('mem_userid')
            ->where('enterprise_id', $enterpriseId)
            ->where('mem_uid', $ucMemId)
            ->where('worker_status', 1)
            ->first();

        if ($storeWorker) {
            $userId = $storeWorker->user_id;
        }
        if (!empty($userId)) {
            return Members::query()
                ->select('members.id')
                ->where('userid', $userId)
                ->where('enterprise_id', $enterpriseId)
                ->first()
                ->id;
        }
        return 0;
    }

    /**
     * 会员分布
     * @param int $enterpriseId
     * @param int $memUid
     * @return array
     * @Date: 2021/3/31
     * @Time: 17:35
     * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
     */
    protected function getMemberDistribution(int $enterpriseId, string $memUid): array
    {
        //$memberCount    = $this->getMemberCount($enterpriseId, $memUid, Constant::IS_MEMBER_STATISTICS);
        //$nonMemberCount = $this->getMemberCount($enterpriseId, $memUid, Constant::NON_MEMBER_STATISTICS);
        if ($memUid) {
            list($vipNum, $totalNum) = $this->getMemberCount($enterpriseId, $memUid);
            return [
                'member_count' => $vipNum,
                'non_member_count' => $totalNum
            ];
        }
        return [
            'member_count' => 0,
            'non_member_count' => 0
        ];
    }

    protected function getMemberCount(int $enterpriseId, string $MemId)
    {
        $userList = WorkExternalUserRelation::query()
            ->select('work_external_user_relation.work_external_user_id')
            ->where('work_external_user_relation.enterprise_id', $enterpriseId)
            ->where('work_external_user_relation.member_id', $MemId)
            ->where('work_external_user_relation.relation_status', 0)
            ->pluck('work_external_user_id')
            ->toArray();
        $vipNum = WxUserCorpExternal::query()
            ->setConn($enterpriseId)
            ->select('wx_user_corp_external.id')
            ->whereIn('wx_user_corp_external.id', $userList)
            ->where('wx_user_corp_external.grade', '<>', '')
            ->pluck('id')
            ->toArray();
        $vipNum = count($vipNum);
        $totalNum = count($userList);

        return [$vipNum, ($totalNum - $vipNum)];
    }

    /**
     * 获取企业订单数据
     * @param $params
     *
     * @return array
     * @throws Exception
     */
    protected function getOrderData($params)
    {

        $extensionOrder = false;
        $serviceOrder = false;

        //查询首页推广订单和待抢单是否新的订单
        $serviceOrderNumbersByTime = (new Order())->orderOutStatistics([
                "type" => self::WORKER_ORDER_TYPE,
                'params'
                =>
                    [
                        'enterprise_id' => $params['enterprise_id'],
                        'worker_id' => $params['worker_id'],
                        'store_id' => $params['store_id'],
                        'user_start_time' => 0,
                        'extension_start_time' => $params['extension_tab_time'],
                        'receive_start_time' => $params['service_tab_time'],
                        'promoter_ids' => $params['promoter_ids'],
                        'extension_order_status' => self::ORDER_PROMOTER_STATUS
                    ]
            ]
        );

        // 推广订单数,需处理订单数
        $workerOrder = (new Order())->orderOutStatistics([
                "type" => self::WORKER_ORDER_TYPE,
                'params'
                =>
                    [
                        'enterprise_id' => $params['enterprise_id'],
                        'worker_id' => $params['worker_id'],
                        'store_id' => $params['store_id'],
                        'user_start_time' => 0,
                        'extension_start_time' => 0,
                        'receive_start_time' => 0,
                        'promoter_ids' => $params['promoter_ids'],
                        'extension_order_status' => self::ORDER_PROMOTER_STATUS
                    ]
            ]
        );

        if ($serviceOrderNumbersByTime['extension_order_count'] > 0) {
            $extensionOrder = true;
        }
        if ($serviceOrderNumbersByTime['receive_order_count'] > 0) {
            $serviceOrder = true;
        }

        return [$extensionOrder, $serviceOrder, $workerOrder];
    }

    /**
     * 好物集订单
     * @param $params
     *
     * @return array
     */
    protected function getHaoWuJiOrderData($params)
    {

        $extensionOrder = false;

        $extensionOrderCount = CrsOrder::query()
            ->whereIn('status', self::ORDER_PROMOTER_STATUS)
            ->where('store_staff_id', $params['worker_id'])
            ->where('created_at', $params['extension_tab_time'])
            ->where('sync_source', CrsOrder::HAO_WU_JI_ORDER)
            ->count();

        if ($extensionOrderCount > 0) {
            $extensionOrder = true;
        }

        $workerOrder['extension_order_count'] = CrsOrder::query()
            ->whereIn('status', self::ORDER_PROMOTER_STATUS)
            ->where('store_staff_id', $params['worker_id'])
            ->where('sync_source', CrsOrder::HAO_WU_JI_ORDER)
            ->count();

        return [$extensionOrder, false, $workerOrder];

    }

    /**
     * 获取会员或非会员数量
     * @param int $enterpriseId
     * @param int $memUid
     * @param bool $isMember
     * @return int
     * @Date: 2021/3/31
     * @Time: 17:36
     * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
     */
    protected function getMemberCountBAK(int $enterpriseId, string $memUid, int $isMember = Constant::IS_MEMBER_STATISTICS): int
    {
        return MemberUserRelation::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('mem_uid', $memUid)
            ->when(!empty($isMember), function ($query) use ($isMember) {
                if ($isMember == Constant::IS_MEMBER_STATISTICS) {
                    $query->where('union_code', '<>', '');
                } else if ($isMember == Constant::NON_MEMBER_STATISTICS) {
                    $query->where('union_code', '');
                }
            })
            ->where('delete_type', 0)
            ->count();
    }

}
