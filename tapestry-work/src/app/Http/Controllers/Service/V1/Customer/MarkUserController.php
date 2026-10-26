<?php

namespace App\Http\Controllers\Service\V1\Customer;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Models\GuideUserRelation;
use Illuminate\Http\Request;

class MarkUserController extends BaseController
{
    /**
     * 标记是否复制手机号
     *
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request)
    {
        $userId = $request->input('user_id', 0);
        $memUid = $request->input('mem_uid', 0);

        if (empty($userId) && empty($memUid)) {
            throw new Error(10000, '参数user_id和参数mem_uid不能为空');
        }

        // 是否已复制手机号（默认1）:0、否, 1、是
        $isCopyMobile = $request->input('is_copy_mobile', 1);

        if (!in_array($isCopyMobile, [0, 1])) {
            throw new Error(20000, '参数is_copy_mobile不在范围');
        }

        GuideUserRelation::query()
            ->when($memUid, function ($query) use ($memUid) {
                $query->where('mem_uid', $memUid);
            })
            ->when($userId, function ($query) use ($userId) {
                $query->where('user_id', $userId);
            })
            ->update([
                'is_copy_mobile' => $isCopyMobile
            ]);

        return success(['data' => []]);
    }
}
