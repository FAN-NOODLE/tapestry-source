<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/11/23
 * Time: 18:25
 */

namespace App\Models;

class UserGroup extends Base
{
    /**
     * 表名
     * @var string
     */
    protected $table = 'user_group';

    /**
     * 主键
     * @var string
     */
    public $primaryKey = 'id';

    /**
     * @var array
     */
    protected $fillable = [
        'id',
        'enterprise_id',
        'name',
        'mem_uid'
    ];

    /**
     * @return \Illuminate\Database\Eloquent\Relations\HasMany
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/12/21 10:56 AM
     */
    public function members()
    {
        return $this->hasMany('App\Models\UserGroupRelation', 'group_id');
    }
}
