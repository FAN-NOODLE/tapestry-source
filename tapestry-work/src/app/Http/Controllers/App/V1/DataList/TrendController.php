<?php
/**
 * Created by PhpStorm.
 * User: admin
 * Date: 2019/3/21
 * Time: 11:20
 */

namespace App\Http\Controllers\App\V1\DataList;


use App\Http\Controllers\App\V1\BaseController;
use Illuminate\Http\Request;
use App\Services\MemberUserRelationService;
use App\Common\Constant;
use CrsSdk\Shop\Order;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Schema;
use App\Models\MemberUserRelation;

class TrendController extends BaseController
{
    /**
     * 本周
     */
    const WEEK = 1;
    /**
     * 本月
     */
    const MONTH = 2;
    /**
     * 本年
     */
    const YEAR = 3;

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

    /**
     * 趋势数据
     *
     * @param Request $request
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index(Request $request)
    {
        //表名
        $table = Constant::TABLE_NAME.$this->enterpriseId;
        // 接值
        $input = $request->only('type');
        // 参数验证规则
        $rules = [
            'type'  => 'required | integer | in:1,2,3',
        ];
        // 提示语
        $message = [
            'type.required'  => 'type不能为空',
            'type.integer'  => 'type必须是正整数',
            'type.in'  => 'type不合法'
        ];
        $this->validate($request, $rules, $message);

        $params = [
            'enterprise_id' => $this->enterpriseId,
            'mem_uid' => $this->memUid
        ];
        // 默认值
        // 组装数据
        $data = [];
        $userNumTotal = 0;
        $orderPriceTotal = 0;
        $orderNumTotal = 0;
        $scanNum = 0;
        $data['user_num']['total'] = 0;
        if ($this->storeStatus == Constant::STORE_OPEN) {
            $data['order_price']['total'] = 0;
            $data['order_num']['total'] = 0;
            $data['scan_num']['total'] = 0;
        }

        $start = mktime(0,0,0,date('m'),date('d'),date('Y'));
        $now = time();
        // 判断type
        if ($input['type'] == self::WEEK) {
            $starts = $start - 6*86400;
        }
        if ($input['type'] == self::MONTH) {
            $starts = $start - 29*86400;
        }
        if ($input['type'] == self::YEAR) {
            $starts = $start - 89*86400;
        }
        // 检测是否有表
        $this->checkTable();
        // 检测是否已加客户
        $member = MemberUserRelation::query()
            ->where([
                'enterprise_id' => $this->enterpriseId,
                'mem_uid' => $this->memUid
            ])
            ->first();
        if (empty($member)) {
            DB::table($table)->where('mem_uid', $this->memUid )->delete();
            $this->initData($start);
        }
        // 会员数据异常重新初始化
        if (!empty($member)) {
            $trend = DB::table($table)->where('mem_uid', $this->memUid )->first();
            if (empty($trend)) {
                $this->trendInit($table);
            }
        }
        $dataList = DB::table($table)->where($params)->whereBetween('day_time', [$starts, $start])->get()->toArray();
        $dataList = array_reverse($dataList);
        // 封装数据
        foreach($dataList as $value) {
            // 客户数
            $data['user_num']['array'][] = ['data' => $value->day_time, 'value' => $value->customer_num];
            $userNumTotal += $value->customer_num;
            $data['user_num']['total'] = $userNumTotal;
            $data['user_num']['name'] = '好友客户数';
            // 订单数
            if ($this->storeStatus == Constant::STORE_OPEN) {
                $data['order_price']['array'][] = ['data' => $value->day_time, 'value' => doubleval($value->order_price)];
                $data['order_num']['array'][] = ['data' => $value->day_time, 'value' => $value->order_num];
                $data['scan_num']['array'][] = ['data' => $value->day_time, 'value' => $value->scan_num];

                $orderPriceTotal += doubleval($value->order_price);
                $orderNumTotal += $value->order_num;
                $scanNum += $value->scan_num;

                $data['order_price']['total'] = $orderPriceTotal;
                $data['order_num']['total'] = $orderNumTotal;
                $data['scan_num']['total'] = $scanNum;
                $data['order_price']['name'] = '成交额';
                $data['order_num']['name'] = '成交订单数';
                $data['scan_num']['name'] = '浏览人数';
            }
        }
        // 今日数据
        // 获取客户数
        $userNum = $this->memberUserRelationService->timeUserCount(
            $this->enterpriseId,
            $this->memUid,
            $start,
            $now
        );
        $data['user_num']['array'][] = ['data' => $start, 'value' => $userNum];
        $data['user_num']['total'] += $userNum;
        // 如果开通了商城\
        $param = [
            'enterprise_id' => $this->enterpriseId,
            'mem_uid' => $this->memUid,
            'start_time' => $start,
            'end_time' => $now,
            'type' => 1
        ];
        if ($this->storeStatus == Constant::STORE_OPEN) {
            $list = (new Order())->orderByTime($param);
            // 封装数据
            if (empty($list['order_price'])) {
                $list['order_price'] = '';
            }
            if (empty($list['order_num'])) {
                $list['order_num'] = 0;
            }
            $data['order_price']['array'][] = ['data' => $start, 'value' => doubleval($list['order_price'])];
            $data['order_num']['array'][] = ['data' => $start, 'value' => $list['order_num']];
            $data['scan_num']['array'][] = ['data' => $start, 'value' => $list['scan_person_num']];
            $data['order_price']['total'] += doubleval($list['order_price']);
            $data['order_num']['total'] += $list['order_num'];
            $data['scan_num']['total'] += $list['scan_person_num'];
        }

        return success(['data' => $data]);
    }

    /**
     * 初始化
     *
     * @param $start
     */
    public function initData($start)
    {
        //表名
        $table = Constant::TABLE_NAME.$this->enterpriseId;
        $now = time();
        $start = $start - 86400;
        $sql = "INSERT INTO `".$table."` (`enterprise_id`,`mem_uid`,`day_time`,`order_price`,`customer_num`,`order_num`,`scan_num`,`created_at`) VALUES (".$this->enterpriseId.",'".$this->memUid."',".$start.",0,0,0,0,".$now.")";
        for($i = 1; $i <= 89; $i++) {
            $starts = $start - $i * 86400;
            $sql .= ",(" . $this->enterpriseId . ",'" . $this->memUid . "'," . $starts . ",0,0,0,0," . $now . ")";
        }
        DB::statement($sql);
    }

