<?php
/**
 * 获取用户系统标签列表
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/1/24
 * Time: 21:40
 */

namespace App\Http\Controllers\App\V1\Label;

use App\Http\Controllers\App\V1\BaseController;
use App\Models\MemberUserRelation;
use CrsSdk\Shop\UserLabel;
use Illuminate\Http\Request;

class SystemListController extends BaseController
{
    /**
     * 获取用户系统标签列表
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

        $search['enterprise_id'] = $this->enterpriseId;

        $list = (new UserLabel())->listSystem([
            'page' => $page,
            'size' => $size,
            'search' => $search
        ]);

        return  success($list);
    }
}
