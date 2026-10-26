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

class QueryMemberPointsDetailListController extends DataLakeBaseController
{
    protected $pointTypeGroup;
    protected $Startdate;
    protected $Enddate;
    protected $pageNo;
    protected $limit;

    public function checkParam()
    {
        parent::checkParam();
        $rules   = [
            'pointTypeGroup' => 'required',
            'startDate'      => 'required',
            'endDate'        => 'required',
            'pageNo'         => 'required',
        ];
        $message = [
            'pointTypeGroup.required' => '缺少参数',
            'Startdate.required'      => '缺少参数',
            'Enddate.required'        => '缺少参数',
            'pageNo.required'         => '缺少参数',
        ];
        $this->validate($this->request, $rules, $message);

        $this->pointTypeGroup = request()->input('pointTypeGroup');
        $this->Startdate      = request()->input('Startdate');
        $this->Enddate        = request()->input('Enddate');
        $this->pageNo         = request()->input('pageNo');
        //$this->limit          = request()->input('limit', "10");
        $this->limit = '10';
    }

    /**
     * @return array
     * @throws \Throwable
     */
    public function requestApi()
    {
        return $this->dataLakeService->queryPointsList(
            $this->type,
            $this->value,
            $this->pointTypeGroup,
            $this->Startdate,
            $this->Enddate,
            $this->pageNo,
            $this->limit
        );

    }
}