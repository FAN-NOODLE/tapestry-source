<?php
/**
 * 删除用户印象标签列表
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/1/22
 * Time: 14:39
 */

namespace App\Http\Controllers\App\V1\Label;

use App\Http\Controllers\App\V1\BaseController;
use CrsSdk\Shop\UserLabel;
use Illuminate\Http\Request;

class DeleteController extends BaseController
{
    /**
     * 删除用户印象标签列表
     * @param Request $request
     * @return array
     */
    public function index(Request $request)
    {
        $labelIds = $request->input('label_id');

        // 参数验证规则
        $rules = [
            'label_id' => 'required|array',
        ];
        // 提示语
        $message = [
            'label_id.required' => '印象标签ID不能为空',
            'label_id.array' => '印象标签ID必须是数组',
        ];

        $this->validate($request, $rules, $message);

        // 此处传入当前导购memUid ，只允许删自己打的标签
        $result = (new UserLabel())->delete([
            'label_id' => $labelIds,
            'guide_id' => $this->memUid
        ]);
        return success($result);
    }
}
