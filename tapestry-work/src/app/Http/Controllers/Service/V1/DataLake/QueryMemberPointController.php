<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2021/3/25 0025
 * Time: 20:17
 */

namespace App\Http\Controllers\Service\V1\DataLake;


use App\Exceptions\Error;

class QueryMemberPointController extends DataLakeBaseController
{
    /**
     * @return array
     * @throws \Throwable
     */
    public function requestApi()
    {
        return $this->dataLakeService->getCurrentPoint($this->type, $this->value);
    }
}