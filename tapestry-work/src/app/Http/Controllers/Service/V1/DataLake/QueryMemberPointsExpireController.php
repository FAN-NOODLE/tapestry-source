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

class QueryMemberPointsExpireController extends DataLakeBaseController
{
    protected $pointTypeGroup;
    protected $Startdate;
    protected $Enddate;

    public function checkParam()
    {
        parent::checkParam();
        $rules   = [
            'pointTypeGroup' => 'required',
            'Startdate'      => 'required',
            'Enddate'        => 'required',
        ];
        $message = [
            'pointTypeGroup.required' => '缺少积分分组类型参数',
            'Startdate.required'      => '缺少开始时间参数',
            'Enddate.required'        => '缺少结束时间参数',
        ];
        $this->validate($this->request, $rules, $message);

        $this->pointTypeGroup = request()->input('pointTypeGroup');
        $this->Startdate      = request()->input('Startdate', 0);
        $this->Enddate        = request()->input('Enddate', 0);
    }

    /**
     * @return array
     * @throws \Throwable
     */
    public function requestApi()
    {
        return $this->dataLakeService->queryPointsExpire(
            $this->type, $this->value, $this->pointTypeGroup, $this->Startdate, $this->Enddate
        );

    }
}