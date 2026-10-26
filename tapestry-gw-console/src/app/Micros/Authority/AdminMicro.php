<?php

namespace App\Micros\Authority;

use App\Micros\Common;

class AdminMicro
{
    public static function checkResetPassword($params)
    {
        return Common::authority()->setApi('micro-authority/v1/account/check-reset-password')->post($params);
    }

    public static function resetPassword($params)
    {
        return Common::authority()->setApi('micro-authority/v1/account/reset-password')->post($params);
    }
}

