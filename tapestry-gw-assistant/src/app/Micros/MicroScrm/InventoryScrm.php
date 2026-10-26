<?php

namespace App\Micros\MicroScrm;

use App\Micros\Common;

class InventoryScrm
{
    /**
     * 库存查询
     *
     * @param $params
     *
     * @return array|mixed
     */
    public static function list($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/inventory/list')->post($params);
    }


    /**
     * 库存查询的历史搜索记录
     *
     * @param $params
     *
     * @return array|mixed
     */
    public static function searchKeywords($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/inventory/search-keywords')->post($params);
    }


    /**
     * 库存查询的历史搜索记录
     *
     * @param $params
     *
     * @return array|mixed
     */
    public static function delKeywords($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/inventory/del-keywords')->post($params);
    }


}