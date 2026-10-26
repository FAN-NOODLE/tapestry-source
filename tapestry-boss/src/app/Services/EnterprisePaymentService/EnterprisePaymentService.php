<?php


namespace App\Services\EnterprisePaymentService;


use App\Jobs\SendEmailJob;
use App\Models\Enterprise;
use App\Models\EnterpriseAdmin;
use App\Models\EnterprisePayment;
use App\Models\EnterprisePaymentLog;
use App\Models\EnterpriseSmsStatistics;
use App\Services\BaseService;
use App\Services\EnterpriseNew\EnterpriseNewService;
use App\Services\EnterpriseService;
use Illuminate\Support\Facades\DB;
use CrsSdk\Work\EnterprisePayment AS SDKEnterprisePayment;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Redis;

class EnterprisePaymentService extends BaseService
{

    protected $fields = [
        'id',
        'enterprise_id',
        'parent_id',
        'payment_category',
        'payment_amount',
        'payment_source',
        'collection_amount',
        'available_amount',
        'server_charge',
        'buy_number',
        'buy_day',
        'payment_time',
        'arrival_time',
        'use_time',
        'fail_time',
        'fail_status',
        'payment_type',
        'payment_status',
        'payment_order_sn',
        'remark'
    ];
    /**
     * 添加企业付费信息
     * @param $data
     * @return bool|int
     * @date 2020/11/5
     * @author 杨虎成 yanghucheng@vchangyi.com
     */
    public function add($data)
    {
        // 企业信息查询
        $enterprise = (new EnterpriseService)->enterprise($data['enterprise_id']);
        if (!$enterprise) {
            return 201002;
        }

        $data['account_sync'] = 0;
        // 2：智慧导购，3：SCRM，4：商城，5：综合系统
        if (in_array($data['payment_category'], [EnterprisePayment::PAYMENT_CATEGORY_ZHDG, EnterprisePayment::PAYMENT_CATEGORY_SCRM, EnterprisePayment::PAYMENT_CATEGORY_MALL, EnterprisePayment::PAYMENT_CATEGORY_SYSTEM])) {

            // 验证时间
            $payTime = date('Ymd', $data['payment_time']);
            if ($data['use_time'] < $payTime) {
                return 400101;
            }
            $time = time();
            // 生效时间到达，同步子账号
            if ($data['use_time'] <= $time) {
                $data['account_sync'] = 1;
            }
            // 计算失效时间 智慧导购除外
            if ($data['payment_category'] != EnterprisePayment::PAYMENT_CATEGORY_ZHDG) {
                $data['fail_time'] = $data['use_time'] + ($data['buy_day'] * 86400) - 1;
            }

            $data['fail_status'] = 0;
            // 过期时间小于当前时间，直接过期
            if ($data['fail_time'] <= $time) {
                $data['fail_status'] = 1;
            }

        }
        // 到账判断
        if ($data['payment_status'] == 1) {
            $data['arrival_time'] = time();
        }

        // 组织日志数据
        $logData = [
            'new_data_json' => json_encode($data, JSON_UNESCAPED_UNICODE),
            'enterprise_id' => $data['enterprise_id']
        ];

        // 入库
        try {
            DB::beginTransaction();

            $enterprisePayment = EnterprisePayment::query()->create($data);

            $logData['enterprise_payment_id'] = $enterprisePayment->id;
            EnterprisePaymentLog::query()->create($logData);

            // 购买系统变更 1、试用版企业改为正式版，2、生效时间为当前时间增加子账号数量
            if (in_array($enterprisePayment->payment_category, [EnterprisePayment::PAYMENT_CATEGORY_ZHDG, EnterprisePayment::PAYMENT_CATEGORY_SCRM, EnterprisePayment::PAYMENT_CATEGORY_MALL, EnterprisePayment::PAYMENT_CATEGORY_SYSTEM])) {

                $enterpriseData = [];
                if ($enterprise->version_tag == '试用版') {
                    $enterpriseData['version_tag'] = '正式版';
                    $enterprise->account_limit = 0;
                    $enterprise->account_expire_at = 0;
                    $enterprise->version_tag = '正式版';
                    $emailData = $enterprise->toArray();
                    $emailData['account_expire_at'] = $data['fail_time'];
                    $emailData['account_limit'] = $data['buy_number'];
                    $emailData['emailRemark'] = "请使用当前密码登录";
                    (new SendEmailJob($emailData, 0))->handle();

                }
                if ($data['account_sync'] == 1 && $data['fail_status'] == 0) {
                    $enterpriseData['account_limit'] = $enterprise->account_limit + $data['buy_number'];
                    if ($data['fail_time'] > $enterprise->account_expire_at) {
                        $enterpriseData['account_expire_at'] = $data['fail_time'];
                    }
                }
                if (!empty($enterpriseData)) {
                    Enterprise::query()->where('id', $enterprisePayment->enterprise_id)->update($enterpriseData);
                }
            }

            // 购买红包通知
            if ($enterprisePayment->payment_category == EnterprisePayment::PAYMENT_CATEGORY_RED_PACK  && $enterprisePayment->available_amount > 0) {
                $requestParams = [
                    'enterprise_id' => $enterprisePayment->enterprise_id,
                    'amount' => $enterprisePayment->available_amount, // 传递红包可用金额
                    'voucher_url' => $enterprisePayment->voucher_url ?? '',
                    'income_type' => 1,
                ];
                $requestResult = (new SDKEnterprisePayment)->recharge($requestParams);
                Log::info('购买红包通知：', ['params' => $requestParams, 'result' => $requestResult]);
            }

            // 短信记录统计
            if ($enterprisePayment->payment_category == EnterprisePayment::PAYMENT_CATEGORY_SMS) {
                $smsStatistics = EnterpriseSmsStatistics::query()->where(['enterprise_id' => $enterprisePayment->enterprise_id])->first();
                if (!$smsStatistics) {
                    $smsData = [
                        'enterprise_id' => $enterprisePayment->enterprise_id,
                        'sms_number' => $enterprisePayment->buy_number,
                        'sms_total_number' => $enterprisePayment->buy_number,
                        'sms_payment_total' => $enterprisePayment->payment_amount,
                    ];
                    EnterpriseSmsStatistics::query()->create($smsData);
                } else {
                    $smsData = [
                        'sms_number' => ($smsStatistics->sms_number + $enterprisePayment->buy_number),
                        'sms_total_number' => ($smsStatistics->sms_total_number + $enterprisePayment->buy_number),
                        'sms_payment_total' => ($smsStatistics->sms_payment_total + $enterprisePayment->payment_amount),
                    ];
                    EnterpriseSmsStatistics::query()->where('enterprise_id', $enterprisePayment->enterprise_id)->update($smsData);

                }
            }
        } catch (\Exception $e) {
            DB::rollBack();
            Log::info('创建失败：', ['data' => $data, 'msg' => $e->getMessage()]);
            return 400102;
        }
        DB::commit();

        $redisKey = redisKey('REDIS_ENTERPRISE_PAYMENT_HAS_NEW', ['enterprise_id' => $enterprisePayment->enterprise_id]);
        $redis = Redis::connection();
        $redis->setnx($redisKey, 1);

        return true;
    }


