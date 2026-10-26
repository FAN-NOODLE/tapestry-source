<?php

namespace App\Http\Controllers;

use App\Traits\RequestRecord;
use Illuminate\Support\Facades\Request;
use Laravel\Lumen\Routing\Controller as BaseController;

class Controller extends BaseController
{
    use RequestRecord;

    /**
     * Request
     *
     * @var \Illuminate\Http\Request|null
     */
    protected $request = null;

    public function __construct()
    {
        $this->request = Request::instance();
    }
}
