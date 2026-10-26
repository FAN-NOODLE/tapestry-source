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

class QueryMemberOrderListController extends DataLakeBaseController
{
    protected $page;
    protected $size;
    const DEFAULT_SIZE = 5;

    public function checkParam()
    {
        parent::checkParam();
        $rules   = [
            'page'    => 'integer',
            'size'    => 'integer',
            'tradeNo' => 'string',
        ];
        $message = [
            'page.integer' => '页码为整数',
            'size.integer' => '分页数量为整数',
        ];
        $this->validate($this->request, $rules, $message);

        $this->page    = request()->input('page', 1);
        $this->size    = request()->input('size', static::DEFAULT_SIZE);
        $this->tradeNo = request()->input('tradeNo', '');
    }

    /**
     * @return array
     * @throws \Throwable
     */
    public function requestApi()
    {
        return $this->dataLakeService->getOrderList($this->type, $this->value, $this->page, $this->size, $this->tradeNo);

    }
}