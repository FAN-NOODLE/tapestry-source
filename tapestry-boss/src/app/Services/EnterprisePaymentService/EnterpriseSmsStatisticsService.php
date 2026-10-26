<?php


namespace App\Services\EnterprisePaymentService;


use App\Models\EnterpriseSmsStatistics;
use App\Services\BaseService;

class EnterpriseSmsStatisticsService extends BaseService
{

    /**
     * 扣减短信数量
     * @param $data
     * @return EnterpriseSmsStatistics|\Illuminate\Database\Eloquent\Builder|\Illuminate\Database\Eloquent\Model|object|void|null
     * @date 2020/11/6
     * @author 杨虎成 yanghucheng@vchangyi.com
     */
    public function depleteSms($data)
    {
        $smsStatistics = EnterpriseSmsStatistics::query()->where(['enterprise_id' => $data['enterprise_id']])->first();
        if (!$smsStatistics) {
            return error(400110);
        }
        if ($smsStatistics->sms_number < $data['number']) {
            return error(400111);
        }
        try {
            $smsStatistics->sms_number -= $data['number'];
            $smsStatistics->save();
            return $smsStatistics;
        } catch (\Exception $e) {
            return error(110000);
        }
    }

    /**
     * 获取短信数量
     * @param $data
     * @return EnterpriseSmsStatistics|bool|\Illuminate\Database\Eloquent\Builder|\Illuminate\Database\Eloquent\Model|object|null
     * @date 2020/11/6
     * @author 杨虎成 yanghucheng@vchangyi.com
     */
    public function getSms($data)
    {
        $smsStatistics = EnterpriseSmsStatistics::query()
            ->where(['enterprise_id' => $data['enterprise_id']])
            ->first();
        if (!$smsStatistics) {
            return false;
        }
        return $smsStatistics;
    }
}