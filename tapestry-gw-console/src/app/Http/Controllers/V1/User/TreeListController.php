<?php

namespace App\Http\Controllers\V1\User;

use App\Micros\MicroScrm\GoMicroScrm;
use App\Common\Constant;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;

class TreeListController extends BaseController
{
    /**
     * @return JsonResponse
     * @throws \Exception
     */
    public function index()
    {
        // 分页参数
        $page = $this->request->input('page', Constant::DEFAULT_PAGE);
        $size = $this->request->input('size', 50);
        // 查询条件
        $search = $this->request->input('search', []);

        if (empty($this->qyDomain)) {
            return success([
                'list'  => [],
                'total' => 0,
                'page'  => intval($page),
                'size'  => intval($size),
            ]);
        }

        if (!empty($search['mem_username']) || !empty($search['userid']) || !empty($search['mobile'])){
            $condition = [
                'enterprise_id'  => $this->enterpriseId,
                'page'=>$page,
                'size' => $size,
            ];
            $list['list']  = $list['member_list'] ?? [];
            if (!empty($search['mem_username'])) {
                $condition['name'] = $search['mem_username'];
            }
            if (!empty($search['userid'])) {
                $condition['userid'] = $search['userid'];
            }
            if (!empty($search['mobile'])) {
                $condition['mobile'] = $search['mobile'];
            }
            $list = GoMicroScrm::oldMemberList($condition);
            $list['list']  = $list['member_list'] ?? [];
            // 总数
            $total = $list['total'] ?? 0;
        } else{
            // 调用sdk
            $thirdIds = isset($search['dp_id']) ? [$search['dp_id']] : [];
            $condition = [
                'third_ids' => $thirdIds,
                'enterprise_id' => $this->enterpriseId
            ];
            $list = GoMicroScrm::memberList($condition);
            // 总数
            $total = count($list);
            $list['list']  = $list;
            $page = 1;
            $size = 9999;
        }
        // 用户列表数据格式化
        $data = [];
        foreach ($list['list'] as $member) {
            $data[] = [
                'user_id' => $member['userid'],
                'user_name' => $member['name'],
                'dp_name' => $member['department_name']??$member['forms'],
                'mobile' => $member['mobile'],
                'avatar' => $member['avatar'],
                'is_guide' => 1,
                'mem_subscribe_status' => 1
            ];
        }

        // 有分页
        return success([
            'list'  => $data,
            'total' => $total,
            'page'  => intval($page),
            'size'  => intval($size)
        ]);
    }
}