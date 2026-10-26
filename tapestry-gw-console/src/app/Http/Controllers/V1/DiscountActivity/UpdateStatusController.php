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
use App\Exceptions\Error;

class UpdateStatusController extends BaseController
{
    /**
     * 活动编辑状态
     * @author 胡龙阳
     * @date 2020/5/28 11:16
     * @param Request $request
     * @return array
     */
    public function index(Request $request)
    {
        // 查询参数处理
        $params = $request->only(['id', 'status']);
        try {
            (new DiscountActivity())->updateStatus($params);
        } catch (\Exception $e) {
            throw new Error($e->getCode(), $e->getMessage(), [], 200);
        }
        return success();
    }
}
