<?php
/**
 * 添加用户印象标签
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/1/22
 * Time: 14:20
 */

namespace App\Http\Controllers\App\V1\Label;

use App\Http\Controllers\App\V1\BaseController;
use App\Models\MemberUserRelation;
use CrsSdk\Shop\UserLabel;
use Illuminate\Http\Request;

class CreateController extends BaseController
{

    /**
     * 添加用户印象标签
     * @param Request $request
     * @return array
     */
    public function index(Request $request)
    {
        $userId = $request->input('ec_userid');
        $labelIds = $request->input('label_id');

        // 参数验证规则
        $rules = [
            'label_id' => 'required|array',
            'ec_userid' => 'required',
        ];
        // 提示语
        $message = [
            'label_id.required' => '印象标签ID不能为空',
            'label_id.array' => '印象标签ID必须是数组',
            'ec_userid.required' => '客户ID不能为空',
        ];

        $this->validate($request, $rules, $message);
        // 获取union_id  企业管理后台数据分析用
        $userRelation = MemberUserRelation::query()->where(['ec_userid' => $userId])->first()->toArray();

        $result = (new UserLabel())->add([
            'ec_userid' => $userId,
            'label_id' => $labelIds,
            'guide_id' => $this->memUid,
            'enterprise_id' => $this->enterpriseId,
            'union_id' => $userRelation['union_id']
        ]);

        return success($result);
    }
}
