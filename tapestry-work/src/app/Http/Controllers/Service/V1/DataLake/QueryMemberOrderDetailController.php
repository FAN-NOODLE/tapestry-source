<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2021/3/25 0025
 * Time: 20:19
 */

namespace App\Http\Controllers\Service\V1\DataLake;

use App\Exceptions\Error;
use Illuminate\Support\Facades\Validator;

class QueryMemberOrderDetailController extends DataLakeBaseController
{
    protected $tradeNo;

    public function checkParam()
    {
        parent::checkParam();
        $rules = [
            'trade_no'  => 'required',
        ];

        $message = [
            'trade_no.required' => '交易流水号为空',
        ];
        $this->validate($this->request, $rules, $message);

        $this->tradeNo = request()->input('trade_no');
    }

    /**
     * @return array
     * @throws \Throwable
     */
    public function requestApi()
    {
        return $this->dataLakeService->getOrderDetail($this->type,$this->value,$this->tradeNo);
    }
}