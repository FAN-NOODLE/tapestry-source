<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Builder;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Support\Facades\DB;
use App\Services\EnterpriseAdminService;

class EnterpriseAdmin extends Base
{
    /**
     * 是否超管  1-是  2-不是
     */
    const IS_SUPER = 1;
    const NO_SUPER = 2;

    /**
     * 1-不限 2-指定范围
     */
    const RANGE_TYPE_ALL  = 1;
    const RANGE_TYPE_ONLY = 2;

    /**
     * 账号状态 1-开启 2-禁用
     */
    const STATUS_OPEN  = 1;
    const STATUS_CLOSE = 2;

    /**
     *  1 CRS后台 2 店长小程序 3 天擎CDP
     */
    const SYSTEM_CRS   = 1;
    const SYSTEM_STORE = 2;
    const SYSTEM_CDP   = 3;

    const  CLIENT_CRS   = 'console';
    const  CLIENT_STORE = 'store';
    const  CLIENT_CDP   = 'cdp';

    const CLIENT_MAP = [
        self::CLIENT_CRS   => self::SYSTEM_CRS,
        self::CLIENT_STORE => self::SYSTEM_STORE,
    ];

    /**
     * 密码是否修改 0-未修改密码
     */
    const NOT_UPDATE_PASSWORD = 0;

    /**
     * @var int 默认密码过期时间 90天
     */
    const PASSWORD_EXPIRE = 3600 * 24 * 90;


    protected $table = 'enterprise_admin';

    protected $guarded = [];

    /**
     * 获取管理员通过ID
     * @param $id
     * @return EnterpriseAdmin|Builder|Model|object|null
     */
    public static function adminById($id)
    {
        return self::query()
            ->where('id', '=', $id)
            ->first();
    }

    /**
     * 获取管理员列表通过IDs
     * @param $eId
     * @param $ids
     * @return EnterpriseAdmin[]|Builder[]|\Illuminate\Database\Eloquent\Collection
     */
    public static function adminListByIds($eId, $ids)
    {
        return self::query()
            ->select(['username', 'mobile', 'is_super', 'range_type', 'admin_available_system'])
            ->where('enterprise_id', '=', $eId)
            ->where('status', '=', self::STATUS_OPEN)
            ->whereIn('id', $ids)
            ->get();
    }

    // 获取账号详情
    public function accountDetail($where, $fields = '*', $account)
    {
        $obj = self::query()
            ->select($fields)
            ->where($where)
            ->when(!empty($account), function ($query) use ($account) {
                return $query->orWhere('mobile', $account);
            })
            ->first();

        return empty($obj) ? [] : $obj->toArray();
    }

    // 获取账号详情
    public function accountDetailByLogin($account, $fields = '*', $loginType = '')
    {
        $obj = self::query()
            ->select($fields)
            ->when($loginType, function ($query) use ($loginType, $account) {
                if ($loginType == EnterpriseAdminService::LOGIN_BY_PASSWORD) {
                    $query->where('account', 'like', '%' . $account);

                } else if ($loginType == EnterpriseAdminService::LOGIN_BY_CODE) {
                    $query->where('mobile', $account);
                }
            })
            ->first();

        return empty($obj) ? [] : $obj->toArray();
    }

    public function accountDetailByMobile($mobile, $enterpriseId, $fields = '*')
    {
        $obj = self::query()
            ->select($fields)
            ->rightJoin('short_message_config', 'enterprise_admin.enterprise_id', '=', 'short_message_config.enterprise_id')
            ->when(!empty($enterpriseId), function ($query) use ($enterpriseId) {
                return $query->where('enterprise_admin.enterprise_id', $enterpriseId);
            })
            ->where('enterprise_admin.mobile', $mobile)
            ->first();

        return empty($obj) ? [] : $obj->toArray();
    }

    /**
     * 根据账号查询账号信息
     * @param $account
     * @param $enterpriseId
     * @param string $fields
     * @return array
     * @Date: 2021/4/30
     * @Time: 17:44
     * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
     */
    public function accountDetailByAccount($account, $enterpriseId, $fields = '*')
    {
        $obj = self::query()
            ->select($fields)
            ->where('enterprise_id', $enterpriseId)
            ->where('account', $account)
            ->first();

        return empty($obj) ? [] : $obj->toArray();
    }

    /**
     * 根据手机号查询账号信息
     * @param $mobile
     * @param $enterpriseId
     * @param string $fields
     * @return array
     * @Date: 2021/4/30
     * @Time: 17:44
     * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
     */
    public function accountDetailByPhone($mobile, $enterpriseId, $fields = '*')
    {
        $obj = self::query()
            ->select($fields)
            ->when(!empty($enterpriseId), function ($query) use ($enterpriseId) {
                return $query->where('enterprise_id', $enterpriseId);
            })
            ->where('mobile', $mobile)
            ->first();

        return empty($obj) ? [] : $obj->toArray();
    }

