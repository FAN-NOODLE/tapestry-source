<?php
/**
 * Created by sublime.
 * User: zhaoqi
 * Date: 2019/01/17
 * Time: 10:00
 */
namespace App\Http\Controllers\V1\Label\System;

use CrsSdk\Shop\SystemLabel;
use Illuminate\Http\Request;
use App\Http\Controllers\App\V1\Label\LabelBaseController;

class UpdateController extends LabelBaseController
{
    /**
     * 系统标签更新
     *
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request)
    {
        $params = $request->only([
            "label_id",
            "name",
            "condition",
            "select_condition",
            "goods_type",
            "time_type",
            "date_min",
            "date_max",
            "number_min",
            "number_max",
            "price_min",
            "price_max",
            "score_min",
            "score_max",
            "average_cost_min",
            "average_cost_max",
            "select_goods"
        ]);

        $re = $this->validates($params);

        if ($re != 0) {
            error($re);
        }

        $params['enterprise_id'] = $this->enterpriseId;

        $result = (new SystemLabel())->update($params);

        if (!$result) {
            notice(100000, [], $params);
        }

        // 日志记录
        $this->writeOperatorLog('更新了 ' . count($params['name']) . ' 系统标签');

        return $result;
    }

    private function validates($params)
    {
        // 验证标签 ID
        if (empty($params['label_id']) || !is_numeric($params['label_id']) || $params['label_id'] < 1) {
            return 800000;
        }

        // 标签名称不能为空
        if (empty(preg_replace('# #','',$params['name']))) {
            return 800001;
        }

        // 标签名称介于0~6个字符
        if (mb_strlen($params['name']) < 0 || mb_strlen($params['name']) > 7) {
            return 800002;
        }

        // 条件类型必须为整数且只能是1或者2 1:满足其中之一，2:都满足
        if (!is_int($params['condition']) || !in_array($params['condition'], [1,2])) {
            return 800003;
        }

        // 条件集合为数组,不能为空
        if (!is_array($params['select_condition']) || empty($params['select_condition'])) {
            return 800004;
        }

        // 验证条件数据合法性
        $select_condition = array_filter($params['select_condition'], function ($value) {
            return is_int($value) ? true : false;
        });
        if (array_diff($params['select_condition'], $select_condition)) {
            return 800004;
        }

        // 验证条件合法性 1:上次消费时间,2:累积消费次数,3:累积消费金额,4:累积获得积分,5:平均消费金额,6:购买任意商品,7:购买指定商品
        $select_condition = array_filter($params['select_condition'], function ($value) {
            return in_array($value, [1,2,3,4,5,6,7]) ? true : false;
        });

        if (array_diff($params['select_condition'], $select_condition)) {
            return 800004;
        }

        // 时间条件类型 0:指定时间区间查询label_condition表,其他大于0的值代表最近天数
        if (in_array(1, $params['select_condition']) && !is_int($params['time_type'])) {
            return 800005;
        }

        // 时间条件类型 0:指定时间区间查询label_condition表,其他大于0的值代表最近天数
        if (in_array(1, $params['select_condition']) && $params['time_type'] == 0) {
            if (empty($params['date_min']) || empty($params['date_max'])) {
                return 800005;
            }
        }

        // 当选择累积消费次数 
        if (in_array(2, $params['select_condition'])) {
            if (!is_int($params['number_min']) || !is_int($params['number_max'])) {
                return 800006;
            }

            if ($params['number_min'] < 0 || $params['number_max'] > 100000) {
                return 800006;
            }

            if ($params['number_min'] > $params['number_max']) {
                return 800006;
            }
        }

        // 当选择累积消费金额
        if (in_array(3, $params['select_condition'])) {
            if (!is_int($params['price_min']) || !is_int($params['price_max'])) {
                return 800007;
            }

            if ($params['price_min'] < 0 || $params['price_max'] > 10000000) {
                return 800007;
            }

            if ($params['price_min'] > $params['price_max']) {
                return 800007;
            }
        }

        // 当选择累积获得积分
/*        if (in_array(4, $params['select_condition'])) {
            if (!is_int($params['score_min']) || !is_int($params['score_max'])) {
                return 800008;
            }
        }*/

        // 当选择平均消费金额
        if (in_array(5, $params['select_condition'])) {
            if (!is_int($params['average_cost_min']) || !is_int($params['average_cost_max'])) {
                return 800009;
            }

            if ($params['average_cost_min'] < 0 || $params['average_cost_max'] > 10000000) {
                return 800009;
            }

            if ($params['average_cost_min'] > $params['average_cost_max']) {
                return 800009;
            }
        }

        // 当选择购买任意商品,商品条件类型（1:所有商品,2:指定商品查询label_good表）
        if (in_array(6, $params['select_condition'])) {
            if ($params['goods_type'] != 1) {
                return 800010;
            }
        }

        // 当选择购买指定商品
        if (in_array(7, $params['select_condition'])) {
            if ($params['goods_type'] != 2) {
                return 800010;
            }

            if (!is_array($params['select_goods']) || empty($params['select_goods'])) {
                return 800010;
            }

            $select_good = array_filter($params['select_goods'], function ($value) {
                return is_int($value) ? true : false;
            });

            if (array_diff($params['select_goods'], $select_good)) {
                return 800010;
            }
        }

        return 0;
    }
}
