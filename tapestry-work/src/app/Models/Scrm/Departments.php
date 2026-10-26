<?php
/**
 * Created By PhpStorm
 * @Date: 2021/6/25
 * @Time: 9:15
 * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
 * @File: Departments.php
 */

namespace App\Models\Scrm;


use Illuminate\Database\Eloquent\Relations\HasMany;

class Departments extends ScrmInit
{

    CONST CORPORATION_DEPT_ID_KEY = 'corporation_dept_id';
    /**
     * 表名
     *
     * @var string
     */
    protected $table = 'departments';


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
    public function membersRel(): HasMany
    {
        return $this->hasMany(MemberDepartments::class, 'department_id', 'id');
    }

    //关联表
    public function child(): HasMany
    {
        return $this->hasMany(Departments::class, 'third_parent_id', 'third_id');
    }
}