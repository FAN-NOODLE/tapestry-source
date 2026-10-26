<?php

namespace App\Sdk\Active\CRM\CoachCore;


abstract class SwBase extends Base
{
    protected $prefix = "/v2/cdp";

    public function __construct()
    {
        parent::__construct();
        if (ConfigSW::isProdEnv()) {
            $this->prefix = "/cdp-api-sw";
        }
    }

}