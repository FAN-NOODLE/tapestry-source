<?php

namespace App\Micros\Work;

use App\Micros\Common;

class Work
{
    // 部门列表
    public static function lists($params)
    {
        return Common::Work()->setApi('service/v1/department/lists')->post($params);
    }

    // 导购列表
    public static function getMembersList($params)
    {
        return Common::Work()->setApi('service/v1/member/get-members-list')->post($params);
    }

    // 生成太阳码
    public static function createQrcode($params)
    {
        return Common::Work()->setApi('service/v1/common/create-qrcode')->post($params);
    }

}

