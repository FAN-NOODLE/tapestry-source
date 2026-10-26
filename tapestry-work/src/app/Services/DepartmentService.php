<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/12/3
 * Time: 18:33
 */

namespace App\Services;

use App\Library\Uc\Department;

class DepartmentService
{
    /**
     * 获取门店dptId
     *
     * @param $qyDomain      企业domain
     * @param $appIdentifier 企业应用标识
     *
     * @return string
     * @throws \Exception
     */
    public function getStoreDptId($qyDomain, $appIdentifier)
    {
        // 调用sdk获取部门类型
        $typeList = (new Department())->typeList($qyDomain, $appIdentifier);
        $typeList = array_column($typeList, null, 'dptName');

        return $typeList['门店']['dptId'] ?? '';
    }
}