    /**
     * 统计企业付费信息
     * @param $data
     * @return array
     * @date 2020/11/5
     * @author 杨虎成 yanghucheng@vchangyi.com
     */
    public function statistics($data)
    {
        $startTime = strtotime(date('Y0101'));
        $endTime = strtotime(date('Y1231')) + 86399;

        $result = [
            'all_money' => 0,
            'all_year_money' => 0,
            'all_system_money' => 0,
            'all_year_system_money' => 0,
            'all_account' => 0,
            'seven_account' => 0,
            'fail_time' => 0,
            'fail_day' => 0,
            'red_pack_total_money' => 0,
            'enterprise_payment_id' => 0,
        ];

        // 全部累计
        $allPayment = EnterprisePayment::query()
            ->where(['enterprise_id' => $data['enterprise_id']])
            ->select(['payment_category', DB::raw('SUM(payment_amount) AS money')])
            ->groupBy('payment_category')
            ->get()
            ->toArray();
        if ($allPayment) {
            foreach ($allPayment as $value) {
                $result['all_money'] += $value['money'];
                if ($value['payment_category'] == EnterprisePayment::PAYMENT_CATEGORY_RED_PACK) {
                    $result['red_pack_total_money'] += $value['money'];
                }
                // 系统累计
                if (!in_array($value['payment_category'], [EnterprisePayment::PAYMENT_CATEGORY_RED_PACK, EnterprisePayment::PAYMENT_CATEGORY_SMS])) {
                    $result['all_system_money'] += $value['money'];
                }
            }
        }
        // 当年全部累计
        $allYearPayment = EnterprisePayment::query()
            ->where(['enterprise_id' => $data['enterprise_id']])
            ->whereBetween('payment_time', [$startTime, $endTime])
            ->select(['payment_category', DB::raw('SUM(payment_amount) AS money')])
            ->groupBy('payment_category')
            ->get()
            ->toArray();
        if ($allYearPayment) {
            foreach ($allYearPayment as $value) {
                $result['all_year_money'] += $value['money'];
                // 本年系统累计
                if (!in_array($value['payment_category'], [EnterprisePayment::PAYMENT_CATEGORY_RED_PACK, EnterprisePayment::PAYMENT_CATEGORY_SMS])) {
                    $result['all_year_system_money'] += $value['money'];
                }
            }
        }

        $account = EnterprisePayment::query()
            ->where('enterprise_id', $data['enterprise_id'])
            ->where('fail_status', 0)
            ->where('parent_id', 0)
            ->whereBetween('payment_category', [EnterprisePayment::PAYMENT_CATEGORY_ZHDG,  EnterprisePayment::PAYMENT_CATEGORY_SYSTEM])
            ->orderBy('fail_time')
            ->get(['buy_number', 'fail_time', 'id', 'payment_category'])
            ->toArray();
        $enterprise = (new EnterpriseNewService())->getEnterpriseDetail($data['enterprise_id'], true);
        $result['all_account'] = $enterprise['account_limit'];

        if ($account) {
            $systemCount = count($account);
            $failTime = $account[0]['fail_time'];
            $result['seven_account'] = $account[0]['buy_number'];
            $result['fail_time'] = $failTime; // 时间
            $result['fail_day'] = ceil(($failTime - time()) / 86400);
            if ($systemCount == 1 && $account[0]['payment_category'] != EnterprisePayment::PAYMENT_CATEGORY_ZHDG) {
                $result['enterprise_payment_id'] = $account[0]['id'];
            }
        }
        return $result;
    }

