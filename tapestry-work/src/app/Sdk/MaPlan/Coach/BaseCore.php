<?php

namespace App\Sdk\MaPlan\Coach;

use App\Exceptions\Error;

abstract class BaseCore
{

    const MA_PLAN_TYPE=1;
    const MA_STARTDT_TYPE=2;

    protected $enterpriseId;
    protected $originType;

    public function __construct($enterpriseId,$originType=self::MA_PLAN_TYPE)
    {
        if (empty($enterpriseId)) {
            throw new Error(9999, "企业错误");
        }
        $this->enterpriseId = $enterpriseId;
        $this->originType = $originType;
    }
}