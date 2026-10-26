<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/11/26
 * Time: 16:24
 */

namespace App\Models;

class MemberUserRelationTmp extends Base
{
    /**
     * 是否专属导购：是
     */
    const IS_EXCLUSIVE_TRUE = 1;
    /**
     * 是否专属导购：否
     */
    const IS_EXCLUSIVE_FALSE = 2;

    /**
     * 星标客户 1：是
     */
    const STARS_TRUE = 1;

    /**
     * 星标客户 2：否
     */
    const STARS_FALSE = 2;

    /**
     * 是否是绑定客户：是
     */
    const BINDING_CUSTOMER_YES = 1;
    /**
     * 是否是绑定客户：否
     */
    const BINDING_CUSTOMER_NO = 2;

    /**
     * 表名
     *
     * @var string
     */
    protected $table = 'member_user_relation_tmp';

    /**
     * 主键
     *
     * @var string
     */
    public $primaryKey = 'id';

    // 删除外部联系人
    const DEL_EXTERNAL_CONTACT = 1;

    // 删除跟进成员
    const DEL_FOLLOW_USER = 2;

    // 外部联系人正常状态
    const NORMAL_STATUS = 0;

    /**
     * @var array
     */
    protected $fillable = [
        'id',
        'enterprise_id',
        'store_id',
        'mem_uid',
        'ec_id',
        'ec_userid',
        'ec_name',
        'ec_avatar',
        'ec_gender',
        'user_remark_name',
        'user_name_index',
        'user_description',
        'union_id',
        'relation_at',
        'stars',
        'is_exclusive',
        'source_type',
        'user_id',
        'delete_type',
        'is_binding_customer'
    ];
}
