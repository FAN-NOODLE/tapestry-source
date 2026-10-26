<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2019/1/21
 * Time: 16:39
 */

namespace App\Http\Controllers\V1\Customer;

use App\Http\Controllers\V1\BaseController;
use App\Services\CustomerService;
use App\Services\MemberService;
use CrsSdk\Shop\Customer;
use CrsSdk\Shop\UserLevel;
use Illuminate\Http\Request;
use App\Common\Constant;
use Illuminate\Support\Facades\Log;

class OrderMemberListController extends BaseController
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
     * 下单客户
     *
     * @param Request $request
     *
     * @return array|\Illuminate\Http\JsonResponse
     */
    public function index(Request $request)
    {
        // 分页参数
        $data['page'] = $request->post('page', Constant::DEFAULT_PAGE);
        $data['size'] = $request->post('size', Constant::DEFAULT_SIZE);
        // 查询条件
        $data['search'] = $request->post('search', []);
        Log::debug('下单客户请求参数: ' . json_encode($data));
        $data['enterprise_id'] = $this->enterpriseId;
        $data['store'] = $this->getStores($this->enterpriseAdmin['role'], false) ?? [];

        // 下单客户列表
        $result = (new Customer())->OrderCustomer($data);

        // 下单客户
        $orderCustomer = $result['data'];

        // 获取客户类型
        $userId = array_column($orderCustomer, 'user_id');
        $userType = (new Customer())->UserType([
            'user_id' => $userId,
            'enterprise_id' => $this->enterpriseId,
            'serve_type' => $this->serveType
        ]);
        $userType = array_column($userType['data'], null, 'user_id');

        // 标签
        $customerService = new CustomerService();
        list($ecUserId, $userId, $label) = $customerService
                                            ->getLabel($orderCustomer, $this->enterpriseId, $this->serveType);

        // 等级名称
        $level = (new UserLevel())->userLevel(['enterprise_id' => $this->enterpriseId]);
        $level = array_column($level['data'], null, 'id');

        // 专属导购
        $memUid = array_column($orderCustomer, 'mem_uid');

        // 所属专属导购
        $member = (new MemberService())->getFirstMemberData(
            $this->qyDomain,
            $this->appIdentifier,
            $memUid
        );

        // 获取离职人员员工客户分配记录
        $memChangeUid = (new \CrsSdk\Work\Customer())->alreadyAllocated(['enterprise_id' => $this->enterpriseId, 'search' => ['data' => $data]])['data'];
        $newMemUid = array_column($memChangeUid, 'ec_userid');

        foreach ($result['data'] as &$val) {
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
            $val['ec_id'] = $val['ec_id'] ?? '';
            $val['ec_userid'] = $val['ec_userid'] ?? '';
            $val['level_id'] = $val['level_id'] ?? 0;
            $val['mobile'] = $val['mobile'] ?? '';
            $val['nick_name'] = $val['nick_name'] ?? '';
            $val['real_name'] = $val['real_name'] ?? '';
            $val['union_id'] = $val['union_id'] ?? '';
            $val['confirmed'] = self::NO_TO_BE_CONFIRMED;
            if (isset($member[$val['mem_uid']])) {
                if (in_array($val['ec_userid'], $newMemUid)) {
                    $val['confirmed'] = self::TO_BE_CONFIRMED;
                }
                $val['sales_username'] = $member[$val['mem_uid']]['memUsername'] . '@' . $member[$val['mem_uid']]['dpName'];
            }
        }

        return success($result);
    }
}
