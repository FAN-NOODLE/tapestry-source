<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/12/3
 * Time: 18:33
 */

namespace App\Services;

use App\Models\EnterprisePaymentBudget;

class EnterprisePaymentBudgetService
{
    // 查询当前企业余额相关
    public function getBudgetInfo($enterpriseId)
    {

        $budgetInfo = EnterprisePaymentBudget::query()
            ->where('enterprise_id', $enterpriseId)
            ->orderByDesc('id')
            ->first();

        return $budgetInfo;
    }



}
