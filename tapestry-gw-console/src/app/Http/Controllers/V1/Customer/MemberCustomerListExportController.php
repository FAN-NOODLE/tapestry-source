<?php
/**
 * Created by PhpStorm.
 * User: shizhixiong
 * Date: 2019-01-25
 * Time: 10:16
 */

namespace App\Http\Controllers\V1\Customer;

use App\Http\Controllers\V1\BaseController;
use App\Services\CustomerService;
use App\Services\MemberService;
use CrsSdk\Shop\Customer;
use CrsSdk\Shop\UserLevel;
use Maatwebsite\Excel\Facades\Excel;

class MemberCustomerListExportController extends BaseController
{
    /**
     * 客户确认状态：确认
     */
    const NO_TO_BE_CONFIRMED = 0;

    /**
     * 客户确认状态：待确认
     */
    const TO_BE_CONFIRMED = 1;

    /**
     * 客户类型
     */
    const CUSTOMER_TYPE_TEXT = [
        1 => '订单客户',
        2 => '会员客户',
        3 => '导购好友客户',
    ];

    /**
     * 导出导购好友客户列表
     *
     * @return mixed|\Psr\Http\Message\StreamInterface
     */
    public function index()
    {
        $search = request()->get('search', []);
        $search = $search ? json_decode($search, true) : [];
        $search['enterprise_id'] = $this->enterpriseId;
        $search['export'] = true;
        $store = $this->getStores($this->enterpriseAdmin['role'], false) ?? [];

        $data = (new Customer())->friendCustomer(['search' => $search,  'store' => $store]);

        $data = $this->format($data['data']);


        return $this->export($data);
    }

    /**
     * 数据格式化
     *
     * @param $data
     *
     * @return mixed
     */
    public function format($data)
    {
        $userIds = array_column($data, 'user_id');
        $memUid = array_column($data, 'mem_uid');

        // 客户类型
        $userType = (new Customer())->UserType([
            'user_id' => $userIds,
            'enterprise_id' => $this->enterpriseId
        ]);
        $userType = array_column($userType['data'], 'user_type', 'user_id');

        // 客户等级
        $level = (new UserLevel())->userLevel(['enterprise_id' => $this->enterpriseId]);
        $level = array_column($level['data'], null, 'id');

        // 标签
        list($ecUserId, $userId, $label) = (new CustomerService())->getLabel($data, $this->enterpriseId);

        // 所属专属导购
        $member = (new MemberService())->getFirstMemberData(
            $this->qyDomain,
            $this->appIdentifier,
            $memUid
        );

        // 获取离职人员员工客户分配记录
        $memChangeUid = (new \CrsSdk\Work\Customer())->alreadyAllocated(['enterprise_id' => $this->enterpriseId, 'search' => ['data' => $data]])['data'];
        $newMemUid = array_column($memChangeUid, 'ec_userid');

        foreach ($data as &$value) {
            $value['created_mem_uid_at_text'] = $value['created_mem_uid_at'] ? date('Y-m-d', $value['created_mem_uid_at']) : '';
            $value['customer_type'] = $userType[$value['user_id']] ?? [];
            $value['level_name'] = $level[$value['level_id']]['name'] ?? '';
            $value['sales_username'] = '';
            $value['confirmed'] = self::NO_TO_BE_CONFIRMED;
            if (isset($member[$value['mem_uid']])) {
                if (in_array($value['ec_userid'], $newMemUid)) {
                    $value['confirmed'] = self::TO_BE_CONFIRMED;
                }
                $value['sales_username'] = $member[$value['mem_uid']]['memUsername'] . '@' . $member[$value['mem_uid']]['dpName'];
            }

            $value['label'] = [];
            if (in_array($value['ec_userid'], $ecUserId)) {
                $value['label'] = array_merge(
                    $value['label'],
                    array_column($label['data']['ec_userid'][$value['ec_userid']], 'label_name')
                );
            }

            if (in_array($value['user_id'], $userId)) {
                $value['label'] = array_merge(
                    $value['label'],
                    array_column($label['data']['user_id'][$value['user_id']], 'label_name')
                );
            }
            $value['label'] = array_values(array_unique($value['label']));
        }

        return $data;
    }

    /**
     * 数据导出
     *
     * @param $data
     *
     * @return mixed
     */
    public function export($data)
    {
        $cellData[] = ['昵称', '姓名', '客户类型', '标签', '手机号', '专属导购'];

        foreach ($data as $value) {
            if (empty($value['nick_name'])) {
                continue;
            }

            // 客户类型
            $customerType = array_map(function ($value) {
                return self::CUSTOMER_TYPE_TEXT[$value];
            }, $value['customer_type']);

            $cellData[] = [
                $value['nick_name'],
                $value['real_name'],
                implode(',', $customerType),
                implode(',', $value['label']),
                $value['mobile'],
                $value['sales_username']
            ];
        }

        $name = '导购好友客户列表' . date('YmdHi', time());

        return Excel::create($name, function($excel) use ($cellData) {
            $excel->sheet('score', function($sheet) use ($cellData) {
                $sheet->rows($cellData);
            });
        })->export('xlsx');
    }
}
