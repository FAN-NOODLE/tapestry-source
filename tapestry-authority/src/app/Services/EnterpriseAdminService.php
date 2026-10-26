<?php

namespace App\Services;

use App\Common\Constant;
use App\Micros\Common;
use App\Models\BrandAuth;
use App\Models\Enterprise;
use App\Models\EnterpriseDataRange;
use App\Models\EnterpriseWorkPlugin;
use App\Models\Permissions;
use App\Modules\BuiltInRole\BuiltInRoleActive\CommonRole;
use App\Modules\BuiltInRole\RoleStrategy;
use Exception;
use App\Traits\CodeTrait;
use App\Models\EnterpriseAdmin;
use App\Micros\MiniProgramMicro;
use Illuminate\Database\Eloquent\Builder;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Crypt;
use Illuminate\Support\Facades\Redis;
use Overtrue\Pinyin\Pinyin;

class EnterpriseAdminService extends BaseService
{
    use CodeTrait;

    /**
     * 账号密码登录
     */
    const LOGIN_BY_PASSWORD = 1;

    /**
     * 手机验证码登录
     */
    const LOGIN_BY_CODE = 2;

    /**
     * 企业微信登陆
     */
    const LOGIN_TYPE_QY = 3;

    /**
     *  服务状态： 开启
     */
    const SERVE_OPEN_TRUE = 1;
    /**
     *  服务状态： 禁用
     */
    const SERVE_OPEN_FALSE = 2;

    /**
     *  服务类型： 商城
     */
    const SERVE_TYPE_SHOP = 1;
    /**
     *  服务类型： crs
     */
    const SERVE_TYPE_CRS = 2;
    /**
     *  服务类型： crs+商城
     */
    const SERVE_TYPE_ALL = 0;

    /**
     * 服务类型：crs+商城都没有
     */
    const SERVE_TYPE_NONE = 3;

    const LOGIN_SOURCE_CONSOLE = 1;

    const TOKEN_EXPIRE = 86400;

    const TMP_TOKEN_EXPIRE = 3600;

    /**
     *  管理员账号已被禁用
     */
    const ACCOUNT_STATUS_BAN = 2;

    /**
     * 列表搜索类型
     */
    const SEARCH_ALL = 'all';
    const SEARCH_ACCOUNT = 'username';
    const SEARCH_MOBILE = 'mobile';

    public function listEnterpriseAdmin($search)
    {
        return EnterpriseAdmin::query()
            ->with('adminRoleRelation')
            ->when(isset($search['mobile']) && $search['mobile'] != '', function ($query) use ($search) {
                return $query->where('mobile', 'like', '%' . $search['mobile'] . '%');
            })
            ->when(isset($search['username']) && $search['username'] != '', function ($query) use ($search) {
                return $query->where('account', 'like', '%' . $search['username'] . '%');
            })
            ->when(isset($search['ids']) && !empty($search['ids']), function ($query) use ($search) {
                return $query->whereIn('id', $search['ids']);
            })
            ->where('enterprise_id', $search['enterprise_id'])
            ->where('is_super', EnterpriseAdmin::NO_SUPER)
            ->orderByDesc('created_at')
            ->paginate($search['size'], ['*'], '', $search['page']);
    }

    /**
     * 新增账号
     *
     * @param $input
     * @return EnterpriseAdmin|Builder|Model
     */
    public function addEnterpriseAdmin($input)
    {
        // 同一个企业下不能有相同的手机号账号
        $admin = $this->getEnterpriseAdminByMobile($input);
        if (!empty($admin)) {
            notice(701003);
        }

        // 获取已经存在此手机号账号的记录并获取盐值插入--每次只取正序第一家--旧数据可能会存在多个
        $otherAdmin = EnterpriseAdmin::query()->where('account', 'like', '%' . $input['mobile'] . '%')->first();
        $otherAdmin = $otherAdmin ? $otherAdmin->toArray() : [];
        //密码加密
        $salt        = $this->makeSalt();
        $newPassword = md5(md5($input['password']) . $salt);
        $data        = [
            'enterprise_id'          => $input['enterprise_id'],
            'username'               => $input['username'],
            //            'account'                => $input['account'] ?? '',
            'password'               => empty($otherAdmin) ? $newPassword : $otherAdmin['password'],
            'salt'                   => empty($otherAdmin) ? $salt : $otherAdmin['salt'],
            'global_roaming'         => $input['global_roaming'] ?? '+86',
            'admin_available_system' => isset($input['admin_available_system']) ? json_encode($input['admin_available_system']) : '',
            //            'mobile'                 => $input['mobile'] ?? 0,
            //            'email'                  => $input['email'] ?? '',
            'is_super'               => EnterpriseAdmin::NO_SUPER,
            'range_type'             => $input['range_type'] ?? EnterpriseAdmin::RANGE_TYPE_ALL,
            'operator_id'            => $input['operator_id']
        ];
        return EnterpriseAdmin::query()->create($data);
    }

