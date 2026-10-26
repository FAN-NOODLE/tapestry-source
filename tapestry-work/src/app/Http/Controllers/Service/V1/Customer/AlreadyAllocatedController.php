<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2019/3/21
 * Time: 16:38
 */

namespace App\Http\Controllers\Service\V1\Customer;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\MemberChangeRecord;
use App\Models\MemberUserRelation;
use Illuminate\Http\Request;
use App\Common\Constant;

class AlreadyAllocatedController extends BaseController
{
    /**
     * 获取
     *
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request)
    {
        // 分页参数
        $page = $request->post('page', Constant::DEFAULT_PAGE);
        $size = $request->post('size', 0);
        $enterpriseId = $request->post('enterprise_id', 0);
        // 查询条件
        $data = $request->get('search', [])['data'];
        // 获取当前客户的关系
        $ecUserid = array_unique(array_filter(array_column($data, 'ec_userid')));
        $memberUserRelation = MemberUserRelation::query()->whereIn('ec_userid', $ecUserid)->where('enterprise_id', $enterpriseId)->get()->toArray();
        // 归类同意分类的属性值
        $result = [];
        foreach ($memberUserRelation as $attribute) {
            $result[$attribute['ec_userid']][] = $attribute;
        }
        // 只有一个专属导购
        foreach ($result as $key => $val) {
            if (count($val) > 1) {
                unset($result[$key]);
                continue;
            }
            $result[$key] = $val[0];
        }

        // 查询单一导购的导购好友分配
        $orm = MemberChangeRecord::query()
            ->where('status', MemberChangeRecord::UNALLOCATED)
            ->where('enterprise_id', $enterpriseId);

        $orm->where(function ($query) use ($result) {
            foreach ($result as $val) {
                $query->orwhere(function ($query) use ($val) {
                    $query->where('ec_userid', $val['ec_userid'])
                        ->Where('mem_uid', $val['mem_uid']);
                });
            }
        });

        $data = $orm->get(['ec_userid', 'old_mem_uid', 'mem_uid', 'status']);

        return success([
            'data' => $data,
            'meta' => []
        ]);
    }
}
