<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/11/26
 * Time: 16:24
 */

namespace App\Models;

class GuideUserRelation extends Base
{
    /**
     * 是否复制手机号-否
     * @var int
     */
    const IS_COPY_MOBILE_NO = 0;

    /**
     * 是否复制手机号-是
     * @var int
     */
    const IS_COPY_MOBILE_YES = 1;

    /**
     * 表名
     *
     * @var string
     */
    protected $table = 'guide_user_relation';

    /**
     * 主键
     *
     * @var string
     */
    public $primaryKey = 'id';

    /**
     * @var array
     */
    protected $fillable = [
        'id',
        'enterprise_id',
        'mem_uid',
        'guide_id',
        'union_id',
        'user_id',
        'type',
        'relation_at',
        'created_at',
        'updated_at',
    ];

    /**
     * @var array
     */
    public static $selectField = ['id', 'enterprise_id', 'mem_uid', 'guide_id',
        'union_id', 'user_id', 'type', 'relation_at'];
}