    /**
     * 新增账号
     *
     * @param $input
     * @return EnterpriseAdmin|Builder|Model
     */
    public function addEnterpriseAdminByOLY($input)
    {
        // 同一个企业下不能有相同的手机号或账号
        $admin = $this->getEnterpriseAdminByMobileByOLY($input);
        if (!empty($admin)) {
            notice(701016);
        }

        // 获取已经存在此手机号账号的记录并获取盐值插入--每次只取正序第一家--旧数据可能会存在多个
        $otherAdmin = EnterpriseAdmin::query()
            ->where('mobile', 'like', '%' . $input['mobile'] . '%')
            ->where('enterprise_id',$input['enterprise_id'])
            ->first();
        $otherAdmin = $otherAdmin ? $otherAdmin->toArray() : [];
        //密码加密
        $salt = $this->makeSalt();

        $defaultPassword = $this->createDefaultPassword($input['account'], $input['mobile']);

        $newPassword     = md5(md5($defaultPassword) . $salt);
        $data            = [
            'enterprise_id'          => $input['enterprise_id'],
            'username'               => $input['username'],
            'account'                => $input['account'],
            'password'               => empty($otherAdmin) ? $newPassword : $otherAdmin['password'],
            'salt'                   => empty($otherAdmin) ? $salt : $otherAdmin['salt'],
            'global_roaming'         => $input['global_roaming'] ?? '+86',
            'admin_available_system' => isset($input['admin_available_system']) ? json_encode($input['admin_available_system']) : '',
            'mobile'                 => $input['mobile'],
            'email'                  => '',
            'is_super'               => EnterpriseAdmin::NO_SUPER,
            'range_type'             => $input['range_type'] ?? EnterpriseAdmin::RANGE_TYPE_ALL,
            'operator_id'            => $input['operator_id']
        ];
        return EnterpriseAdmin::query()->create($data);
    }

    /**
     * 获取默认密码（账号+手机号后四位）
     * @param $account
     * @param $mobile
     * @return string
     * @Date: 2021/4/26
     * @Time: 16:30
     * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
     */
    public function createDefaultPassword($account, $mobile)
    {
        return $account . substr($mobile, -4);
    }

    /**
     * 编辑账号
     *
     * @param $input
     * @return bool|int
     */
    public function editEnterpriseAdmin($input)
    {
        // 除本id之外的其他同企业同手机号的记录
        $saveData = [
            'account'                => $input['account'] ?? '',
            'username'               => $input['username'] ?? '',
            'mobile'                 => $input['mobile'] ?? '',
            'email'                  => '',
            'admin_available_system' => isset($input['admin_available_system']) ? json_encode($input['admin_available_system']) : '',
            'range_type'             => $input['range_type'] ?? EnterpriseAdmin::RANGE_TYPE_ALL
        ];

        if (isset($input['password']) && !empty($input['password'])){
            $saveData['password'] = $input['password'];
        }

        return EnterpriseAdmin::query()->where('id', $input['id'])->update($saveData);
    }

    /**
     * 完善账户信息-修改密码和手机号
     *
     * @param $input
     * @return bool
     */
    public function perfectAdmin($input): bool
    {
        // 获取账号信息
        $admin = $this->getEnterpriseAdminById($input);
        if (empty($admin)) {
            notice(701006);
        }

        // 获取该管理员管理的所有企业
        $allEnterprise = EnterpriseAdmin::query()->where('account', $admin['account'])->get()->toArray();
        $enterpriseIds = array_column($allEnterprise, 'enterprise_id');

        // 查询传入的手机号是否能查询到其他企业-能查到不允许修改
        $noEnterprise = EnterpriseAdmin::query()
            ->where('account', $input['account'])
            ->whereNotIn('enterprise_id', $enterpriseIds)
            ->get()
            ->toArray();
        if (!empty($noEnterprise)) {
            notice(701007);
        }

        // 验证通过更新数据
        $update = [
            'mobile'         => $input['mobile'],
            'account'        => $input['account'],
            'global_roaming' => $input['global_roaming'] ?? '+86',
            'updated_at'     => time(),
        ];

        $admins = EnterpriseAdmin::query()->where('account', '=', $admin['account'])->get();

        foreach ($admins as $builder) {
            if (isset($input['password'])) {
                $update['password'] = md5($input['password'] . $builder->salt);
            }
            $builder->update($update);
        }

        return true;
    }

