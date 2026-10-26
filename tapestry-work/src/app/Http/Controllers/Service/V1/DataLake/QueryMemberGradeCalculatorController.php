<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2021/3/25 0025
 * Time: 20:18
 */

namespace App\Http\Controllers\Service\V1\DataLake;

use App\Exceptions\Error;
use Illuminate\Support\Facades\Validator;

class QueryMemberGradeCalculatorController extends DataLakeBaseController
{
    protected $grade;
    protected $upGrade;

    public function checkParam()
    {
        parent::checkParam();
        $rules   = [
            'grade'   => 'required',
            'upGrade' => 'required',
        ];
        $message = [
            'grade.required'   => '缺少当前等级参数',
            'upGrade.required' => '缺少升级等级参数',
        ];
        $this->validate($this->request, $rules, $message);

        $this->grade   = request()->input('grade');
        $this->upGrade = request()->input('upGrade');
    }

    /**
     * @return array
     * @throws \Throwable
     */
    public function requestApi()
    {
        return $this->dataLakeService->queryMemberGradeCalculator($this->type, $this->value, $this->grade, $this->upGrade);

    }
}