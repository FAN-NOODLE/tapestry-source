<?php
/**
 * Created by PhpStorm.
 * User: 大熊
 * Date: 2018/12/3
 * Time: 16:03
 */

namespace App\Services;

use App\Models\CommissionRecord;
use Illuminate\Support\Facades\DB;

class CommissionRecordService
{

    /**
     * @var CommissionRecord
     */
    protected $commissionRecord;

    public function __construct(CommissionRecord $commissionRecord)
    {
        $this->commissionRecord = $commissionRecord;
    }

    /**
     * 获取导购某段时间内佣金数
     *
     * @param int    $enterpriseId 企业ID
     * @param string $memUid       导购ID
     * @param int    $startTime    开始时间
     * @param int    $endTime      结束时间
     *
     * @return int
     */
    public function timeCommissionMoney($enterpriseId, $memUid, $startTime = 0, $endTime = 0)
    {
        // 查询条件
        $condition = [
            ['enterprise_id', '=', $enterpriseId],
            ['mem_uid', '=', $memUid],
        ];

        if ($startTime) {
            $condition[] = ['created_at', '>=', $startTime];
        }

        if ($endTime) {
            $condition[] = ['created_at', '<=', $endTime];
        }

        $money = $this->commissionRecord
            ->query()
            ->select(DB::raw('sum(charge_price) as money'))
            ->where($condition)
            ->first();

        if (isset($money['money']) && $money['money']) {
            $money = $money['money'];
        } else {
            $money = 0;
        }
        return $money;
    }
}