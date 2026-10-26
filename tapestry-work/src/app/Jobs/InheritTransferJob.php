<?php

namespace App\Jobs;

use App\Services\CustomerInheritService;

class InheritTransferJob extends Job
{
    public $queue = "inherit_transfer_job";

    private $inherit_id;
    private $transfer_success_msg;
    private $enterprise_id;

    public function __construct($inherit_id, $transfer_success_msg, $enterprise_id)
    {
        $this->inherit_id = $inherit_id;
        $this->transfer_success_msg = $transfer_success_msg;
        $this->enterprise_id = $enterprise_id;
    }

    public function handle()
    {
        (new CustomerInheritService())->transferCustomer($this->inherit_id, $this->transfer_success_msg, $this->enterprise_id);
    }
}