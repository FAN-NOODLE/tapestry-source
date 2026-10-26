<?php
/**
 * Created by PhpStorm.
 * User: wangle@vchangyi.com
 * Date: 2021/04/26
 */

namespace App\Models\Scrm;

class WorkExternalUserTag extends ScrmInit
{

    /**
     * 表名
     * 外部联系人用户标签关系表
     *
     * @var string
     */
    protected $table = 'work_external_user_tag';


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


    /**
     * 用户标签
     *
     * @param $relationIds
     *
     * @return array
     */
    public static function getUserTagsByRelationIds($relationIds): array
    {

        return self::query()
            ->select('work_external_user_tag.relation_id', 'work_external_tag.id', 'work_external_tag.group_id', 'work_external_tag.tag_id', 'work_external_tag.name', 'work_external_tag.type')
            ->leftJoin('work_external_tag',
                'work_external_user_tag.external_tag_id',
                '=',
                'work_external_tag.id')
            ->whereIn('work_external_user_tag.relation_id', $relationIds)
            ->where('work_external_user_tag.deleted_at', 0)
            ->where('work_external_tag.deleted_at', 0)
            ->get()
            ->toArray();

    }


}

