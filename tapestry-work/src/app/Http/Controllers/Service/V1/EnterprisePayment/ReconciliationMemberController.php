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
use App\Models\EnterprisePaymentRecord;
use App\Services\EnterprisePaymentService;
use App\Services\StaffManagerService;
use CrsSdk\Shop\User;

class ReconciliationMemberController extends BaseController
{
    /**
     * 活动已发送红包员工列表
     *
     * @return array
     */
    public function index()
    {
        $params = $this->request->only(['enterprise_id', 'enterprise_payment_id', 'qy_domain', 'app_identifier']);

        if (empty($params['enterprise_id'])) {
            throw new Error(110003, '企业ID不能为空');
        }

        $list = EnterprisePaymentRecord::query()
            ->where('enterprise_id', $params['enterprise_id'])
            ->when(!empty($params['enterprise_payment_id']), function ($query) use ($params) {
                $query->where('enterprise_payment_id', $params['enterprise_payment_id']);
            })
            ->get()
            ->toArray();

        $data = $this->format($list, $params);

        return success([
            'data' => $data,
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

        sort($memUids);
        $params['memUids'] = $memUids;

        $userList = (new StaffManagerService())->getStaffListFromUc($search['qy_domain'], $search['app_identifier'], $params);
       
        $userData = array_column($userList['list'], null, 'memUid') ?? [];

        $list = [];
        foreach ($memUids as $memUid) {
            $memUser = $userData[$memUid] ?? [];
            $list[]  = [
                'mem_username' => $memUser['memUsername'] ?? '',
                'mem_face'     => $memUser['memFace'] ?? '',
                'mem_uid'      => $memUid ?? 0,
            ];
        }

        return $list;
    }
}
