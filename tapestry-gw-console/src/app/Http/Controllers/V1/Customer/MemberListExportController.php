<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2019/1/25
 * Time: 16:35
 */

namespace App\Http\Controllers\V1\Customer;

use App\Http\Controllers\V1\BaseController;
use App\Services\CustomerService;
use App\Services\MemberService;
use CrsSdk\Shop\Customer;
use CrsSdk\Shop\UserLevel;
use Maatwebsite\Excel\Facades\Excel;
use Illuminate\Http\Request;
use App\Common\Constant;
use Illuminate\Support\Facades\Log;

class MemberListExportController extends BaseController
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
     * 导出会员客户列表
     *
     * @param Request $request
     * @return mixed
     */
    public function index(Request $request)
    {
        // 分页参数
        $data['page'] = $request->get('page', Constant::DEFAULT_PAGE);
        $data['size'] = $request->get('size', Constant::MAX_DEFAULT_SIZE);
        // 查询条件
        $search = $request->get('search', []);
        $search = $search ? json_decode($search, true) : [];
        $data['search'] = $search;
        $data['enterprise_id'] = $this->enterpriseId;
        $data['store'] = $this->getStores($this->enterpriseAdmin['role'], false) ?? [];
        Log::debug('会员导出请求参数: ' . json_encode($data));

        // 下单客户列表
        $result = (new Customer())->UserCustomer($data);
        Log::debug('会员导出返回值: ' . json_encode($result));

        $resultData = $this->format($result['data']);
        Log::debug('会员导出下载: ' . json_encode($resultData));

        return $this->export($resultData);
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
        // 获取客户类型
        $userId = array_column($data, 'user_id');
        $userType = (new Customer())->UserType(['user_id' => $userId, 'enterprise_id' => $this->enterpriseId]);
        $userType = array_column($userType['data'], null, 'user_id');

        // 标签
        $customerService = new CustomerService();
        list($ecUserId, $userId, $label) = $customerService->getLabel($data, $this->enterpriseId);

        // 等级名称
        $level = (new UserLevel())->userLevel(['enterprise_id' => $this->enterpriseId]);
        $level = array_column($level['data'], null, 'id');

        // 专属导购
        $memUid = array_column($data, 'mem_uid');
        // 所属专属导购
        $member = (new MemberService())->getFirstMemberData(
            $this->qyDomain,
            $this->appIdentifier,
            $memUid
        );

        // 获取离职人员员工客户分配记录
        $memChangeUid = (new \CrsSdk\Work\Customer())->alreadyAllocated(['enterprise_id' => $this->enterpriseId, 'search' => ['data' => $data]])['data'];
        $newMemUid = array_column($memChangeUid, 'ec_userid');

        foreach ($data as &$val) {
            $val['label'] = [];
            if (in_array($val['ec_userid'], $ecUserId)) {
                $val['label'] = array_merge(
                    $val['label'],
                    array_column($label['data']['ec_userid'][$val['ec_userid']], 'label_name')
                );
            }
            if (in_array($val['user_id'], $userId)) {
                $val['label'] = array_merge(
                    $val['label'],
                    array_column($label['data']['user_id'][$val['user_id']], 'label_name')
                );
            }
            $val['label'] = array_values(array_unique($val['label']));
            $val['level_name'] = $level[$val['level_id']]['name'] ?? '';
            $val['customer_type'] = $userType[$val['user_id']]['user_type'] ?? '';
            $val['sales_username'] = '';
            $val['avatar'] = $val['avatar'] ?? '';
            $val['mem_uid'] = $val['mem_uid'] ?? '';
            $val['confirmed'] = self::NO_TO_BE_CONFIRMED;
            if (isset($member[$val['mem_uid']])) {
                if (in_array($val['ec_userid'], $newMemUid)) {
                    $val['confirmed'] = self::TO_BE_CONFIRMED;
                }
                $val['sales_username'] = $member[$val['mem_uid']]['memUsername'] . '@' . $member[$val['mem_uid']]['dpName'];
            }
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

        $name = '会员客户列表' . date('YmdHi', time());

        return Excel::create($name, function($excel) use ($cellData) {
            $excel->sheet('score', function($sheet) use ($cellData) {
                $sheet->rows($cellData);
            });
        })->export('xlsx');
    }
}