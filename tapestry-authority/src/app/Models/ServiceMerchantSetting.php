<?php

namespace App\Models;

class ServiceMerchantSetting extends Base
{
    const ENABLE = 1;

    const DISABLE = 2;

    const DEFAULT_NO = 0;

    const DEFAULT_YES = 1;

    protected $table = 'service_merchant_setting';

    protected $guarded = [];
}
