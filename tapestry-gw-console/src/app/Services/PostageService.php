<?php
/**
 * Created by PhpStorm.
 * User: pw
 * Date: 2020/4/2
 * Time: 16:52
 */

namespace App\Services;

use CrsSdk\Shop\PostageRelation;

class PostageService extends BaseService
{

    /**
     * 售卖类型 1：快递配送
     */
    const ON_LINE = 1;

    /**
     * 售卖类型 2：线下自提
     */
    const OFF_LINE = 2;

    /**
     * 售卖类型 3：同城配送
     */
    const SELF_MENTION = 3;


    /**
     * sku售卖方式 1：仅线上售卖
     */
    const ON_LINE_SALE = 1;

    /**
     * sku售卖方式 2：仅线下自提
     */
    const OFF_LINE_SALE = 2;

    /**
     * sku售卖方式 3：线上线下均有售
     */
    const ALL_SALE = 3;

    /**
     * 根据sku_id查询线上线下售卖方式
     *
     * @param $skuId
     * @return int
     */
    public function getSkuSaleType($skuId)
    {

        $relation = (new PostageRelation())->relationList(['object_type' => 1, 'object_id' => $skuId]);

        // sku售卖方式
        $sales = 0;
        if (!empty($relation['list'])) {

            $types = array_column($relation['list'], 'delivery_type', null);

            if (in_array(self::ON_LINE_SALE, $types) || in_array(self::SELF_MENTION, $types)) $sales = self::ON_LINE_SALE;
            if (in_array(self::OFF_LINE_SALE, $types)) $sales = self::OFF_LINE_SALE;
            if ((in_array(self::ON_LINE_SALE, $types) || (in_array(self::SELF_MENTION, $types))) && in_array(self::OFF_LINE_SALE, $types)) $sales = self::ALL_SALE;
        }

        return $sales;
    }

    /**
     * 根据邮费类型查询线上线下售卖方式
     *
     * @param array $deliveryType
     *
     * @return int
     */
    public function getSaleTypeByDeliver(array $deliveryType)
    {
        // sku售卖方式
        $sales = self::ALL_SALE;

        if (empty($deliveryType)) return $sales;

        // 选择了快递或同城
        if (in_array(self::ON_LINE_SALE, $deliveryType) || in_array(self::SELF_MENTION, $deliveryType)) $sales = self::ON_LINE_SALE;
        // 选择了自提
        if (in_array(self::OFF_LINE_SALE, $deliveryType)) $sales = self::OFF_LINE_SALE;
        // 选择了快递或同城，并选择了自提
        if ((in_array(self::ON_LINE_SALE, $deliveryType) || (in_array(self::SELF_MENTION, $deliveryType))) && in_array(self::OFF_LINE_SALE, $deliveryType)) $sales = self::ALL_SALE;

        return $sales;
    }
}
