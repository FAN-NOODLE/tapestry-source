<?php

namespace App\Micros\MicroScrm;

use App\Micros\Common;
use CrsSdk\Sdk;

class GoMicroScrm
{
    /**
     * 外部联系人列表
     *
     * @param $params
     *
     * @return array|mixed
     * @Time: 2022/10/20
     */
    public static function getExternalUserList($params)
    {
        return Common::microscrm()->setApi('micro-scrm/v1/customer/external-user/list')->post($params);
    }
}

