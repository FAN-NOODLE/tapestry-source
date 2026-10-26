<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/11/16
 * Time: 9:37 AM
 */

namespace App\Http\Controllers\Service\V1;

use App\Http\Controllers\Controller;

abstract class BaseController extends Controller
{
    /**
     * BaseController constructor.
     *
     * @throws \Exception
     */
    public function __construct()
    {
        parent::__construct();
    }
}
