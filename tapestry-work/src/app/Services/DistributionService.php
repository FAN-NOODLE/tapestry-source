<?php
namespace App\Services;
use App\Models\Distribution;

class DistributionService extends BaseService
{
    /**
     * 获取分销结算标准
     * @param array $params
     * @return array
     */
    public function getDistribution($params = [])
    {
        //获取分销信息
        $dataOld = Distribution::query()
            ->where('enterprise_id', $params['enterprise_id'])
            ->first(['settlement_standard']);

        return $dataOld ? $dataOld->toArray() : array();
    }

    /**
     * 分销配置详情
     *
     * @param int $enterpriseId
     * @param int $distributionId
     *
     * @return array
     * @throws \Exception
     */
    public function detail(int $enterpriseId, &$distributionId = 0)
    {
        $distribution = [];
        // 查询分销配置详情
        if ($distributionId == 0) {
            $detail = Distribution::query()
                ->where('enterprise_id', $enterpriseId)
                ->first(['id', 'distribution_title','distribution_text', 'settlement_standard', 'poster_img']);
        } else {
            $detail = Distribution::query()
                ->where('enterprise_id', $enterpriseId)
                ->where('id', $distributionId)
                ->first(['id', 'distribution_title','distribution_text', 'settlement_standard', 'poster_img']);
        }

        if (is_object($detail)) {
            $distribution   = $detail->toArray();
            $distributionId = $detail->id;
        }
        return $distribution;
    }

}