<?php

namespace App\Console\Commands;

use App\Models\Enterprise;
use App\Models\EnterpriseAdmin;
use App\Models\EnterprisePayment;
use App\Services\EnterpriseLogService;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;


class EnterpriseAccountSync extends Command
{
    protected $signature = 'auto:enterprise_account_sync';

    protected $description = '企业付费账号同步';

    public function handle()
    {
        $this->failHandle();
        $this->effectiveHandle();

    }

    /**
     * 过期处理
     * @date 2020/11/5
     * @author 杨虎成 yanghucheng@vchangyi.com
     */
    public function failHandle()
    {
        $time = time();
        $list = EnterprisePayment::query()
            ->where('fail_time', '<=', $time)
            ->where('fail_status', 0)
            ->where('parent_id', 0)
            ->whereIn('payment_category', [EnterprisePayment::PAYMENT_CATEGORY_ZHDG, EnterprisePayment::PAYMENT_CATEGORY_SCRM, EnterprisePayment::PAYMENT_CATEGORY_MALL, EnterprisePayment::PAYMENT_CATEGORY_SYSTEM])
            ->select(['enterprise_id', 'buy_number', 'id'])
            ->get()
            ->toArray();

        if (empty($list)) {
            return;
        }
        // 处理数据
        $failList = [];
        $ids = [];
        foreach ($list as $value) {
            if (isset($failList[$value['enterprise_id']])) {
                $failList[$value['enterprise_id']] +=  $value['buy_number'];
            } else {
                $failList[$value['enterprise_id']] = $value['buy_number'];
            }

            $ids[] = $value['id'];
        }
        try {
            DB::beginTransaction();
            foreach ($failList as $enterpriseId => $value) {
                // 禁用管理员
                EnterpriseAdmin::query()
                    ->where('enterprise_id', $enterpriseId)
                    ->where('is_super', 2)
                    ->limit($value['total_number'])
                    ->update(['status' => 2]);
                // 更新子账号数量
                Enterprise::query()->where('id', $enterpriseId)->decrement('account_limit', $value);
            }

            // 更新付费记录为失效
            EnterprisePayment::query()
                ->whereIn('id', $ids)
                ->update(['fail_status' => 1]);
        } catch (\Exception $e) {
            DB::rollBack();
            Log::info('企业付费账号同步-失效账号处理异常：' . $e->getMessage(), $list);
            EnterpriseLogService::info('企业付费账号同步-失效账号处理异常', [
                'msg' => $e->getMessage(),
                'data' => $list,
            ]);
            return;
        }
        DB::commit();
        return;
    }

    /**
     * 生效处理
     * @date 2020/11/5
     * @author 杨虎成 yanghucheng@vchangyi.com
     */
    public function effectiveHandle()
    {
        $time = time();
        $list = EnterprisePayment::query()
            ->where('use_time', '<=', $time)
            ->where('account_sync', 0)
            ->where('fail_status', 0)
            ->where('parent_id', 0)
            ->whereIn('payment_category', [EnterprisePayment::PAYMENT_CATEGORY_ZHDG, EnterprisePayment::PAYMENT_CATEGORY_SCRM, EnterprisePayment::PAYMENT_CATEGORY_MALL, EnterprisePayment::PAYMENT_CATEGORY_SYSTEM])
            ->select(['enterprise_id', 'buy_number', 'id', 'fail_time'])
            ->orderBy('fail_time')
            ->get()
            ->toArray();

        if (empty($list)) {
            return;
        }

        // 处理数据
        $effectiveList = [];
        $ids = [];
        foreach ($list as $value) {
            if (isset($effectiveList[$value['enterprise_id']])) {
                $effectiveList[$value['enterprise_id']]['buy_number'] +=  $value['buy_number'];
                $effectiveList[$value['enterprise_id']]['fail_time'] =  $value['fail_time'];
            } else {
                $effectiveList[$value['enterprise_id']]['buy_number'] =  $value['buy_number'];
                $effectiveList[$value['enterprise_id']]['fail_time'] =  $value['fail_time'];
            }
            $ids[] = $value['id'];
        }
        try {
            DB::beginTransaction();
            foreach ($effectiveList as $enterpriseId => $value) {

                // 更新子账号数量
                $enterprise = Enterprise::query()
                    ->where('id', $enterpriseId)
                    ->get(['account_limit', 'account_expire_at'])
                    ->first();
                if ($enterprise) {
                    $enterprise = $enterprise->toArray();
                    $data = [
                        'account_limit' => ($enterprise['account_limit'] + $value['buy_number']),
                    ];
                    if ($enterprise['account_expire_at'] < $value['fail_time']) {
                        $data['account_expire_at'] = $value['fail_time'];
                    }
                    Enterprise::query()->where('id', $enterpriseId)->update($data);
                }

            }
            // 更新付费记录的同步状态
            EnterprisePayment::query()
                ->whereIn('id', $ids)
                ->update(['account_sync' => 1]);
        } catch (\Exception $e) {
            DB::rollBack();
            Log::info('企业付费账号同步-生效账号处理异常：' . $e->getMessage(), $list);
            EnterpriseLogService::info('企业付费账号同步-生效账号处理异常', [
                'msg' => $e->getMessage(),
                'data' => $list,
            ]);
            return;
        }
        DB::commit();
        return;
    }
}