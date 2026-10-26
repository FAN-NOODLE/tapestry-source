<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2021/3/25 0025
 * Time: 19:54
 */

namespace App\Http\Controllers\Service\V1\DataLake;

use App\Exceptions\Error;
use Illuminate\Support\Facades\Validator;

class QueryMemberInfoController extends DataLakeBaseController
{
    /**
     * @return array
     * @throws \Throwable
     */
    public function requestApi()
    {
        return $this->dataLakeService->queryMemberInfo($this->type,$this->value);
    }
}