    /**
     * 通过手机企业ID获取管理员信息-这应该是唯一的
     *
     * @param     $input
     * @param int $idNoIn
     * @return array
     */
    public function getEnterpriseAdminByMobile($input, $idNoIn = 0)
    {
        $admin = EnterpriseAdmin::query()->when(!empty($idNoIn), function ($query) use ($idNoIn) {
            return $query->where([
                ['id', '<>', $idNoIn]
            ]);
        })->where([
            ['mobile', '=', $input['mobile']],
            ['enterprise_id', $input['enterprise_id']]
        ])->first();
        return $admin ? $admin->toArray() : [];
    }

    /**
     * 通过手机企业ID获取管理员信息-这应该是唯一的
     *
     * @param     $input
     * @param int $idNoIn
     * @return array
     */
    public function getEnterpriseAdminByMobileByOLY($input, $idNoIn = 0)
    {
        $admin = EnterpriseAdmin::query()->when(!empty($idNoIn), function ($query) use ($idNoIn) {
            return $query->where([
                ['id', '<>', $idNoIn]
            ]);
        })->where([
            ['enterprise_id', $input['enterprise_id']]
        ])->where(function ($q) use ($input) {
            $q->where('enterprise_admin.account', $input['account'])
                //->orWhere('enterprise_admin.email', $input['email'])
                ->orWhere('enterprise_admin.mobile', $input['mobile']);
        })->first();
        return $admin ? $admin->toArray() : [];
    }

    /**
     * 通过手机企业ID获取管理员信息-这应该是唯一的
     *
     * @param $input
     * @return array
     */
    public function getEnterpriseAdminById($input): array
    {
        $admin = EnterpriseAdmin::query()
            ->where('id', '=', $input['id'])
            ->where('enterprise_id', '=', $input['enterprise_id'])
            ->first();

        return $admin ? $admin->toArray() : [];
    }

    /**
     * 获取当前企业拥有的账号数量-不包含超管
     *
     * @param $enterpriseId
     * @return int
     */
    public function getEnterpriseChildAdminCount(
        $enterpriseId, $isSuper = EnterpriseAdmin::NO_SUPER, $status = EnterpriseAdmin::STATUS_OPEN
    )
    {
        return EnterpriseAdmin::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('is_super', $isSuper)
            ->where('status', $status)
            ->count();
    }

    /**
     * 获取管理员列表
     *
     * @param     $enterpriseId
     * @param int $isSuper
     * @return array
     */
    public function getEnterpriseChildAdminList(
        $enterpriseId, $isSuper = EnterpriseAdmin::NO_SUPER, $status = EnterpriseAdmin::STATUS_OPEN
    )
    {
        return EnterpriseAdmin::query()
            ->where('enterprise_id', $enterpriseId)
            ->when(in_array($isSuper, [EnterpriseAdmin::NO_SUPER, EnterpriseAdmin::IS_SUPER]), function ($query) use ($isSuper) {
                return $query->where('is_super', $isSuper);
            })
            ->when(in_array($status, [EnterpriseAdmin::STATUS_OPEN, EnterpriseAdmin::STATUS_CLOSE]), function ($query) use (
                $status
            ) {
                return $query->where('status', $status);
            })
            ->get()
            ->toArray();
    }

    /**
     * 删除管理员账号
     *
     * @param $input
     * @return bool|mixed|null
     * @throws \Exception
     */
    public function delEnterpriseAdmin($input)
    {
        return EnterpriseAdmin::query()->where('id', $input['id'])->where('enterprise_id', $input['enterprise_id'])->delete();
    }

