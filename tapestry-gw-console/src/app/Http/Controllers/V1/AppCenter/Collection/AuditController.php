<?php
/**
 * Created by PhpStorm.
 * User: zhuhongda
 * Date: 2019/5/14
 * Time: 10:04
 */

namespace App\Http\Controllers\V1\AppCenter\Collection;

use CrsSdk\Shop\Collection;
use Illuminate\Http\Request;
use App\Http\Controllers\V1\BaseController;
use Illuminate\Validation\Rule;

class AuditController extends BaseController
{
    /**
     * 用户提交收藏有礼审核
     * @author 朱虹达 zhuhongda@vchangyi.com
     * @date 2019/05/14 11:16
     * @param Request $request
     * @return array
     */
    public function index(Request $request)
    {
        // 接受参数
        $input = $request->only([
            'state',
            'id'
        ]);

        // 参数验证规则
        $rules = [
            'state' => ['required', Rule::in([1, 2, 3])],
            'id'    => 'required'
        ];

        // 提示语
        $message = [
            'state.required' => '审核状态不能为空',
            'id.required'    => '列表id不能为空'
        ];

        $this->validate($request, $rules, $message);

        $input['enterprise_id'] = $this->enterpriseId;

        return (new Collection())->audit($input);
    }

}
