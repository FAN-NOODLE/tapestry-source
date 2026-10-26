<?php


namespace App\Http\Controllers;

/**
 * Class HealthController
 * @package App\Http\Controllers
 */
class HealthController
{
    /**
     * 健康检查
     * @return string
     */
    public function index()
    {
        return 'OK';
    }
}