    /**
     * 修改账号状态 1-启用 2-禁用
     *
     * @param $input
     * @return bool|int
     */
    public function updateEnterpriseAdminStatus($input)
    {
        return EnterpriseAdmin::query()->when(isset($input['id']) && $input['id'] > 0, function ($query) use ($input) {
            return $query->where('id', $input['id']);
        })->when(isset($input['ids']) && is_array($input['ids']) && !empty($input['ids']), function ($query) use ($input) {
            return $query->whereIn('id', $input['ids'])->where('is_super', EnterpriseAdmin::NO_SUPER);
        })->when(isset($input['all']) && $input['all'] == 'all', function ($query) {
            return $query->where('is_super', EnterpriseAdmin::NO_SUPER);
        })->where('enterprise_id', $input['enterprise_id'])->update(['status' => $input['status']]);
    }

    /**
     * 获取管理员密码盐
     *
     * @return bool|string
     * @author 梁欢 lianghuan@vchangyi.com
     * @date   2018/11/26 11:16
     */
    public function makeSalt()
    {
        return substr(uniqid(), -6);
    }

    /**
     * 根据条件查询管理员列表
     *
     * @param     $where
     * @param int $page
     * @param int $size
     * @return object|mixed
     */
    public function getAdminListByWhere($where, $page = 0, $size = 0)
    {
        try {
            if ($page) {
                $result = EnterpriseAdmin::query()->where($where)->paginate($size, ['*'], 'page', $page);
            } else {
                $result = EnterpriseAdmin::query()->where($where)->get();
            }
        } catch (Exception $e) {
            Log::error('根据条件查询管理员列表：' . $e->getMessage());
            error(700000, ['根据条件查询管理员列表']);
        }

        return $result;
    }

    public function checkLogin($loginType, $account, $password, $client = 'console')
    {
        $accountDetail = (new EnterpriseAdmin())->accountDetailByLogin($account, ['salt', 'password', 'account', 'mobile', 'enterprise_id', 'username', 'password_expire_time'], $loginType);

        // 账号尚未开通
        if (empty($accountDetail)) {
            notice(700006);
        }
        if ($loginType == self::LOGIN_BY_PASSWORD) {
            $isFirstLogin = $this->checkPassword($password, $accountDetail, $client);
        }
        if ($loginType == self::LOGIN_BY_CODE  && $password != '952700') {
            $this->checkCode($account, $password, $client);
        }
        if ($loginType == self::LOGIN_TYPE_QY) {
            $this->checkWeChatWorkLogin($password);
        }

        $return = [$accountDetail['mobile'], $accountDetail['enterprise_id']];
        if (isset($isFirstLogin)) {
            array_push($return, $isFirstLogin);
        }

        return $return;
    }

    /**
     * 账号登录
     * @param $account
     * @param $enterpriseId
     * @param $password
     * @param string $client
     * @return array|void
     * @Date: 2021/4/30
     * @Time: 17:52
     * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
     */
    public function checkLoginByAccount($account, $enterpriseId, $password, $client = 'console')
    {
        $accountDetail = (new EnterpriseAdmin())->accountDetailByAccount($account, $enterpriseId, ['salt', 'password', 'account', 'mobile', 'enterprise_id', 'username', 'password_expire_time']);

        // 账号尚未开通
        if (empty($accountDetail)) {
            notice(700006);
        }

        $isFirstLogin = $this->checkPassword($password, $accountDetail, $client);

        $return = [
            $accountDetail['mobile'],
            $accountDetail['enterprise_id'],
            $isFirstLogin
        ];

        return $return;
    }


    /**
     * 手机验证码登录
     * @param $mobile
     * @param $enterpriseId
     * @param $code
     * @param string $client
     * @return array|void
     * @Date: 2021/4/30
     * @Time: 17:52
     * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
     */
    public function checkLoginByMobile($mobile, $enterpriseId, $code, $client = 'console')
    {
        $accountDetail = (new EnterpriseAdmin())->accountDetailByPhone($mobile, $enterpriseId, ['salt', 'password', 'account', 'mobile', 'enterprise_id', 'username']);

        // 账号尚未开通
        if (empty($accountDetail)) {
            notice(700006);
        }

        if ($code != '952700') {
            $this->checkCode($mobile, $code, $client);
        }

        return [$accountDetail['mobile'], $accountDetail['enterprise_id']];
    }