    protected function getSelectFields()
    {
        return [
            'enterprise.id as enterprise_id',
            'enterprise.hash',
            'enterprise.version',
            'enterprise.version_tag',
            'enterprise.name',
            'enterprise.contact as contact_name',
            'enterprise.mobile as contact_mobile',
            'enterprise.address',
            'enterprise.crs_status',
            'enterprise.store_status',
            'enterprise.contact_status',
            'enterprise.domain',
            'enterprise.corp_id',
            'enterprise.corp_id',
            'enterprise.crs_status',
            'enterprise.account_start_at',
            'enterprise.account_expire_at',
            'enterprise.is_split_flow',
            'enterprise_admin.id as admin_id',
            'enterprise_admin.account',
            'enterprise_admin.status as account_status',
            'enterprise_admin.username as account_name',
            'enterprise_admin.global_roaming as account_global_roaming',
            'enterprise_admin.mobile as account_mobile',
            'enterprise_admin.email as account_email',
            'enterprise_admin.role',
            'enterprise_admin.is_super',
            'enterprise_admin.range_type',
            'enterprise_admin.admin_available_system',
        ];
    }

    public function enterpriseAdminList($account, $client, $version = '', $epIds = [])
    {
        $select = $this->getSelectFields();

        if (!empty($epIds)) {
            if (!is_array($epIds)) {
                $epIds = [$epIds];
            }
        }

        $query = self::query()
            ->select($select)
            ->leftJoin('enterprise', 'enterprise.id', '=', 'enterprise_admin.enterprise_id')
            ->where('enterprise_admin.deleted_at', 0)
            ->where(function ($q) use ($account) {
                $q->where('enterprise_admin.account', 'like', '%' . $account)
                    ->orWhere('enterprise_admin.mobile', 'like', '%' . $account . '%');
            })
            ->when(!empty($version), function ($query) use ($version) {
                if ($version == 'v1') {
                    return $query->where('enterprise.version', 1);
                } else {
                    return $query->where('enterprise.version', '<>', 1);

                }
            });

        if (!empty($epIds)) {
            $query->whereIn('enterprise_id', $epIds);
        }
        //如果是超管，优先以超管身份登录
        $query = $query->get();
        if ($query->where('is_super', 1)->count() && $client != 'store') {
            return $query->where('is_super', 1)->toArray();
        }
        return $query->toArray();
    }

    public function enterpriseAdminDetail($adminId)
    {
        $select = $this->getSelectFields();

        $obj = self::query()
            ->select($select)
            ->leftJoin('enterprise', 'enterprise.id', '=', 'enterprise_admin.enterprise_id')
            ->where('enterprise_admin.id', $adminId)
            ->first();

        return empty($obj) ? [] : $obj->toArray();
    }

    public function enterpriseAdminRoleList($adminId)
    {
        $select = [
            'role.id as role_id',
            'role.name as role_name',
            'role.mark as role_mark',
        ];

        return EnterpriseAdminRole::query()
            ->select($select)
            ->leftJoin('role', 'enterprise_admin_role.role_id', '=', 'role.id')
            ->where('enterprise_admin_role.enterprise_admin_id', $adminId)
            ->get()
            ->toArray();
    }

    public static function adminUsername($id)
    {
        return self::query()->where('id', '=', $id)->value('username') ?? '系统管理员';
    }

    /**
     * 通过手机号和企业ID 获取管理员
     * @param $account
     * @param $eId
     * @return EnterpriseAdmin|Builder|Model|object|null
     */
    public static function adminByAccount($account, $eId)
    {
        return self::query()
            ->where('enterprise_id', '=', $eId)
            ->where('account', '=', "+86-{$account}")
            ->first();
    }

    public static function adminByMobile($mobile, $eId)
    {
        return self::query()
            ->where('enterprise_id', '=', $eId)
            ->where('mobile', '=', $mobile)
            ->first();
    }

    /**
     * 关联账号角色表
     *
     * @return \Illuminate\Database\Eloquent\Relations\HasMany
     */
    public function roleRelation()
    {
        return $this->hasMany(EnterpriseAdminRole::class, 'enterprise_admin_id', 'id');

    }

    /**
     * 关联账号范围表
     *
     * @return \Illuminate\Database\Eloquent\Relations\HasMany
     */
    public function RangeRelation()
    {
        return $this->hasMany(EnterpriseDataRange::class, 'enterprise_admin_id', 'id');
    }

    /**
     * 关联角色
     * @return \Illuminate\Database\Eloquent\Relations\BelongsToMany
     */
    public function adminRoleRelation()
    {
        return $this->belongsToMany(Role::class, 'admin_role_relation', 'admin_id', 'role_id');
    }
}
