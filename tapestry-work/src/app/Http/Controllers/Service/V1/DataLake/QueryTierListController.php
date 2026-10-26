<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2021/3/31 0031
 * Time: 15:18
 */

namespace App\Http\Controllers\Service\V1\DataLake;


class QueryTierListController extends DataLakeBaseController
{

    /**
     * @return array
     * @throws \Throwable
     */
    public function requestApi()
    {
        return $this->dataLakeService->getTierList();
    }

    /**
     * @throws \Illuminate\Validation\ValidationException
     */
    public function checkParam()
    {
        $rules = [
            'enterprise_id' => 'required'
        ];

        $message = [
            'enterprise_id.required' => '企业id缺失',
        ];
        $this->validate($this->request, $rules, $message);
    }
}