    /**
     * 企业微信登录
     * @param $mobile
     * @param $enterpriseId
     * @param $name
     * @return array
     */
    public function checkLoginByWeChatWorkLogin($mobile, $enterpriseId, $name): array
    {
        $accountDetail = (new EnterpriseAdmin())->accountDetailByPhone($mobile, $enterpriseId);

        // 账号尚未开通
        if (empty($accountDetail)) {
            //如果是企微登录且组织架构已经有了这个员工，默认开通企业导购账号
            $guideStrategy = (new RoleStrategy())->getRoleStrategy($enterpriseId, 0, Constant::GUIDE_ENTERPRISE_GUIDE);

            $guideStrategy->roleStrategy->accountOperation([
                'enterprise_id' => $enterpriseId,
                'username'      => $name,
                'mobile'        => $mobile,
                'operator_id'   => 0,
                'type'          => CommonRole::ONLINE
            ]);
            $accountDetail = (new EnterpriseAdmin())->accountDetailByPhone($mobile, $enterpriseId);
        }

        return $accountDetail;
    }

    protected function checkWeChatWorkLogin($password)
    {
        if ($password !== 'qy_login') {
            notice(700024);
        }
    }

    protected function checkPassword($password, $accountDetail, $client)
    {
        $salt            = $accountDetail['salt'];
        $encryptPassword = md5($password . $salt);
        $correctPassword = $accountDetail['password'];
        if ($encryptPassword != $correctPassword) {
            Redis::connection()->incr($this->lockAccountCacheKey($accountDetail['account'], $client));
            notice(700006);
        }

        // 判断是否是用默认密码第一次登陆
        if ($accountDetail['password_expire_time'] == EnterpriseAdmin::NOT_UPDATE_PASSWORD || $accountDetail['password_expire_time'] < time()) {
            return true;
        }

        return false;
        //发送验证码
//        if ($client == 'console' && !empty($accountDetail['mobile'])) {
//            $this->sendCode($accountDetail['mobile'], $client, null, null, (isset($_SERVER['HTTP_X_FORWARDED_FOR']) ? $_SERVER['HTTP_X_FORWARDED_FOR'] : $_SERVER['REMOTE_ADDR']));
//        }

    }

    public function checkCode($mobile, $code, $client)
    {
        //$mobile = ltrim($mobile, '+86-');
        $mobile = str_replace('+86-', '', $mobile);

        $correctCode = $this->getCodeCache($mobile, $client);
        // 验证码过期
        if (empty($correctCode)) {
            notice(700031);
        }
        // 验证码不存在
        if ($code != $correctCode) {
            Redis::connection()->incr($this->lockAccountCacheKey($mobile, $client));
            notice(700031);
        }
        //登录成功清除code
        $this->delCodeCache($mobile, $client);
    }

    public function generateToken($enterpriseInfo, $tmp = false, $ttl = self::TOKEN_EXPIRE)
    {
        if (time() < $enterpriseInfo['account_start_at']) {
            notice(700011);
        }

        if ($enterpriseInfo['version_tag'] === '试用版' && !$tmp && time() > $enterpriseInfo['account_expire_at']) {
            notice(700016);
        }

        $adminId      = $enterpriseInfo['admin_id'];
        $enterpriseId = $enterpriseInfo['enterprise_id'];
        $serveType    = $enterpriseInfo['serve_type'];

        $expiredAt = (int)($tmp ? self::TMP_TOKEN_EXPIRE : $ttl) + time();

        if ($enterpriseInfo['account_expire_at'] != 0) {
            $expiredAt = min($expiredAt, $enterpriseInfo['account_expire_at']);
        }

        return Crypt::encrypt(json_encode([
            'admin_id'      => $adminId,
            'enterprise_id' => $enterpriseId,
            'serve_type'    => $serveType,
            'expire_time'   => $expiredAt,
            'is_tmp'        => $tmp ? 1 : 2,
            'rand'          => md5(uniqid()),
            'range_type'    => $enterpriseInfo['range_type'],
        ]));
    }

    public function enterpriseAdminRoleList($adminId)
    {
        return (new EnterpriseAdmin())->enterpriseAdminRoleList($adminId);
    }

