<?php
/**
 * Created By PhpStorm
 * @Date: 2021/6/25
 * @Time: 9:15
 * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
 * @File: MemberDepartments.php
 */

namespace App\Models\Scrm;


class MemberDepartments extends ScrmInit
{
    /**
     * 表名
     *
     * @var string
     */
    protected $table = 'member_departments';


    /**
     * 主键
     *
     * @var string
     */
    public $primaryKey = 'id';

    /**
     * @var array
     */
    protected $guarded = [];

    //关联表
    public function departments(): \Illuminate\Database\Eloquent\Relations\HasOne
    {
        return $this->hasOne(Departments::class, 'id', 'department_id');
    }

    //关联表
    public function members(): \Illuminate\Database\Eloquent\Relations\HasOne
    {
        return $this->hasOne(Members::class, 'id', 'member_id');
    }
}