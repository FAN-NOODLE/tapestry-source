<?php
/**
 * Created by PhpStorm.
 * User: XuHuitao
 * Date: 2020/7/13
 * Time: 17:04
 */

namespace App\Models;

/**
 * Class MemberOrderStatistics
 *
 * @package App\Models
 * @property $id
 * @property $union_id
 * @property $total_order_num
 * @property $total_order_amount
 * @property $last_update_data
 * @property $enterprise_id
 * @property $latest_buy_time
 */
class MemberOrderStatistics extends Base {
    /**
     * 表名
     *
     * @var string
     */
    protected $table = 'member_order_statistics';

    /**
     * @var array
     */
    protected $fillable = [
        'id',
        'union_id',
        'user_id',
        'total_order_num',
        'total_order_amount',
        'last_update_data',
        'enterprise_id',
        'latest_buy_time'
    ];

    /**
     * getDataByUserId
     * @desc 通过userId查询用户订单数据，返回model对象
     * @param $userId
     *
     * @return \Illuminate\Database\Eloquent\Model|null|object|static
     * @date 2020-07-13
     * @author XuHuitao
     */
    public static function getDataByUserId($userId, $enterpriseId, $unionId = '') {
        if (empty($userId) && empty($unionId)) {
            return [];
        }
        $query = self::query()->where('enterprise_id', $enterpriseId);
        if (!empty($userId)) {
            $query->where('user_id', $userId);
        }
        if (!empty($unionId)) {
            $query->where('union_id', $unionId);
        }
        return $query->first();
    }

    /**
     * createUserOrder
     * @desc 在数据库新增数据
     * @param        $userId
     * @param        $unionId
     * @param        $totalOrder
     * @param        $totalAmount
     * @param string $orderNo
     *
     * @return $this|\Illuminate\Database\Eloquent\Model
     * @date 2020-07-13
     * @author XuHuitao
     */
    public static function createUserOrder($userId, $unionId, $totalOrder, $totalAmount, $orderNo, $enterpriseId, $createdAt) {
        $data = [
            'enterprise_id' => $enterpriseId,
            'user_id' => $userId,
            'union_id' => $unionId,
            'total_order_num' => $totalOrder,
            'total_order_amount' => $totalAmount,
            'last_update_data' => json_encode([$orderNo]),
            'latest_buy_time' => $createdAt,
        ];
        return self::query()->create($data);
    }

    /**
     * isRecord
     * @desc 判断订单是否已经记录过
     * @param $orderNo
     *
     * @return bool
     * @date 2020-07-25
     */
    public function isRecord($orderNo) {
        $history = json_decode($this->last_update_data, true);
        $history = $history ?? [];
        return in_array($orderNo, $history);
    }

    /**
     * addOrder
     * @desc 增加订单操作
     * @param $orderNo
     * @param $orderAmount
     *
     * @return bool
     * @date 2020-07-25
     */
    public function addOrder($orderNo, $orderAmount, $createOrderAt) {
        $this->total_order_num += 1;
        $this->total_order_amount += $orderAmount;
        $this->latest_buy_time = $createOrderAt > $this->latest_buy_time ? $createOrderAt : $this->latest_buy_time;
        $history = json_decode($this->last_update_data, true);
        $history = $history ?? [];
        $history[] = $orderNo;
        $this->last_update_data = json_encode($history);
        return $this->save();
    }
}