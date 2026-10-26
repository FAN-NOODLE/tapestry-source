<?php
/**
 * 微信优惠卷详情
 * Created by PhpStorm.
 * User: cyj
 * Date: 2019-05-23
 * Time: 17:48
 */

namespace App\Http\Controllers\V1\AppCenter\WxCoupon;

use App\Common\Constant;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\MiniProgram;
use CrsSdk\Shop\WxCoupon;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Redis;

class DetailController extends BaseController
{
    /**
     * @param Request $request
     *
     * @return
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $input = $request->only(['id', 'page', 'size']);
        if (empty($input['id'])) {
            notice(770004);
        }

        $id = intval($input['id']);
        $page = $input['page'] ?? Constant::DEFAULT_PAGE;
        $size = $input['size'] ?? Constant::DEFAULT_SIZE;

        $redisKey = sprintf(Constant::WX_COUPON_DETAIL, $id);
        $redis = Redis::connection();
        $data = $redis->get($redisKey);
        if (empty($data)) {
            // 微信券详情
            $detail = (new WxCoupon())->detail([
                'id'            => $id,
                'enterprise_id' => $this->enterpriseId
            ]);
            $detail = $detail['data'] ?? [];

            $redis->setex($redisKey, 300, json_encode($detail));
        } else {
            $detail = json_decode($data, true);
        }

        $skuList = [
            'list'  => [],
            'total' => 0,
            'page'  => $page,
            'size'  => $size
        ];

        if (!empty($detail) && !empty($detail['tag_goods'])) {
            $skuIds = $detail['tag_goods'];
            if (!is_array($skuIds)) {
                $skuIds = json_decode($skuIds, true);
            }

            $skuList = (new MiniProgram())->getSkuListByIds([
                'enterprise_id' => $this->enterpriseId,
                'sku_ids'       => $skuIds,
                'sort_field'    => 'default',
                'sort'          => 'DESC',
                'page'          => $page,
                'size'          => $size
            ]);

            if (!empty($skuList['list'])) {
                $list = $skuList['list'];
                foreach ($list as $key => $value) {
                    $list[$key]['sku_id'] = $value['id'];
                    $list[$key]['sku_name'] = $value['sku_name'];
                    $list[$key]['sku_attr'] = implode(',', $value['attribute_item_names']);
                    $list[$key]['market_price'] = $value['sign_price'];
                    $list[$key]['sell_price'] = $value['on_line_price'];
                }

                $skuList['list'] = $list;
            }
        }

        return success([
            'data'     => $detail,
            'sku_list' => $skuList
        ]);
    }
}
