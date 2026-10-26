<?php

namespace App\Http\Controllers;

use App\Exceptions\Error;
use Illuminate\Support\Facades\Crypt;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Request;
use \Exception;
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
