<?php

namespace App\Models\Scrm;

class CustomerOnJobInherit extends ScrmInit
{
    const STATUS_ON = 0;
    const STATUS_SUCCESS = 1;
    const STATUS_ERR = 2;
    /**
     * 主键
     *
     * @var string
     */
    public $primaryKey = 'id';

    //分配中
    /**
     * 表名
     *
     * @var string
     */
    protected $table = 'customer_onjob_inherit';

    //分配成功
    /**
     * @var array
     */
    protected $guarded = [];

    //分配失败
    /**
     * @var string[]
     */
    protected $hidden = ['deleted_at'];

    public function record()
    {
        return $this->hasMany(CustomerOnJobInheritRecord::class, 'inherit_id', 'id');
    }

    public function info()
    {
        return $this->hasOne(CustomerOnJobInheritRecord::class, 'inherit_id', 'id');
    }
}