    public function checkTable()
    {
        //表名
        $table = Constant::TABLE_NAME.$this->enterpriseId;
        if (!Schema::hasTable($table)) {
            $sql = "CREATE TABLE `".$table."` (
            `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
              `enterprise_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '企业id',
              `mem_uid` char(32) NOT NULL DEFAULT '' COMMENT '导购id',
              `day_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '天数时间戳',
              `order_price` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '订单金额',
              `customer_num` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '客户数',
              `order_num` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '订单数',
              `scan_num` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '浏览人数',
              `created_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
              `updated_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
              `deleted_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '删除时间',
              PRIMARY KEY (`id`),
              KEY `enterprise_id_index` (`enterprise_id`)
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数据记录表';";

            DB::statement($sql);
        }

        return true;
    }

    public function trendInit($table)
    {
        $start = strtotime('today 00:00:00');
        $end = strtotime('today 23:59:59');
        $now = time();
        // 模拟数据
        for ($i=1; $i<=90; $i++) {
            $data = [];
            $starts = $start - $i * 86400;
            $ends = $end - $i * 86400;
            $userNum = $this->memberUserRelationService->timeUserCount(
                $this->enterpriseId,
                $this->memUid,
                $starts,
                $ends
            );
            $param = [
                'enterprise_id' => $this->enterpriseId,
                'mem_uid' => $this->memUid,
                'start_time' => $starts,
                'end_time' => $ends,
                'type' => 1
            ];
            $list = (new Order())->orderByTime($param);
            // 封装数据
            if (empty($list['order_price'])) {
                $list['order_price'] = 0;
            }
            if (empty($list['order_num'])) {
                $list['order_num'] = 0;
            }
            $data = [
                'enterprise_id' => $this->enterpriseId,
                'mem_uid' => $this->memUid,
                'day_time' => $starts,
                'order_price' => doubleval($list['order_price']),
                'customer_num' => $userNum,
                'order_num' => $list['order_num'],
                'scan_num' => $list['scan_person_num'],
                'created_at' => $now
            ];
            DB::table($table)->insert($data);
        }
    }
}
