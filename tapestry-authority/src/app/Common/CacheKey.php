<?php

namespace App\Common;

class CacheKey
{
    /**
     * 角色
     * %s 企业ID
     * %s 角色ID
     */
    const ROLE = 'crs:micro-authority:role:%s:%s';
}
