<?php

namespace App\Http\Controllers\Service\V1\Customer;

use App\Common\Constant;
use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\Member;

class StaffListController extends BaseController
{
    public function staffList()
    {
        $params  = $this->request->all();
        $name    = $params['name'] ?? null;
        $dep_uid = $params['dep_uid'] ?? null;

        $size = $params['size'] ?? Constant::DEFAULT_SIZE;
        $page = $params['page'] ?? Constant::DEFAULT_PAGE;

        $contains = [
            'pageNum'  => $page,
            'pageSize' => $size,
        ];
        if (!empty($dep_uid)) {
            $contains = array_merge($contains, ['dpId' => $dep_uid]);
        }

        if (!empty($name)) {
            $contains['memUsername'] = $name;
        }

        $result = (new Member())->ignorePermissionsList($params['qy_domain'], $params['app_identifier'], $contains);

        $data = [];
        if (isset($result['list']) && count($result['list']) > 0) {
            foreach ($result['list'] as $item) {
                $memIds[] = $item['memUid'];
                $data[]   = [
                    'name'    => $item['memUsername'],
                    'men_uid' => $item['memUid'],
                ];
            }
        }

        return success([
            'data' => $data,
            'meta' => [
                'page'  => (int)$page,
                'size'  => (int)$size,
                'total' => (int)$result['total'] ?? 0,
            ]
        ]);
    }
}