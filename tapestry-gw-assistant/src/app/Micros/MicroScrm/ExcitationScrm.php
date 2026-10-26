<?php

namespace App\Micros\MicroScrm;

use App\Micros\Common;

class ExcitationScrm
{
    /**
     * 首页
     *
     * @param $params
     *
     * @return array|mixed
     */
    public static function index($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/excitation/member-data/data-index')->post($params);
    }

    /**
     * 等级信息
     *
     * @param $params
     *
     * @return array|mixed
     */
    public static function gradePoints($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/excitation/member-data/grade-points')->post($params);
    }

    /**
     * 积分列表
     *
     * @param $params
     *
     * @return array|mixed
     */
    public static function pointsList($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/excitation/points/list')->post($params);
    }

    /**
     * 金币列表
     *
     * @param $params
     *
     * @return array|mixed
     */
    public static function coinsList($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/excitation/coins/mini-coins-list')->post($params);
    }

    /**
     * 金币配置
     *
     * @param $params
     *
     * @return array|mixed
     */
    public static function coinsSettingsInfo($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/excitation/coins/coins-settings-info')->post($params);
    }

    /**
     * 积分规则说明
     * @param $params
     *
     * @return array|mixed
     */
    public static function integralRuleDescription($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/excitation/integral-rule/description-detail')->post($params);
    }

    /**
     * 积分规则信息获取
     * @param $params
     *
     * @return array|mixed
     */
    public static function getPointsRules($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/excitation/member-data/get-points-rules')->post($params);
    }

    /**
     * 等级信息
     *
     * @param $params
     *
     * @return array|mixed
     */
    public static function levelRecord($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/excitation/member-data/level-record')->post($params);
    }

}