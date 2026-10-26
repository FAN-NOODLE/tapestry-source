<?php


namespace App\Models;


class StoreEmployee extends Base
{
    /**
     * 门店员工在职、离职状态 1 在职 2 离职
     */
    const STATUS_ENABLE = 1;
    const STATUS_DISABLE = 2;

    CONST SYNC_SUCCESS = 1;


    protected $table = 'store_employee';

    protected $guarded = [];
    
}
