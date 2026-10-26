<?php
/**
 * Created by PhpStorm.
 * Date: 2018/11/15
 * Time: 4:05 PM
 */

namespace App\Models;

use App\Models\ConstTrait\MiniProgram;

class EnterpriseMiniProgramCodeLibrary extends Base
{
    /**
     * 表名称
     *
     * @var string
     */
    protected $table = 'enterprise_mini_program_code_library';

    /**
     * 引入公共常量
     *
     */
    use MiniProgram;

    const PENDING = 'PENDING'; // 待审核/待上线

    const AUDITING = 'AUDITING'; // 审核中
    const AUDITED  = 'AUDITED'; // 审核通过
    const REFUSED  = 'REFUSED'; // 审核失败
    const WITHDRAW = 'WITHDRAW'; // 撤审

    const ONLINE  = 'ONLINE'; // 上线
    const OFFLINE = 'OFFLINE'; // 下线

    /**
     * 自动上线
     */
    const ONLINE_TYPE_AUTOMATIC = 1;

    /**
     * 手动上线
     */
    const ONLINE_TYPE_MANUALLY = 2;

    /**
     * 是否在使用该版本：是
     */
    const USING_STATUS_TRUE = 1;

    /**
     * 是否在使用该版本：否
     */
    const USING_STATUS_FALSE = 0;

    /**
     * 发布模式：全量发布
     */
    const RELEASE_MODEL_ALL = 1;

    /**
     * 发布模式：分阶段发布
     */
    const RELEASE_MODEL_GRAY = 2;

    /**
     * 分阶段发布状态：执行完毕
     */
    const RELEASE_STATUS_COMPLETED = 3;

    /**
     * 分阶段发布状态：被删除
     */
    const RELEASE_STATUS_DELETED = 4;

    /**
     * @var array
     */
    protected $fillable = [
        'id',
        'account_id',
        'version',
        'description',
        'extra_info',
        'template_id',
        'audit_status',
        'online_status',
        'using',
        'online_type',
        'version_id',
        'release_tag',
        'release_branch'
    ];

    /**
     * 关联查询审核日志
     */
    public function auditLog()
    {
        return $this->hasOne(EnterpriseMiniProgramAuditLog::class, 'code_library_id', 'id');
    }
}