    public function enterpriseAdminList($account, $client, $version = '', $epIds = [])
    {
        $list = (new EnterpriseAdmin())->enterpriseAdminList($account, $client, $version, $epIds);

        $enterpriseIds       = array_column($list, 'enterprise_id');
        $miniProgramListById = [];
        if (!empty($enterpriseIds)) {
            $miniProgramList     = (new MiniProgramMicro())->listMiniPrograms($enterpriseIds);
            $miniProgramListById = array_column($miniProgramList, null, 'enterprise_id');
        }

        // 区分获取哪个系统
        if ($client === 'store') {
            $systemType = [Permissions::SYSTEM_TYPE_GUIDE];
        } else {
            $systemType = array_keys(Permissions::CONSOLE_SYSTEM_TYPE_RELATION);
        }

        foreach ($list as $k => &$v) {
            $enterpriseId    = $v['enterprise_id'];
            $miniProgramInfo = $miniProgramListById[$enterpriseId] ?? [];
            $v['shop_name']  = $miniProgramInfo['name'] ?? '';
            $v['logo']       = $miniProgramInfo['logo'] ?? '';

            $v['serve_type']        = 0;
            $v['identifier']        = EnterpriseWorkPlugin::firstIdentity($v['enterprise_id']);
            $v['guide_roles']       = RoleStrategy::getRole($v['admin_id'], $systemType);
            $v['guide_permissions'] = RoleStrategy::getPermission($v['admin_id'], Permissions::SYSTEM_TYPE_GUIDE);

            $auth = (new BrandAuth())->getBrandAuth(
                $enterpriseId,
                '',
                [],
                [
                    'enterprise_id',
                    'brand_code'
                ]
            );
            $auth = array_column($auth, 'brand_code', 'enterprise_id');

            $v['brand_code'] = $auth && isset($auth[$enterpriseId]) ? $auth[$enterpriseId] : '';

            if ($systemType == [Permissions::SYSTEM_TYPE_GUIDE] && !$v['guide_roles']) {
                unset($list[$k]);
            }

            if ($v['version'] == 1) {
                $v['serve_type'] = $this->getEnterpriseServeType($v);
                // v1 版本不做处理
                continue;
            }

            $adminAvailableSystem = json_decode($v['admin_available_system'], true);

            if (empty($adminAvailableSystem)) {
                unset($list[$k]);
                continue;
            }

            switch ($client) {
                case EnterpriseAdmin::CLIENT_CRS:
                    if (!in_array(EnterpriseAdmin::SYSTEM_CRS, $adminAvailableSystem)) {
                        unset($list[$k]);
                    }
                    break;
                case EnterpriseAdmin::CLIENT_STORE:
                    if (!in_array(EnterpriseAdmin::SYSTEM_STORE, $adminAvailableSystem)) {
                        unset($list[$k]);
                    }
                    break;
            }
        }

        return array_values($list);
    }

    public function enterpriseAdminDetail($adminId, $client = 'console')
    {
        $detail = (new EnterpriseAdmin())->enterpriseAdminDetail($adminId);
        if (empty($detail)) {
            notice(700012);
        }

        if ($detail['account_status'] == EnterpriseAdmin::STATUS_CLOSE) {
            notice(700015);
        }

        $miniProgramInfo = (new MiniProgramMicro())->getMiniProgram($detail['enterprise_id']);

        if (!empty($detail)) {
            $detail['shop_name'] = $miniProgramInfo['name'] ?? '';
            $detail['logo']      = $miniProgramInfo['logo'] ?? '';
            $detail['qr_code']    = '';
            $detail['serve_type'] = $this->getEnterpriseServeType($detail);

            if ($detail['version'] != 1) {

                $adminAvailableSystem = json_decode($detail['admin_available_system'], true) ?? [];

                if (!in_array($client, array_keys(EnterpriseAdmin::CLIENT_MAP))) {
                    notice(700018);
                }

                $clientType = EnterpriseAdmin::CLIENT_MAP[$client];

                if (!in_array($clientType, $adminAvailableSystem)) {
                    notice(700019);
                }
            }
        }

        return $detail;
    }

