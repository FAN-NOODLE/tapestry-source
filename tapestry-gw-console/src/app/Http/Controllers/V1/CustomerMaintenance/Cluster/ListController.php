<?php

/**
 * 群管理
 * User: 王鑫阳
 * Date: 2020-02-23
 * Time: 10:34
 */

namespace App\Http\Controllers\V1\CustomerMaintenance\Cluster;

use App\Http\Controllers\V1\BaseController;
use App\Services\UcService;
use CrsSdk\Work\Customer;
use Illuminate\Http\Request;

class ListController extends BaseController
{
    public $service = null;

    /**
     * 企业成员
     */
    const STAFF = 1;

    public function __construct(Request $request)
    {
        parent::__construct();
        $this->request = $request;
        $this->service = new UcService($this->qyDomain);
    }

    public function index()
    {
        $params = $this->request->post();

        if (isset($params['tag'])) {
            $tag = $params['tag'];
            unset($params['tag']);
        } else {
            notice(1000010, '参数缺失');
        }

        $data = $this->service->send($tag, json_encode($params));

        $groupData = json_decode($data, true);

        if (empty($groupData['data']['list'])) {
            return success($groupData['data']);
        }

        if ($tag == 'groupMember') {
            // 拿到所有外部联系人ec_id
            $ecId = [];
            foreach ($groupData['data']['list'] as $datum) {
                if ($datum['egmType'] == self::STAFF) {
                    continue;
                }
                $ecId['ec_id'][] = $datum['ecmId'];
            }
            $userIds = $this->isGroupMember($ecId);
            foreach ($groupData['data']['list'] as &$item) {
                $item['union_id'] = '';
                $item['id']       = '';
                if (!empty($userIds)) {
                    foreach ($userIds as $user) {
                        if ($item['ecmId'] == $user['ec_id']) {
                            $item['union_id'] = $user['union_id'];
                            $item['id']       = $user['id'] ?? '';
                        }
                    }
                }
            }
        }

        return success([
            'list'  => $groupData['data']['list'],
            'page'  => $groupData['data']['pageNum'],
            'size'  => $groupData['data']['pageSize'],
            'total' => $groupData['data']['total'],
        ]);
    }

    protected function isGroupMember($ec_ids)
    {
        if (empty($ec_ids['ec_id'])) {
            return [];
        } else {
            $userIds = (new Customer)->getEcIdAndUserId($ec_ids);
            return $userIds['data'];
        }
    }
}
