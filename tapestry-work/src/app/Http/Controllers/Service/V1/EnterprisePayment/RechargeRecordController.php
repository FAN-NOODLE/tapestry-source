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
use App\Models\EnterprisePaymentBudget;
use Illuminate\Http\Request;

class RechargeRecordController extends BaseController
{
    /**
     * 充值记录接口
     *
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request)
    {
        $params = $request->only(['page', 'size']);
        $search = $request->input('search', []);

        if (empty($search['enterprise_id'])) {
            throw new Error(110003, '企业ID不能为空');
        }

        $list = EnterprisePaymentBudget::query()
            ->where('enterprise_id', $search['enterprise_id'])
            ->when(!empty($search['type']), function ($query) use ($search) {
                $query->where('type', $search['type']);
            })
            ->when(!empty($search['start_time']), function ($query) use ($search) {
                $query->where('created_at', '>=', $search['start_time']);
            })
            ->when(!empty($search['end_time']), function ($query) use ($search) {
                $query->where('created_at', '<=', $search['end_time']);
            })
            ->orderByDesc('id')
            ->paginate($params['size'] ?? 10)
            ->toArray();

        return success([
            'data' => $list['data'],
            'meta' => [
                'total' => $list['total'],
                'size'  => (int)$list['per_page'],
                'page'  => $list['current_page']
            ]
        ]);
    }
}
