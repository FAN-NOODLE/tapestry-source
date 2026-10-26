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
use App\Library\Uc\Member;
use App\Models\EnterprisePayment;
use App\Models\EnterprisePaymentBudget;
use App\Models\EnterprisePaymentRecord;
use App\Models\MemberUserRelation;
use App\Services\EnterprisePaymentBudgetService;
use App\Services\EnterprisePaymentService;
use App\Services\StaffManagerService;
use CrsSdk\Shop\User;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;

class ReconciliationListController extends BaseController
{
    /**
     * 活动对账记录
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

        $list = EnterprisePaymentRecord::query()
            ->where('enterprise_id', $search['enterprise_id'])
            ->when(!empty($search['start_time']), function ($query) use ($search) {
                $query->where('send_time', '>=', $search['start_time']);
            })
            ->when(!empty($search['enterprise_payment_id']), function ($query) use ($search) {
                $query->where('enterprise_payment_id', $search['enterprise_payment_id']);
            })
            ->when(!empty($search['end_time']), function ($query) use ($search) {
                $query->where('send_time', '<=', $search['end_time']);
            })
            ->when(!empty($search['mem_uid']), function ($query) use ($search) {
                $query->where('mem_uid', $search['mem_uid']);
            })
            ->orderByDesc('id')
            ->paginate($params['size'] ?? 10)
            ->toArray();

        $data = $this->format($list['data'], $search);

        return success([
            'data' => $data,
            'meta' => [
                'total' => $list['total'],
                'size'  => (int)$list['per_page'],
                'page'  => $list['current_page']
            ]
        ]);
    }

    /**
     * 格式化代码
     *
     * @param array $data
     * @param array $search
     *
     * @return array
     */
    public function format($data = [], $search = [])
    {

        if (empty($data)) {
            return $data;
        }

        $memUids = array_unique(array_column($data, 'mem_uid'));
        $unionIds = array_unique(array_column($data, 'union_id'));

        sort($memUids);
        $params['memUids'] = $memUids;

        $userList = (new StaffManagerService())->getStaffListFromUc($search['qy_domain'], $search['app_identifier'], $params);
        $userData = array_column($userList['list'], null, 'memUid') ?? [];

        // 获取用户信息
        $userList = MemberUserRelation::query()
                    ->whereIn('union_id', $unionIds)
                    ->get(['ec_name', 'ec_avatar', 'union_id'])
                    ->toArray();
        $userList = array_column($userList, null, 'union_id');

        $list = [];
        foreach ($data as $value) {
            $memUser = $userData[$value['mem_uid']] ?? [];
            $user    = $userList[$value['union_id']] ?? [];
            $list[]  = [
                'mem_username' => $memUser['memUsername'] ?? '',
                'mem_face'     => $memUser['memFace'] ?? '',
                'send_time'    => $value['send_time'] ?? 0,
                'nick_name'    => $user['ec_name'] ?? '',
                'avatar'       => $user['ec_avatar'] ?? '',
                'status'       => $value['status'],
                'amount'       => $value['amount']
            ];
        }

        return $list;
    }
}
