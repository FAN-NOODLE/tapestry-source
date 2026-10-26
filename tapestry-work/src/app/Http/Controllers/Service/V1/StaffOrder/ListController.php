<?php


namespace App\Http\Controllers\Service\V1\StaffOrder;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\StaffOrderRelation;
use App\Services\StaffManagerService;
use App\Services\StaffOrderService;
use Illuminate\Http\Request;

class ListController extends BaseController
{
    /**
     * 内购员工订单列表
     *
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request)
    {
        $input = $request->only(['page', 'size', 'search']);

        $searchData = $input['search'] ?? [];
        $search     = [
            'mem_mobile'     => $searchData['mem_mobile'] ?? '',
            'mem_username'   => $searchData['mem_username'] ?? '',
            'type'           => $searchData['type'] ?? 0,
            'start_time'     => $searchData['start_time'] ?? '',
            'end_time'       => $searchData['end_time'] ?? '',
            'enterprise_id'  => $searchData['enterprise_id'] ?? '',
            'key_word'       => $searchData['key_word'] ?? '',
            'app_identifier' => $searchData['app_identifier'] ?? 'contact',
            'qy_domain'      => $searchData['qy_domain'] ?? '',
        ];


        // 初始化数据
        $page = $input['page'] ?? 1;
        $size = $input['size'] ?? 10;

        list($startTime, $endTime) = (new StaffOrderService())->orderTimeWhere($search);

        $memUids=[];
        // 搜索用户名或者手机号
        if (!empty($search['key_word']) || $search['key_word'] == '0') {

            $paramsMobile = [
                'memMobile' => $search['key_word']
            ];
            // 通过手机号查询
            $memMobileUids = $this->selectUserUids($search, $paramsMobile);
            $paramsName    = [
                'memUsername' => $search['key_word']
            ];
            // 通过用户名号查询
            $memNameUids = $this->selectUserUids($search, $paramsName);

            $memUids = array_merge($memMobileUids, $memNameUids);

            // 因搜索不到用户uid时候查询，需要返回误结果
            $memUids[]='select_select';
        }

        $obj = StaffOrderRelation::query()
            ->when(!empty($startTime), function ($query) use ($startTime) {
                return $query->where('created_at', '>=', $startTime);
            })
            ->when(!empty($endTime), function ($query) use ($endTime) {
                return $query->where('created_at', '<', $endTime);
            })
            ->where('order_pay_time', '>', 0)
            ->when(!empty($memUids), function ($query) use ($memUids) {
                return $query->whereIn('mem_uid', $memUids);
            })
            ->where('order_no', '<>', 1)
            ->where('enterprise_id', $search['enterprise_id'])
            ->orderBy('created_at', 'desc')
            ->paginate($size, ['*'], '', $page);

        $list = $this->format($obj->items(), $search);

        $data = [
            'meta' => [
                'total' => $obj->total(),
                'page'  => (int)$page,
                'size'  => (int)$size
            ],
            'data' => $list
        ];

        return success($data);
    }

    /**
     * 封装查询方法
     * @param array $search
     * @param array $params
     * @return array
     */
    public function selectUserUids($search = [], $params = [])
    {

        $userData = (new StaffManagerService())->getStaffListFromUc($search['qy_domain'], $search['app_identifier'], $params);
        return array_column($userData['list'], 'memUid') ?? [];
    }

    /**
     * 格式化代码
     * @param array $data
     * @param array $search
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
        $list     = [];
        foreach ($data as $value) {
            $user   = $userData[$value['mem_uid']] ?? [];
            $list[] = [
                'mem_mobile'   => $user['memMobile'] ?? '',
                'mem_username' => $user['memUsername'] ?? '',
                'mem_face'     => $user['memFace'] ?? '',
                'dp_names'     => !empty($user['dpName']) ? array_column($user['dpName'], 'dpName') : [],
                'order_no'     => $value['order_no'] ?? '',
                'pay_money'    => $value['pay_money'] ?? 0,
            ];
        }

        return $list;
    }

}