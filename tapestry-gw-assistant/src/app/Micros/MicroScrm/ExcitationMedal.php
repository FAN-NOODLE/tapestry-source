<?php

namespace App\Micros\MicroScrm;

use App\Micros\Common;

class ExcitationMedal
{
    /**
     * 奖章墙
     *
     * @param $params
     *
     * @return array|mixed
     */
    public static function MedalWall($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/excitation/medal/medal-wall')->post($params);
    }

    /**
     * 奖章墙详情
     *
     * @param $params
     *
     * @return array|mixed
     */
    public static function MedalWallDetail($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/excitation/medal/medal-wall-detail')->post($params);
    }

    /**
     * 领取奖章数量和弹框
     *
     * @param $params
     *
     * @return array|mixed
     */
    public static function RecordAlreadyGot($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/excitation/medal/record-already-got')->post($params);
    }

    /**
     * 领取奖章
     * @param $params
     *
     * @return array|mixed
     */
    public static function RecordMemberReceived($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/excitation/medal/record-member-received')->post($params);
    }
}