<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/12/3
 * Time: 3:12 PM
 */

namespace App\Http\Controllers\Service\V1\EnterprisePayment;

use App\Http\Controllers\Service\V1\BaseController;
use App\Http\Requests\EnterprisePaymentRequest;
use App\Exceptions\Error;
use App\Models\EnterprisePayment;
use App\Models\EnterprisePaymentBudget;
use App\Services\EnterprisePaymentBudgetService;
use App\Services\EnterprisePaymentService;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;

class ListController extends BaseController
{
    /**
     * 企业微信红包列表
     *
     * @return array
     */
    public function index()
    {
        $params = $this->request->only(['page', 'size']);
        $search = $this->request->input('search', []);

        if (empty($search['enterprise_id'])) {
            throw new Error(110003, '企业ID不能为空');
        }

        $list = EnterprisePayment::query()
            ->where('enterprise_id', $search['enterprise_id'])
            ->where('is_deleted', $search['is_deleted'])
            ->when(!empty($search['start_time']), function ($query) use ($search) {
                $query->where('created_at', '>=', $search['start_time']);
            })
            ->when(!empty($search['end_time']), function ($query) use ($search) {
                $query->where('created_at', '<=', $search['end_time']);
            })
            ->when(!empty($search['is_sale']), function ($query) use ($search) {
                $query->where('is_sale', $search['is_sale']);
            })
            ->when(!empty($search['number']), function ($query) use ($search) {
                $query->where('number', 'like', '%'.$search['number'].'%');
            })
            ->when(!empty($search['name']), function ($query) use ($search) {
                $query->where('name', 'like', '%'.$search['name'].'%');
            })
            ->orderByDesc('id')
            ->paginate($params['size'] ?? 10)
            ->toArray();

        // 初始化类
        $service = (new EnterprisePaymentService());

        $data = [];
        foreach ($list['data'] as $value) {
            // 获取列表数量
            $statisticalNum = $service->getStatisticalNum($value);
            $value          = array_merge($value, $statisticalNum);
            $value['path']  = 'other/red-envelope/red-envelope?a='.$value['id'];
            $data[]         = $value;
        }

        return success([
            'data' => $data,
            'meta' => [
                'total' => $list['total'],
                'size'  => (int)$list['per_page'],
                'page'  => $list['current_page']
            ]
        ]);
    }
}
