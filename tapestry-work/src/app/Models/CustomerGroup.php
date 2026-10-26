<?php


namespace App\Models;


class CustomerGroup extends Base
{
    protected $table = 'customer_group';

    protected $guarded = [];

    /**
     * 插入数据
     * @param $input
     * @return \Illuminate\Database\Eloquent\Builder|\Illuminate\Database\Eloquent\Model
     */
    public function insertCustomerGroup($input)
    {
        return self::query()->create($input);
    }

}