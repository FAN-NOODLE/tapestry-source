<?php

namespace App\Services;

class SystemService
{
    /**
     * Shop分流
     *
     * @param $enterpriseHash string 企业Hash
     */
    public static function shopSplitFlow($enterpriseHash)
    {
        $_SERVER['enterprise_hash'] = $enterpriseHash;
    }
}
