<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/26
 * Time: 10:04
 */

namespace App\Http\Controllers\V1\DiscountActivity;

use Illuminate\Http\Request;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\DiscountActivity;

class DetailController extends BaseController
{
    /**
     * 活动详情
     * @author 胡龙阳
     * @date 2020/5/28 11:16
     * @param Request $request
     * @return array
     */
    public function index(Request $request)
    {
        // 查询参数处理
        $params = $request->only(['id']);

        $rules    = [
            'id' => 'required',
        ];
        $messages = [
            'id.required' => '活动ID不能为空',
        ];
        $this->validate($this->request, $rules, $messages);

        // 查询
        $data = (new DiscountActivity())->detail($params);

        return $data;
    }
}
