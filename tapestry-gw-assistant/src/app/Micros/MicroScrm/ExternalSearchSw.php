<?php

namespace App\Micros\MicroScrm;

use App\Micros\Common;

class ExternalSearchSw
{
    public static function choiceOption($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/external-search-sw/choice-option')->post($params);
    }

    public static function list($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/external-search-sw/list')->post($params);
    }

    public static function getOpenIds($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/external-search-sw/get-open-ids')->post($params);
    }


    public static function addCrowd($params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/external-crowd-sw/add')->post($params);
    }

    public function crowdList(array $params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/external-crowd-sw/list')->post($params);
    }

    public function crowdInfo(array $params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/external-crowd-sw/info')->post($params);
    }

    public function crowdOpenIds(array $params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/external-crowd-sw/get-open-ids')->post($params);
    }

    public function crowdEdit(array $params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/external-crowd-sw/edit')->post($params);
    }
    public function crowdDel(array $params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/external-crowd-sw/del')->post($params);
    }
    public function crowdUpdateSort(array $params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/external-crowd-sw/update-sort')->post($params);
    }

    public function crowdUpdateContact(array $params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/external-crowd-sw/update-contact')->post($params);
    }
    public function crowdBatchUpdateContact(array $params)
    {
        return Common::MicroScrm()->setApi('micro-scrm/v1/customer/external-crowd-sw/batch-update-contact')->post($params);
    }
}

