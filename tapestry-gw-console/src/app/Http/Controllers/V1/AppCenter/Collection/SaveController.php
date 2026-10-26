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

class SaveController extends BaseController
{
    /**
     * 保存收藏有礼
     * @param Request $request
     *
     * @return mixed|\Psr\Http\Message\StreamInterface
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index(Request $request)
    {
        // 接受参数
        $input = $request->only([
            'is_popup',
            'coupon_id'
        ]);

        // 参数验证规则
        $rules = [
            'is_popup'    => ['required', Rule::in([1, 2])],
            'coupon_id' => 'required'
        ];

        // 提示语
        $message = [
            'is_popup.required'    => '是否弹窗不符合规范',
            'coupon_id.required' => '优惠券id不能为空'
        ];

        $this->validate($request, $rules, $message);

        $input['enterprise_id'] = $this->enterpriseId;

        // 保存收藏有礼
        return (new Collection())->save($input);

    }

}
