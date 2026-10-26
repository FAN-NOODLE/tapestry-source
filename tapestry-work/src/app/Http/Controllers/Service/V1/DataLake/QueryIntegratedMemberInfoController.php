<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2021/3/29 0029
 * Time: 12:49
 */

namespace App\Http\Controllers\Service\V1\DataLake;


class QueryIntegratedMemberInfoController extends DataLakeBaseController
{
    /**
     * 整合用户信息接口
     *
     * @return array
     * @throws \Throwable
     */
    public function requestApi()
    {
        return  [
            'member_info'=>$this->dataLakeService->queryMemberInfo($this->type,$this->value),
            'grade'=>$this->dataLakeService->queryGradeInfo($this->type,$this->value),
            'point'=>$this->dataLakeService->getCurrentPoint($this->type,$this->value)
        ];
    }
}