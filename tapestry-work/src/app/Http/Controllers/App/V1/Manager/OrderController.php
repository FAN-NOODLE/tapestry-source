<?php
/**
 * 订单列表
 * Created by PhpStorm.
 * User: 刘瑞
 * Date: 2018/11/26
 * Time: 10:04
 */

namespace App\Http\Controllers\App\V1\Manager;

use App\Common\Constant;
use App\Http\Controllers\App\V1\BaseController;
use App\Services\OrderService;
use CrsSdk\Shop\Order;
use \Illuminate\Http\Request;
use App\Models\MemberUserRelation;
use App\Library\Uc\Member;
use App\Exceptions\Error;

class OrderController extends BaseController
{
    /**
     * 是否有用户信息 1：是，2：否 2
     */
    const ORDER_USER_TRUE = 1;

    /**
     * 接口入口方法
     *
     * @param Request $request
     * @param OrderService $orderService
     * @return array|mixed|\Psr\Http\Message\StreamInterface
     * @throws \Illuminate\Validation\ValidationException
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/12/21 10:24 AM
     */
    public function index(Request $request, OrderService $orderService)
    {
        $input = $request->only(['search', 'page', 'size']);
        $input['search']['enterprise_id'] = $this->enterpriseId;
        $input['size'] = $input['size'] ?? Constant::DEFAULT_SIZE;

        // 接值
        $rules = [
            'search.status' => 'integer',
            'search.order_type' => 'integer',
            'page' => 'integer',
            'size' => 'integer',
        ];
        $messages = [
            'search.status.integer' => '订单状态必须为正整数',
            'search.order_type.integer' => '订单类型必须为正整数',
            'page.integer' => '页码必须为正整数',
            'size.integer' => '每页数量必须为正整数',
        ];
        $this->validate($request, $rules, $messages);

        // 导购ID
        $input['search']['mem_uid'] = $this->byListMemUid($orderService);
        if (empty($input['search']['mem_uid'])) {
            return success([
                'data' => [],
                'meta' => [
                    'page' => $input['page'] ?? Constant::DEFAULT_PAGE,
                    'size' => $input['size'],
                    'total' => 0
                ]
            ]);
        }
        // 订单列表
        $list = (new Order())->orderList($input);

        // 格式化
        $list['data'] = $this->formalData($list['data'], $input['size']);

        return $list;
    }

    /**
     * 获取导购UID
     *
     * @param $orderService
     * @return array|mixed
     * @throws \Exception
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/12/21 10:25 AM
     */
    private function byListMemUid($orderService)
    {
        // 获取当前登录用户的详情信息
        $ucResult = (new Member())->detail($this->qyDomain, $this->appIdentifier, ['memUid' => $this->memUid]);
        if (empty($ucResult)) {
            return error(600103);
        }
        $dpId = '';
        if (isset($ucResult['dpName'])) {
            // 取出所属部门id
            $dpId = $ucResult['dpName'][0]['dpId'];
        }
        if (empty($dpId)) {
            // 部门信息不存在
            return error(600102);
        }

        // 获取人员列表
        $memList = $orderService->listAllMember(
            ['qyDomain' => $this->qyDomain, 'appIdentifier' => $this->appIdentifier, 'enterpriseId' => $this->enterpriseId],
            ['dpId' => $dpId, 'memAll' => Constant::MEM_ALL_UC]
        );

        return empty($memList) ? [] : array_keys($memList);
    }

    /**
     * 格式化数据
     *
     * @param array $data
     * @param int $size
     * @return array
     * @throws \Exception
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/12/21 10:25 AM
     */
    private function formalData($data = [], $size = Constant::DEFAULT_SIZE)
    {
        $resultData = [];
        // 会员信息
        $userData = $this->byList($data);

        // 导购信息
        $memberData = (new Member())->ignorePermissionsList(
            $this->qyDomain,
            $this->appIdentifier,
            ['memUids' => array_values(array_filter(array_column($data, 'mem_uid'))),
                'pageNum' => Constant::DEFAULT_PAGE,
                'pageSize' => $size,
                'memAll' => Constant::MEM_ALL_UC
            ]);
        $memberData = empty($memberData['list']) ? [] : array_column($memberData['list'], NULL, 'memUid');

        // 导购信息
        foreach ($data as $key => $val) {
            // 昵称
            $name = $val['nick_name'];
            if (!empty($val['real_name'])) {
                $name .= ' (' . $val['real_name'] . ')';
            } else {
                // 备注
                $name .= empty($userData[$val['union_id']]) ? '' : ' (' . $userData[$val['union_id']] . ')';
            }

            $resultData[$key] = [
                'id' => $val['id'],
                'order_no' => $val['order_no'],
                'status' => $val['status'],
                'final_amount' => $val['final_amount'],
                'pay_at' => (int)$val['pay_at'],
                'created_at' => (int)$val['created_at'],
                'quantity' => $val['quantity'],
                'user_name' => $name,
                'mobile' => $val['mobile'],
                'mem_name' => $memberData[$val['mem_uid']]['memUsername'] ?? '',
                'mem_status' => $memberData[$val['mem_uid']]['memStatus'] ?? Constant::MEM_DEALETE_UC,
                'goods_data' => $val['goods_data'],
            ];
            if (!empty($val['group_id'])) {
                $resultData[$key]['type'] = $val['type'];
                $resultData[$key]['group_id'] = $val['group_id'];
            }

        }

        return $resultData;
    }

    /**
     * 用户备注
     * @param array $list
     * @return array
     */
    public function byList($list = [])
    {
        // 取出没有名称的用户unionId
        $unionIds = [];
        foreach ($list as $key => $val) {
            if (empty($val['real_name'])) {
                $unionIds[] = $val['union_id'];
            }
        }

        if (!empty($unionIds)) {
            // 查询备注
            $unionIdsData = MemberUserRelation::select('union_id', 'user_remark_name')
                ->whereIn('union_id', $unionIds)
                ->get()->toArray();
            return array_column($unionIdsData, 'user_remark_name', 'union_id');
        }

        return [];
    }

}
