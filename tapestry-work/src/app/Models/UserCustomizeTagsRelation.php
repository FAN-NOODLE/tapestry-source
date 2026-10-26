<?php
/**
 * Description:用户自定义标签关联表
 * User: lintong
 * Email: 18511245965@163.com
 * Date: 2020/8/2
 * Time: 17:37
 */

namespace App\Models;


class UserCustomizeTagsRelation extends Base
{
    /**
     * 表名
     * @var string
     */
    protected $table = 'user_customize_tags_relation';

    /**
     * 主键
     * @var string
     */
    public $primaryKey = 'id';

    /**
     * 可以被批量赋值的属性。
     *
     * @var array
     */
    protected $fillable = ['enterprise_id', 'worker_id', 'user_tags_id', 'union_id'];
}