    /**
     * 企业付费信息列表|续费列表
     * @param $params
     * @param $page
     * @param $perPage
     * @param array $orderBy
     * @return EnterprisePayment|\Illuminate\Database\Eloquent\Builder|mixed
     * @date 2020/11/5
     * @author 杨虎成 yanghucheng@vchangyi.com
     */
    public function getList($params, $page, $perPage)
    {
        $list = EnterprisePayment::query()
            ->where('enterprise_id', $params['enterprise_id'])
            ->when(isset($params['payment_status']) && in_array($params['payment_status'], [EnterprisePayment::PAYMENT_STATUS_NO, EnterprisePayment::PAYMENT_STATUS_YES]) && is_int($params['payment_status']), function ($query) use ($params) {
                return $query->where('payment_status', $params['payment_status']);
            });

        if (!empty($params['start_time']) && !empty($params['end_time'])) {
            $field = 'payment_time';
            if ($params['time_type'] == 1) {
                $field = 'arrival_time';
            }
            $list = $list->whereBetween($field, [$params['start_time'], $params['end_time']]);
        }
        if (!empty($params['renew'])) {
            $list = $list->whereIn('payment_category', [EnterprisePayment::PAYMENT_CATEGORY_ZHDG, EnterprisePayment::PAYMENT_CATEGORY_SCRM, EnterprisePayment::PAYMENT_CATEGORY_MALL, EnterprisePayment::PAYMENT_CATEGORY_SYSTEM])
                ->where('parent_id', 0)
                ->orderBy('fail_status', 'asc')
                ->orderBy('fail_time');
        } else {
            $list = $list->orderBy('id', 'desc');
        }

        $list = $list
            ->paginate($perPage, $this->fields, 'page', $page)
            ->toArray();

        return $list;
    }

    /**
     * 企业付费信息详情
     * @param $params
     * @return array|bool
     * @date 2020/11/5
     * @author 杨虎成 yanghucheng@vchangyi.com
     */
    public function detail($params)
    {
        $detail = EnterprisePayment::query()
            ->where('id', $params['id'])
            ->get($this->fields)
            ->first();
        if (!$detail) {
            return false;
        }
        return $detail;
    }

    /**
     * 企业付费收款
     * @param $params
     * @return array|bool|int
     * @date 2020/11/5
     * @author 杨虎成 yanghucheng@vchangyi.com
     */
    public function receivePayment($data)
    {
        $detail = EnterprisePayment::query()
            ->where('id', $data['id'])
            ->where('payment_status', 0)
            ->exists();
        if (!$detail) {
            return 400404;
        }
        $id = $data['id'];
        unset($data['id']);
        $data['payment_status'] = 1;
        $data['arrival_time'] = time();

        EnterprisePayment::query()->where('id', $id)->update($data);
        return true;
    }

