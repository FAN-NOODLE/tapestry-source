<?php
/**
 * Created by PhpStorm.
 * User: shizhixiong
 * Date: 2019-01-22
 * Time: 14:19
 */

namespace App\Http\Controllers\Service\V1\Customer;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\MemberUserRelation;
use Illuminate\Support\Facades\Log;

class MemberRelationController extends BaseController
{
    public function index()
    {
        $params = request()->input();
        $condition = [
            ['is_exclusive', MemberUserRelation::IS_EXCLUSIVE_TRUE]
        ];

        if (isset($params['start_time']) && !empty($params['start_time'])) {
            $condition[] = ['updated_at', '>=', $params['start_time']];
        }

        if (isset($params['end_time']) && !empty($params['end_time'])) {
            $condition[] = ['updated_at', '<=', $params['end_time']];
        }
        Log::debug('最近十分钟导购好友修改数据: ' . json_encode($condition));
        $columns = ['id', 'enterprise_id', 'mem_uid', 'ec_name','ec_gender', 'union_id', 'updated_at', 'ec_avatar', 'ec_userid', 'ec_id'];
        $data = MemberUserRelation::query()->where($condition)->get($columns);

        return success(['data' => $data]);
    }
}