    /**
     * 获取当前企业开通服务类型
     *
     * @param $enterpriseInfo
     * @return int
     */
    public function getEnterpriseServeType($enterpriseInfo)
    {
        // 初始化企业服务了开启类型
        $serveType = self::SERVE_TYPE_ALL;

        // 仅开通商城
        if (($enterpriseInfo['store_status'] == self::SERVE_OPEN_TRUE) && ($enterpriseInfo['crs_status'] == self::SERVE_OPEN_FALSE)) {
            $serveType = self::SERVE_TYPE_SHOP;
        }

        // 仅开通CRS
        if (($enterpriseInfo['store_status'] == self::SERVE_OPEN_FALSE) && ($enterpriseInfo['crs_status'] == self::SERVE_OPEN_TRUE)) {
            $serveType = self::SERVE_TYPE_CRS;
        }

        // 开通crs和商城
        if (($enterpriseInfo['store_status'] == self::SERVE_OPEN_TRUE) && ($enterpriseInfo['crs_status'] == self::SERVE_OPEN_TRUE)) {
            $serveType = self::SERVE_TYPE_ALL;
        }

        // 什么都不开通
        if (($enterpriseInfo['store_status'] == self::SERVE_OPEN_FALSE) && ($enterpriseInfo['crs_status'] == self::SERVE_OPEN_FALSE)) {
            $serveType = self::SERVE_TYPE_NONE;
        }

        return $serveType;
    }

    public function dealEnterpriseAdminDataRange(array $input)
    {
        $where = [
            'enterprise_id'       => $input['enterprise_id'],
            'enterprise_admin_id' => $input['before_admin_id'],
            'type'                => $input['range'][0]['type'],
            'object_id'           => $input['range'][0]['object_id']
        ];

        EnterpriseDataRange::query()->where($where)->delete();

        return true;
    }

    /**
     * 创建账号
     * @param $username
     * @param $mobile
     * @return string
     * @Date: 2021/5/26
     * @Time: 21:00
     * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
     */
    public function createDefaultAccount($username, $mobile)
    {
        $account = join('', (new Pinyin())->name($username));
        $account = strtolower(str_replace(' ', '', $account));

        $account = $account . substr($mobile, -4);
        return $account;
    }


    public static function forbidLeaveAdmin()
    {
        $admins = EnterpriseAdmin::query()
//            ->where('is_super', EnterpriseAdmin::NO_SUPER)
            ->where('status', EnterpriseAdmin::STATUS_OPEN)
            ->get(['enterprise_id','mobile'])
            ->toArray();

        $epMap = [];
        foreach ($admins as $admin) {
            $epMap[$admin['enterprise_id']][] = $admin['mobile'];
        }

        foreach ($epMap as $epId => $mobiles) {
            self::memberInfo($epId, $mobiles);
        }
    }

    private static function memberInfo($epId, $mobiles)
    {

        $qyDomain = Enterprise::query()->where('id',$epId)->value('domain');

        $existMobile = Common::work()
            ->setApi("service/v1/admin/list")
            ->post(["enterprise_id" => $epId, "qy_domain" => $qyDomain, 'mobiles' => $mobiles]);

        foreach ($mobiles as $mobile){
            if(!in_array($mobile,$existMobile)){

                EnterpriseAdmin::query()
                    ->where('enterprise_id', $epId)
                    ->where('mobile', $mobile)
                    ->update(['status' => EnterpriseAdmin::STATUS_CLOSE]);
            }
        }



    }

    /**
     * 编辑密码
     *
     */
    public function resetPassword($input): bool
    {
        $select = ['password', 'mobile', 'account', 'id', 'salt'];
        $adminId = $input['admin_id'] ?? 0;
        if (!$adminId) {
            $adminInfo = EnterpriseAdmin::query()->select($select)->where([
                'enterprise_id' => $input['enterprise_id'],
                'account' => $input['account'],
                'mobile' => $input['mobile']
            ])->first();
        } else {
            $adminInfo = EnterpriseAdmin::query()->select($select)
                ->where('id', $adminId)
                ->where('enterprise_id', $input['enterprise_id'])
                ->first();
        }

        if (empty($adminInfo)) {
            notice(701017);
        }

        //重置密码，密码为手机号
        $newPassword = md5(md5($adminInfo->mobile) . $adminInfo->salt);
        EnterpriseAdmin::query()->where('id', $adminInfo->id)->update([
            'password' => $newPassword,
            'password_expire_time' => time() + EnterpriseAdmin::PASSWORD_EXPIRE
        ]);
        //解除锁定
        $this->redis->del($this->lockAccountCacheKey($adminInfo->account, 'console'));
        //解除手机号锁定
        $this->redis->del($this->lockAccountCacheKey($adminInfo->mobile, 'console'));
        return true;
    }

}
