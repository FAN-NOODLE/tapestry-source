<?php
/**
 * Created by PhpStorm.
 * User: wangle@vchangyi.com
 * Date: 2021/04/26
 */

namespace App\Models\Scrm;

class WorkExternalUserRelation extends ScrmInit
{
    /**
     * 表名
     * 成员和外部联系人关系表
     *
     * @var string
     */
    protected $table = 'work_external_user_relation';

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

    public function getWorkExternalUserRelation(int $enterpriseId, int $startStatTime, int $endStatTime, array $memberIds, $fields = ['*'])
    {
        $query = self::query()
            ->select($fields)
            ->where('enterprise_id', $enterpriseId)
            ->whereIn('member_id', $memberIds)
            ->when(!empty($startStatTime), function ($query) use ($startStatTime) {
                return $query->where('createtime', '>=', $startStatTime);
            })
            ->when(!empty($endStatTime), function ($query) use ($endStatTime) {
                return $query->where('createtime', '<=', $endStatTime);
            })
            ->get()
            ->toArray();

        return $query;
    }

    /**
     * @return \Illuminate\Database\Eloquent\Relations\BelongsTo
     * 关联用户表
     */
    public function relationToMember()
    {
        return $this->belongsTo(Members::class,'member_id','id');
    }
}
