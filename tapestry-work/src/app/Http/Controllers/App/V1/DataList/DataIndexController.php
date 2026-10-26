<?php
/**
 * Created by PhpStorm.
 * User: admin
 * Date: 2019/3/19
 * Time: 16:51
 */

namespace App\Http\Controllers\App\V1\DataList;


use App\Http\Controllers\App\V1\BaseController;
use App\Services\MemberUserRelationService;
use App\Common\Constant;
use CrsSdk\Shop\Order;

class DataIndexController extends BaseController
{
    /**
     * @var MemberUserRelationService
     */
    protected $memberUserRelationService;

    /**
     * CountController constructor.
     *
     * @param MemberUserRelationService $memberUserRelationService
     */
    public function __construct(
        MemberUserRelationService $memberUserRelationService
    ) {
        parent::__construct();
        $this->memberUserRelationService = $memberUserRelationService;
    }

    public function index()
    {
        // 今日时间戳
        $dayStartTime = strtotime(date('Y-m-d', time()) . ' 00:00:00');
        $dayEndTime = strtotime(date('Y-m-d', time()) . ' 23:59:59');

        // 昨日时间戳
        $yesterdayStartTime = strtotime('yesterday');
        $yesterdayEndTime = strtotime('yesterday') + 86400;

        // 今日好友客户数
        $dayUserCount = $this->memberUserRelationService->timeUserCount(
            $this->enterpriseId,
            $this->memUid,
            $dayStartTime,
            $dayEndTime
        );
        // 昨日好友客户数
        $yesterdayUserCount = $this->memberUserRelationService->timeUserCount(
            $this->enterpriseId,
            $this->memUid,
            $yesterdayStartTime,
            $yesterdayEndTime
        );
        // 累计好友客户数
        $userCount = $this->memberUserRelationService->timeUserCount(
            $this->enterpriseId,
            $this->memUid
        );

        // 如果开通了商城
        // 组装返回数据
        $data = [];

        if ($this->storeStatus == Constant::STORE_OPEN) {
            // 参数
            $DayParams = [
                'enterprise_id' => $this->enterpriseId,
                'mem_uid' => $this->memUid
            ];

            $data['total']['order_price'] = 0;
            $data['total']['order_num'] = 0;
            $data['total_day']['order_price'] = 0;
            $data['total_day']['order_num'] = 0;
            $data['total_yesterday']['order_price'] = 0;
            $data['total_yesterday']['order_num'] = 0;

            $list = (new Order())->orderData($DayParams);

            if (!empty($list['data']['total'])) {
                $data['total'] = $list['data']['total'];
            }
            if (!empty($list['data']['total_day'])) {
                $data['total_day'] = $list['data']['total_day'];
            }
            if (!empty($list['data']['total_yesterday'])) {
                $data['total_yesterday'] = $list['data']['total_yesterday'];
            }
            $data['total_day']['scan_num'] = $list['data']['scan_day'];
            $data['total_yesterday']['scan_num'] = $list['data']['scan_yesterday'];
            // 转化一下类型
            $data['total']['order_price'] = doubleval($data['total']['order_price']);
            $data['total_day']['order_price'] = doubleval($data['total_day']['order_price']);
            $data['total_yesterday']['order_price'] = doubleval($data['total_yesterday']['order_price']);
        }
        $data['total']['user_num'] = $userCount;
        $data['total_day']['user_num'] = $dayUserCount;
        $data['total_yesterday']['user_num'] = $yesterdayUserCount;

        return success(['data' => $data]);
    }
}
