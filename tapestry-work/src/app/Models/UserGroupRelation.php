<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/11/23
 * Time: 18:25
 */

namespace App\Models;

class UserGroupRelation extends Base
{
    /**
     * 表名
     * @var string
     */
    protected $table = 'user_group_relation';

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
        'ec_userid',
        'group_id'
    ];

    /**
     * @return \Illuminate\Database\Eloquent\Relations\BelongsTo
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/12/21 10:57 AM
     */
    public function userGroup()
    {
        return $this->belongsTo(UserGroup::class, 'group_id');
    }

    /**
     * @return \Illuminate\Database\Eloquent\Relations\HasMany
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/12/21 10:57 AM
     */
    public function members()
    {
        return $this->hasMany(MemberUserRelation::class, 'ec_userid', 'ec_userid');
    }
}
