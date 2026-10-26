<?php

namespace App\Micros\MicroScrm;

use App\Micros\Common;

class GoMicroScrm
{
    public static function rank($params)
    {

        return Common::MicroScrm()->setApi('micro-scrm/v1/datacompass/rank_data')->post($params);

    }

    public static function importManualPoints($params)
    {

        return Common::MicroScrm()->setApi('micro-scrm/v1/excitation/points/import-manual-points')->post($params);

    }

    public static function importManualCoins($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/excitation/coins/import-manual-coins')->post($params);
    }

    public static function pointsRecordList($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/excitation/points/record-list')->post($params);
    }

    public static function medalRecordList($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/excitation/medal/record-list')->post($params);
    }
    public static function memberList($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/member/list-by-third-ids')->post($params);
    }

    public static function oldMemberList($params){
        return Common::MicroScrm()->setApi('micro-scrm/v1/members/list')->post($params);
    }
    public static function getMemberListByStoreCodes($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/member/get-member-list-by-store-codes')->post($params);
    }


}

