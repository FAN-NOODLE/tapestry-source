<?php


namespace App\Http\Controllers\Service\V1\UserAction;

use App\Http\Controllers\Controller;

class OpenController extends Controller
{
    public function index()
    {
        $params    = $this->request->all();
        return $params;
    }
}