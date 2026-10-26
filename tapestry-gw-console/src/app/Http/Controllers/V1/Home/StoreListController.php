<?php
/**
 * Created by PhpStorm.
 * User: liupengbo
 * Date: 2020/4/21
 * Time: 16:24
 */

namespace App\Http\Controllers\V1\Home;

use App\Common\Constant;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\Store;

/**
 * @Notes:获取当前登录帐号权限下的门店列表
 * Class StoreListController
 * @package App\Http\Controllers\V1\Home
 * @author: liupengbo
 * @Time: 2020/4/22   15:24
 */
class StoreListController extends BaseController
{
    /**
     * 有商城+门店的权限
     *
     * @var int
     */
    const AUTH_SHOP_AND_ENTERPRISE = 1;

    /**
     * 只有商城权限
     *
     * @var int
     */
    const AUTH_SHOP = 2;

    /**
     * 只有门店权限
     *
     * @var int
     */
    const AUTH_STORE = 3;

    /**
     * 分批调用门店信息
     *
     * @var int
     */
    const ARRAY_CHUNK_LIMIT = 100;

    /**
     * 是商城
     *
     * @var int
     */
    const IS_SHOP = 1;

    /**
     * 是门店
     *
     * @var int
     */
    const IS_STORE = 2;

    /**
     * @return \Illuminate\Http\JsonResponse
     */
    public function index()
    {
        $data = [
            'type' => 0,
            'list' => []
        ];

        $defaultEnterprise = [
            'store_id'   => $this->enterprise['id'],
            'name'       => $this->enterprise['name'],
            'short_name' => '',
            'type'       => self::IS_SHOP
        ];
        // 判断当前登录管理员的权限
        if (Constant::HAS_SHOP_AUTH == $this->hasShopAuth && !empty($this->storeIds)) {
            $data['type']   = self::AUTH_SHOP_AND_ENTERPRISE;
            $data['list'][] = $defaultEnterprise;
        } else if (!empty($this->storeIds)) {
            $data['type'] = self::AUTH_STORE;
        } else {
            $data['type'] = self::AUTH_SHOP;
            $data['list'][] = $defaultEnterprise;
        }

        // 批量查询门店信息
        $storeData = [];
        if (!empty($this->storeIds)) {
            $storeIdChunk = array_chunk($this->storeIds, self::ARRAY_CHUNK_LIMIT);
            foreach ($storeIdChunk as $storeIds) {
                $stores    = app(Store::class)->lists([
                    'store_ids'     => $storeIds,
                    'enterprise_id' => $this->enterpriseId,
                    'page'          => Constant::DEFAULT_PAGE,
                    'size'          => count($storeIds)
                ]);
                $storeData = array_merge($storeData, $stores['list']);
            }

            // 格式化返回字段
            array_walk($storeData, function ($store) use (&$data) {
                array_push($data['list'], [
                    'name'       => $store['name'],
                    'short_name' => $store['short_name'],
                    'store_id'   => $store['id'],
                    'type'       => self::IS_STORE
                ]);
            });
        }
        return success($data);
    }
}