    /**
     * 企业付费续费
     * @param $data
     * @return bool|int
     * @date 2020/11/5
     * @author 杨虎成 yanghucheng@vchangyi.com
     */
    public function renew($data)
    {
        $detail = $this->detail($data);
        if (!$detail) {
            return 400404;
        }

        // 企业信息查询
        $enterprise = (new EnterpriseService)->enterprise($detail->enterprise_id);
        if (!$enterprise) {
            return 201002;
        }
        $id = $data['id'];
        unset($data['id']);

        $data['parent_id'] = $id;
        // 到账判断
        $data['arrival_time'] = 0;
        if ($data['payment_status'] == 1) {
            $data['arrival_time'] = time();
        }

        $failTime = $detail->fail_time + ($data['buy_day'] * 86400);
        $data['fail_time'] = $failTime;
        $data['enterprise_id'] = $detail->enterprise_id;
        $data['payment_category'] = $detail->payment_category;
        $data['use_time'] = $detail->fail_time;
        $data['buy_number'] = $detail->buy_number;
        $data['account_sync'] = 1; // 续费不同步账号
        $data['payment_source'] = '手动创建';


        // 组织日志数据
        $logData = [
            'old_data_json' => json_encode($detail->toArray(), JSON_UNESCAPED_UNICODE),
            'new_data_json' => json_encode($data, JSON_UNESCAPED_UNICODE),
            'enterprise_id' => $detail->enterprise_id,
            'enterprise_payment_id' => $id,
        ];
        try {
            DB::beginTransaction();
            // 新增付费记录
            EnterprisePayment::query()->create($data);
            // 更新数据
            EnterprisePayment::query()->where('id', $id)->update(['fail_time' => $failTime]);
            // 记录日志
            EnterprisePaymentLog::query()->create($logData);

            // 更新企业失效时间
            if ($enterprise->account_expire_at < $failTime) {
                Enterprise::query()->where('id', $detail->enterprise_id)->update(['account_expire_at' => $failTime]);

            }

        } catch (\Exception $e) {
            DB::rollBack();
            Log::info('更新失败：', ['data' => $data, 'msg' => $e->getMessage()]);
            return 400103;
        }
        DB::commit();
        return true;

    }

    /**
     * 获取导出数据
     * @param array $params
     * @return EnterprisePayment|array|\Illuminate\Database\Eloquent\Builder|mixed
     * @date 2020/11/6
     * @author 杨虎成 yanghucheng@vchangyi.com
     */
    public function exportList($params = [])
    {
        $list = EnterprisePayment::query()
            ->where('enterprise_id', $params['enterprise_id'])
            ->when(!empty($params['payment_status']), function ($query) use ($params) {
                return $query->where('payment_status', $params['payment_status']);
            });
        if (!empty($params['start_time']) && !empty($params['end_time'])) {
            $list = $list->whereBetween('payment_time', [$params['start_time'], $params['end_time']]);
        }
        $list = $list->orderBy('id', 'DESC')
            ->get($this->fields)->toArray();
        return $list;
    }

    /**
     * 智慧导购过期处理
     * @param array $params
     * @return bool|int
     * @date 2020/11/6
     * @author 杨虎成 yanghucheng@vchangyi.com
     */
    public function failAccount($params = [])
    {
        $row = EnterprisePayment::query()
            ->where('enterprise_id', $params['enterprise_id'])
            ->where('payment_order_sn', $params['payment_order_sn'])
            ->get($this->fields)
            ->toArray();
        if (empty($row)) {
            return 400404;
        }
        try {
            DB::beginTransaction();
            // 更新记录为失效
            EnterprisePayment::query()->where('id', $row['id'])->update(['fail_status' => 1]);
            // 更新子账号数量
            Enterprise::query()->where('id', $row['enterprise_id'])->decrement('account_limit', $row['buy_number']);
            // 禁用管理员
            EnterpriseAdmin::query()
                ->where('enterprise_id', $row['enterprise_id'])
                ->where('is_super', 2)
                ->limit($row['total_number'])
                ->update(['status' => 2]);
        } catch (\Exception $e) {
            DB::rollBack();
            Log::info('智慧导购付费账号同步-失效账号处理异常：' . $e->getMessage(), $params);
            return 400103;
        }
        DB::commit();
        return true;
    }

}