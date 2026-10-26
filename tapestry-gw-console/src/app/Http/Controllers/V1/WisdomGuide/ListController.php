<?php

namespace App\Http\Controllers\V1\WisdomGuide;

use App\Common\Constant;
use App\Micros\Work\Work;
use CrsSdk\Work\StoreWorker;
use CrsSdk\Work\GuideManage;
use CrsSdk\Work\Member;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;

class ListController extends BaseController
{
    /**
     * 导购列表
     *
     * @return JsonResponse
     * @throws \Exception
     */
    public function index(): JsonResponse
    {
        $params                  = $this->request->only(['search', 'page', 'size']);
        $params['enterprise_id'] = $this->enterpriseId;

        // 导购搜索条件
        $search  = [
            'page'          => $params['page'] ?? Constant::DEFAULT_PAGE,
            'size'          => $params['size'] ?? Constant::DEFAULT_SIZE,
            'enterprise_id' => $this->enterpriseId,
            'worker_status' => Constant::IN_SERVICE
        ];
        $memList = $guideMemUid = [];

        // 员工搜索条件
        $condition = [
            'qy_domain'      => $this->qyDomain,
            'app_identifier' => $this->appIdentifier,
            'params'         => [
                'pageNum'  => Constant::DEFAULT_PAGE,
                'pageSize' => 9999,
            ],
            'enterprise_id'  => $this->enterpriseId
        ];

        if (!empty($params['search'])) {
            if (!empty($params['search']['mem_name'])) {

                $condition['params']['memUsername'] = $params['search']['mem_name'];
                $memList                            = (new Member())->memberList($condition);

                $nameFilterMemUids = array_column($memList['list'], 'memUid');

            }
            if (!empty($params['search']['mem_mobile'])) {
                $search['mobile'] = $params['search']['mem_mobile'];
            }
            if (!empty($params['search']['dp_id'])) {
                $condition['params']['dpId'] = $params['search']['dp_id'];
                // 获得当前及子目录的所有数据
                $condition['params']['departmentChildrenFlag'] = 1;
                $memList                                       = (new Member())->memberList($condition);

                if (empty($memList['list'])) {
                    return success([
                        'list'  => [],
                        'total' => 0,
                        'page'  => $search['page'],
                        'size'  => $search['size']
                    ]);
                }
                $dpFilterMemUids = array_column($memList['list'], 'memUid');
            }
            if (isset($dpFilterMemUids) && isset($nameFilterMemUids)) {
                $search['mem_uids'] = array_intersect($dpFilterMemUids, $nameFilterMemUids);
            }
            if (isset($dpFilterMemUids)) {
                $search['mem_uids'] = $dpFilterMemUids;
            }

            if (isset($nameFilterMemUids)) {
                $search['mem_uids'] = $nameFilterMemUids;
            }
            if (isset($search['mem_uids']) && !$search['mem_uids']) {
                return success([
                    'list'  => [],
                    'total' => 0,
                    'page'  => $search['page'],
                    'size'  => $search['size']
                ]);
            }
            // 查询支付即导购列表
            if (!empty($params['search']['filter_type']) && $params['search']['filter_type'] == Constant::GUIDE_TAG) {
                $guideList   = (new GuideManage())->guideBaseLIst([
                    'params' => [
                        'enterprise_id' => $this->enterpriseId,
                        'page'          => Constant::DEFAULT_PAGE,
                        'size'          => Constant::MAX_DEFAULT_SIZE
                    ]
                ]);
                $guideMemUid = array_column($guideList['data'], 'mem_uid');
            }

            if (!empty($params['search']['worker_no'])) {
                $search['worker_no'] = $params['search']['worker_no'];
            }

            if (!empty($params['search']['type'])) {
                $search['type'] = $params['search']['type'];
            }

            if (!empty($params['search']['worker_status'])) {
                $search['worker_status'] = $params['search']['worker_status'];
            }
        }

        $list = (new StoreWorker())->getList($search);

        if (empty($list['list'])) {
            return success([
                'list'  => [],
                'total' => 0,
                'page'  => $search['page'],
                'size'  => $search['size']
            ]);
        }

        if (empty($memList)) {
            $condition['params']['memUids'] = array_column($list['list'], 'mem_uid');
            $memList                        = (new Member())->memberList($condition);
        }
        $memList = array_column($memList['list'], null, 'memUid');

        $return = [];
        foreach ($list['list'] as $item) {
            $dpName   = array_column($memList[$item['mem_uid']]['dpName'] ?? [], 'dpName');
            $tmp      = [
                'id'            => $item['id'],
                'mem_uid'       => $item['mem_uid'],
                'mem_userid'    => $item['mem_userid'],
                'mem_name'      => $memList[$item['mem_uid']]['memUsername'] ?? '',
                'mem_avatar'    => $memList[$item['mem_uid']]['memFace'] ?? '',
                'mem_mobile'    => isset($item['mobile']) && !empty($item['mobile']) ? $item["mobile"] : '',
                'dp_name'       => implode('；', $dpName) ?: '-',
                'type'          => $item['type'],
                'worker_no'     => $item['worker_no'],
                'qr_code'       => $item['qr_code'],
                'created_at'    => (int)$item['created_at'],
                'is_guide'      => in_array($item['mem_uid'], $guideMemUid) ? 1 : 0,
                'store_id'      => $item['store_id'],
                'store_name'    => $item['store_name'] ?? '',
                'worker_status' => isset($memList[$item['mem_uid']]) ? Constant::IN_SERVICE : Constant::OUT_SERVICE,
            ];
            $return[] = $tmp;
        }

        return success([
            'list'  => $return,
            'total' => $list['total'],
            'page'  => $list['page'],
            'size'  => $list['size']
        ]);
    }

    /**
     * member 导购列表
     *
     * @return JsonResponse
     * @Time: 2022/10/19
     */
    public function list(): JsonResponse
    {
        $params = $this->request->only(['search', 'page', 'size']);

        // 导购搜索条件
        $search = [
            'page'          => $params['page'] ?? Constant::DEFAULT_PAGE,
            'size'          => $params['size'] ?? Constant::DEFAULT_SIZE,
            'enterprise_id' => $this->enterpriseId,
        ];

        if (!empty($params['search'])) {

            if (!empty($params['search']['mem_name'])) {
                $search['member_name'] = $params['search']['mem_name'];
            }
            if (!empty($params['search']['mem_mobile'])) {
                $search['mobile'] = $params['search']['mem_mobile'];
            }
            // 导购类型 3 企业导购 4 企业管理
            if (!empty($params['search']['type'])) {
                $search['type'] = $params['search']['type'];
            }
            // 导购在职状态
            if (!empty($params['search']['worker_status'])) {
                $search['worker_status'] = $params['search']['worker_status'];
            }
            // 导购所在的企微部门id
            if (!empty($params['search']['dp_third_id'])) {
                $search['dp_third_id'] = $params['search']['dp_third_id'];
            }
            // 导购的企微userid
            if (!empty($params['search']['worker_no'])) {
                $search['worker_no'] = $params['search']['worker_no'];
            }
        }

        // 查询member
        $list = Work::getMembersList($search);

        return success($list);
    }

}
