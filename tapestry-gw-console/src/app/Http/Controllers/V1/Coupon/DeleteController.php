<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/26
 * Time: 10:04
 */

namespace App\Http\Controllers\V1\Coupon;

use Illuminate\Http\Request;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\Coupon;

class DeleteController extends BaseController
{
    /**
     * 礼劵删除
     *
     * @param Request $request
     *
     * @return array|\Illuminate\Http\JsonResponse
     *@author 梁欢 lianghuan@vchangyi.com
     * @date 2018/11/30 11:16
     */
    public function index(Request $request)
    {
        $params = $request->only(['coupon_id']);

        // 礼劵ID不能为空
        if (!isset($params['coupon_id']) || empty($params['coupon_id'])) {
            error(400001);
        }

        // 查询礼券详情
        $data = (new Coupon())->detail(['coupon_id' => $params['coupon_id']]);

        $this->writeOperatorLog('删除礼券:' . $data['name']);

        $result = (new Coupon())->delete(['coupon_id' => [$params['coupon_id']]]);

        if (empty($result['error_ids'])) {
            return success([]);
        }

        // 礼券被占位符占用，无法删除
        if (!empty($result['placeholder_error_ids'])) {
            error(400003, '占位符');
        }

        // 礼券被收藏有礼占用，无法删除
        if (!empty($result['collection_error_ids'])) {
            error(400003, '收藏有礼');
        }
    }
}
