<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2021/3/31 0031
 * Time: 21:07
 */

namespace App\Http\Controllers\Service\V1\DataLake;


class QueryBaListController extends DataLakeBaseController
{
    protected $baCode;
    protected $page;
    protected $size;
    const DEFAULT_SIZE = 5;
    /**
     * @return array
     * @throws \Throwable
     */
    public function requestApi()
    {
        return $this->dataLakeService->queryBAList($this->baCode,$this->page,$this->size);
    }

    /**
     * @throws \Illuminate\Validation\ValidationException
     */
    public function checkParam()
    {
        $rules = [
            'enterprise_id' => 'required',
            'ba_code'       => 'required'
        ];

        $message = [
            'enterprise_id.required' => '企业id缺失',
            'ba_code.required'       => '导购code缺失',
        ];
        $this->validate($this->request, $rules, $message);
        $this->baCode = request()->input('ba_code');
        $this->page = request()->input('page',1);
        $this->size = request()->input('size',self::DEFAULT_SIZE);
    }
}