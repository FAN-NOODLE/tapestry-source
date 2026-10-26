<?php

namespace App\Micros\MicroScrm;

use App\Micros\Common;

class ExternalSearch
{
    public static function choiceOption($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/external-search/choice-option')->post($params);
    }

    public static function list($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/external-search/list')->post($params);
    }

    public static function getOpenIds($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/external-search/get-open-ids')->post($params);
    }

    public static function detail($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/external-search/detail')->post($params);
    }


    public static function addCrowd($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/external-crowd/add')->post($params);
    }

    public function crowdList(array $params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/external-crowd/list')->post($params);
    }

    public function crowdInfo(array $params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/external-crowd/info')->post($params);
    }

    public function crowdOpenIds(array $params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/external-crowd/get-open-ids')->post($params);
    }

    public function crowdEdit(array $params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/external-crowd/edit')->post($params);
    }
    public function crowdDel(array $params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/external-crowd/del')->post($params);
    }
    public function crowdUpdateSort(array $params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/external-crowd/update-sort')->post($params);
    }

    public function crowdUpdateContact(array $params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/external-crowd/update-contact')->post($params);
    }
}

