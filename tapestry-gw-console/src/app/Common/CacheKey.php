<?php

namespace App\Common;

class CacheKey
{
    /**
     * 发货最后一次使用的物流公司标识
     * %s 企业ID
     */
    const LAST_LOGISTIC_KEY = 'crs:gw-console:last-logistic-key:%s';
}
