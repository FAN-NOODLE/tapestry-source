<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/11/9
 * Time: 10:00 AM
 */

namespace App\Http\Controllers;

class TestController extends Controller
{
    public function index()
    {
        phpinfo();
    }
}
