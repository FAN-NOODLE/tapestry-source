<?php

namespace App\Http\Controllers\V1\User;

use App\Common\Constant;
use CrsSdk\Work\Member;
use Illuminate\Http\Request;
use App\Http\Controllers\V1\BaseController;

class GuideVUserListController extends BaseController
{
    /**
     * 导购版员工列表
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     */
    public function index(Request $request)
    {
        // 分页参数
        $page = $request->post('page', Constant::DEFAULT_PAGE);
        $size = $request->post('size', 10);
        // 查询条件
        $search = $request->post('search', []);

        if (empty($this->qyDomain)) {
            return success([
                'list'  => [],
                'page'  => intval($page),
                'size'  => intval($size),
                'total' => 0,
            ]);
        }

        // 参数列表
        $params = [
            'pageNum'  => $page,
            'pageSize' => $size
        ];

        // 按删除时间倒序
        //   $params['orderList'][] = ['column' => 'memDeleted', 'orderType' => 'DESC'];
        // 在职员工，按创建时间倒序
        $params['orderList'][] = ['column' => 'memCreated', 'orderType' => 'DESC'];

        $params['memAll'] = 1;
        // 姓名
        if (isset($search['mem_username'])) {
            $params['memUsername'] = $search['mem_username'];
        }

        // 手机号码
        if (isset($search['mem_mobile'])) {
            $params['memMobile'] = $search['mem_mobile'];
        }

        // 门店id
        if (isset($search['dp_id']) && $search['dp_id']) {
            $params['dpId'] = $search['dp_id'];
        }

        // 查询是否在职
        if (!empty($search['mem_subscribe_status'])) {
            // 查询在职人员需要去掉
            if ($search['mem_subscribe_status'] == 1) {
                unset($params['memAll']);
            }
            $params['memSubscribeStatus'] = $search['mem_subscribe_status'];
        }

        // 查询状态
        if (!empty($search['eq_mem_status'])) {
            $params['eqMemStatus'] = $search['eq_mem_status'];
        }

        // 调用sdk
        $condition = [
            'qy_domain'      => $this->qyDomain,
            'app_identifier' => $this->appIdentifier,
            'params'         => $params,
            'enterprise_id'  => $this->enterpriseId
        ];

        // 没有部门查询没有部门接口
        if(!empty($search['no_dp'])){
            $list = (new Member())->memberNoDpList($condition);
        }else{
            $list = (new Member())->memberList($condition);
        }

        // 总数
        $total = intval($list['total']);

        // 用户列表数据格式化
        $data = $this->format($list['list']);

        $userData = [
            'list'  => $data,
            'page'  => intval($page),
            'size'  => intval($size),
            'total' => $total,
        ];

        return success($userData);
    }

    /**
     * 格式化员工列表
     * @param $list
     * @return array
     */
    protected function format($list)
    {
        // 初始化待返数据
        $data = [];
        if (!$list) {
            return $data;
        }

        // 获取离职员工mem_uid集合
        foreach ($list as $member) {
            // 员工id
            $memUid = $member['memUid'];
            $data[] = $this->formatMemberTrue($memUid, $member);
        }

        return $data;
    }

    /**
     * 格式化 在职员工数据
     * @param $memUid
     * @param $member
     * @return array
     */
    public function formatMemberTrue($memUid, $member)
    {
        //dd($member);
        return [
            'mem_uid'              => $memUid,
            'mem_user_id'          => $member['memUserid'],
            'mem_username'         => $member['memUsername'],
            'mem_face'             => $member['memFace'],
            'dp_name'              => $member['dpName'] ?? [],
            'mem_role'             => $member['memRole'] ?? '',
            'mem_gender'           => $member['memGender'],
            'mem_mobile'           => $member['memMobile'],
            'mem_active'           => $member['memActive'],
            'mem_subscribe_status' => $member['memSubscribeStatus']
        ];
    }
}