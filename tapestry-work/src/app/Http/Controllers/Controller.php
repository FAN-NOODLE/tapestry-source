<?php

namespace App\Http\Controllers;

use Illuminate\Support\Facades\Request;
use Laravel\Lumen\Routing\Controller as BaseController;

class Controller extends BaseController
{
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
