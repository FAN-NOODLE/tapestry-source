<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/11/26
 * Time: 16:24
 */

namespace App\Models\Scrm;

class Members extends ScrmInit
{

    /**
     * 员工激活状态: 0=刚创建 1=已激活，2=已禁用，4=未激活，5=退出企业
     */
    const MEMBERS_STATUS_CREATE = 0;
    const MEMBERS_STATUS_ACTIVATED = 1;
    const MEMBERS_STATUS_DISABLED = 2;
    const MEMBERS_STATUS_NOT_ACTIVE = 4;
    const MEMBERS_STATUS_OUT_ENTERPRISE = 5;

    /**
     * 员工状态：1=在职，2=离职
     */
    const IN_SERVICE  = 1;
    const OUT_SERVICE = 2;

    /**
     * 管理权限标识 3、企业导购，4、店长或企业管理
     */
    const STORE_SA = 3;
    const STORE_SM = 4;

    /**
     * coach职位
     */
    const COACH_STORE_SA = [
        'Part Time Field Sales',
        'Sales AssociatePT',
        'Seasonal Sales Associate',
        'SALES',
        'Sales Associate',
        'Retail Management Trainee SA',
        'Sr Sales Associate',
        'Senior Sales Associate',
        'Senior Sales',
        'Retail Management Trainee SSA'
    ];

    const COACH_STORE_SM = [
        'Store Supervisor',
        'Assistant SM',
        'Store Manager'
    ];

    const COACH_STORE_BOSS = [
        'General Manager',
        'City Manager',
        'Area Manager',
        'Regional Leader'
    ];

    /**
     * 表名
     *
     * @var string
     */
    protected $table = 'members';


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
    public function departmentsRel(): \Illuminate\Database\Eloquent\Relations\HasMany
    {
        return $this->hasMany(MemberDepartments::class, 'member_id', 'id');
    }

    public static function updateMember(array $where, array $data): int
    {
        return self::query()
            ->where($where)
            ->update($data);
    }
}

