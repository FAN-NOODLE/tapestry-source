<?php
/**
 * 获取用户全部标签列表
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/1/24
 * Time: 21:27
 */

namespace App\Http\Controllers\App\V1\Label;

use App\Http\Controllers\App\V1\BaseController;
use App\Models\MemberUserRelation;
use CrsSdk\Shop\UserLabel;
use Illuminate\Http\Request;

class ListAllController extends BaseController
{
    /**
     * 获取用户全部标签列表
     * @param Request $request
     * @return array
     */
    public function index(Request $request)
    {
        $page = $request->post('page', 1);
        $size = $request->post('size', 10);
        $search = $request->post('search', []);

        if (!isset($search['ec_userid']) || empty($search['ec_userid'])) {
            return success([
                'data' => [],
                'meta' => [
                    'total' => 0,
                    'page' => $page,
                    'size' => $size
                ]
            ]);
        }

        // 拼接查询参数
        $userRelation = MemberUserRelation::query()->where(['ec_userid' => $search['ec_userid']])->first()->toArray();

        if (isset($userRelation['union_id'])) {
            $search['union_id'] = $userRelation['union_id'];
        }
        // 拼接企业、员工信息
        $search['guide_id'] = $this->memUid;
        $search['enterprise_id'] = $this->enterpriseId;
        $search['store_status'] = $this->storeStatus;

        $list = (new UserLabel())->listAll([
            'page' => $page,
            'size' => $size,
            'search' => $search
        ]);

        return success($list);
    }
}
