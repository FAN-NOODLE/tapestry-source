<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/12/3
 * Time: 3:12 PM
 */

namespace App\Http\Controllers\Service\V1\EnterprisePayment;

use App\Http\Controllers\Service\V1\BaseController;
use App\Exceptions\Error;
use App\Models\EnterprisePayment;
use App\Models\EnterprisePaymentBudget;
use Illuminate\Support\Facades\DB;

class BudgetByEnterpriseController extends BaseController
{
    /**
     * 企业充值对账
     */
    public function index()
    {
        $enterpriseIds = $this->request->input('enterprise_ids', []);

        if (empty($enterpriseIds)) {
            throw new Error(110003, '企业ID不能为空');
        }

        $ids = EnterprisePaymentBudget::query()
            ->select(DB::raw('max(id) as id'))
            ->whereIn('enterprise_id', $enterpriseIds)
            ->groupBy('enterprise_id')
            ->orderBy('id')
            ->get()
            ->toArray();
        $ids = array_column($ids, 'id');

        $result = EnterprisePaymentBudget::query()
            ->whereIn('id', $ids)
            ->get(['id', 'enterprise_id', 'balance'])
            ->toArray();

        return success([
            'data' => $result,
        ]);
    }
}
