<?php
/**
 * Created by PhpStorm.
 * User: 大熊
 * Date: 2018/11/23
 * Time: 18:20
 */

namespace App\Services;

use App\Models\EnterpriseAdmin;
use App\Http\Resources\EnterpriseAdminResource;

class EnterpriseAdminService extends BaseService
{
    /**
     * 超级管理员密码
     */
    const PASSWORD = 888888;

    /**
     * 企业管理员列表分页查询
     *
     * @param int $limit 每页显示条数
     * @param array $where 查询条件
     *
     * @return mixed
     * @author 梁欢 lianghuan@vchangyi.com
     * @date   2018/11/26 11:16
     *
     */
    public function paginate($limit, $where)
    {
        $query = EnterpriseAdmin::query();

        // 条件组合
        if (!empty($where)) {
            foreach ($where as $key => $value) {
                if (!empty($value) || $value == 0) {
                    switch ($key) {
                        case 'keyword':
                            $query = $query->where(function ($query) use ($value) {
                                $query->orWhere('username', 'like', '%' . $value . '%')
                                    ->orWhere('account', 'like', '%' . $value . '%');
                            });
                            break;
                        case 'role':
                            if (is_array($value)) {
                                $query = $query->whereIn($key, $value);
                            } else {
                                $query = $query->where($key, '!=', $value);
                            }
                            break;
                        case 'id':
                            $query = $query->whereIn($key, $value);
                            break;
                        default:
                            $query = $query->where($key, $value);
                    }
                }
            }
        }

        // 排序
        $query = $query->orderBy('id', 'desc');

        try {
            return EnterpriseAdminResource::collection($query->paginate((int)$limit));
        } catch (\Exception $e) {
            return error(110000);
        }
    }

    /**
     * 根据企业管理员ID获取企业管理员详情
     *
     * @param int $adminId 企业管理员ID
     *
     * @return mixed
     * @author 梁欢 lianghuan@vchangyi.com
     * @date   2018/11/26 11:16
     *
     */
    public function enterpriseAdmin($adminId)
    {
        try {
            return EnterpriseAdmin::query()
                ->where('id', $adminId)
                ->first();
        } catch (\Exception $e) {
            return error(110000);
        }
    }

    /**
     * 根据管理员ID查询管理员详情
     *
     * @param array $where 查询的条件
     *
     * @return mixed
     * @author 梁欢 lianghuan@vchangyi.com
     * @date   2018/11/26 11:16
     *
     */
    public function adminDetail($where)
    {
        try {
            return EnterpriseAdmin::query()
                ->with('enterprise')
                ->where($where)
                ->first();
        } catch (\Exception $e) {
            return error(110000);
        }
    }

    /**
     * 根据手机号查询管理员列表
     *
     * @param $mobile
     * @return mixed
     */
    public function getAdminListByMobile($mobile)
    {
        try {
            return EnterpriseAdmin::query()
                ->with('enterprise')
                ->where('account', '=', $mobile)
                ->get()
                ->toArray();
        } catch (\Exception $e) {
            return error(110000);
        }
    }

    /**
     * 根据企业管理员ID更新登录时间
     *
     * @param int $adminId 企业管理员ID
     *
     * @return mixed
     * @author 梁欢 lianghuan@vchangyi.com
     * @date   2018/11/26 11:16
     *
     */
    public function loginUpdate($adminId)
    {
        try {
            return EnterpriseAdmin::query()
                ->where('id', $adminId)
                ->update(['last_login_time' => time()]);
        } catch (\Exception $e) {
            return error(110000);
        }
    }

    /**
     * 管理员密码修改
     *
     * @param int $adminId 企业管理员ID
     *
     * @return mixed
     * @author 梁欢 lianghuan@vchangyi.com
     * @date   2018/11/29 11:16
     *
     */
    public function passwordUpdate($adminId, $oldPassword, $password)
    {
        $enterpriseAdmin = $this->enterpriseAdmin($adminId);
        if (md5($oldPassword . $enterpriseAdmin->salt) != $enterpriseAdmin->password) {
            return error(202105);
        }

        try {
            return EnterpriseAdmin::query()
                ->where('id', $adminId)
                ->update(['password' => md5($password . $enterpriseAdmin->salt)]);
        } catch (\Exception $e) {
            return error(110000);
        }
    }

    /**
     * 添加企业管理员信息
     *
     * @param array $data 添加的数据
     *
     * @return mixed
     * @author 梁欢 lianghuan@vchangyi.com
     * @date   2018/11/26 11:16
     *
     */
    public function adminAdd($data)
    {
        $isAdmin = EnterpriseAdmin::query()->where([
            'account'       => $data['account'],
            'enterprise_id' => $data['enterprise_id']
        ])->first();
        if ($isAdmin) {
            return error(202003);
        }

        try {
            return EnterpriseAdmin::query()->create($data);
        } catch (\Exception $e) {
            return error(110000);
        }
    }

    /**
     * 编辑企业管理员信息
     *
     * @param int $adminId 修改的企业管理员ID
     * @param array $data 修改的信息
     *
     * @return mixed
     * @author 梁欢 lianghuan@vchangyi.com
     * @date   2018/11/26 11:16
     *
     */
    public function enterpriseAdminSave($adminId, $data)
    {
        if (isset($data['account'])) {
            $enterpriseAdmin = EnterpriseAdmin::query()
                ->where([
                    'account'       => $data['account'],
                    'enterprise_id' => $data['enterprise_id']
                ])
                ->first();

            if (!empty($enterpriseAdmin) && $enterpriseAdmin->id != $adminId) {
                return error(202003);
            }
        }

        try {
            return EnterpriseAdmin::query()->where('id', $adminId)->update($data);
        } catch (\Exception $e) {
            return error(110000);
        }
    }

    /**
     * 企业管理员密码重置
     *
     * @param int $adminId 修改的企业管理员ID
     *
     * @return mixed
     * @author 梁欢 lianghuan@vchangyi.com
     * @date   2018/11/26 11:16
     *
     */
    public function enterpriseAdminPassword($adminId)
    {
        $data['salt']     = $this->makeSalt();
        $data['password'] = md5(md5('888888') . $data['salt']);

        try {
            return EnterpriseAdmin::query()
                ->where('id', $adminId)
                ->update($data);
        } catch (\Exception $e) {
            return error(110000);
        }
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
     * 查询企业超级管理员
     *
     * @param int $adminId 企业管理员ID
     *
     * @return mixed
     * @author zhaoqi zhaoqi@vchangyi.com
     * @date   2019/02/22 15:20
     *
     */
    public function superAdmin($enterpriseId)
    {
        try {
            return EnterpriseAdmin::query()
                ->where([
                    'enterprise_id' => $enterpriseId,
                    'role'          => EnterpriseAdmin::SUPER_ADMIN
                ])->first();
        } catch (\Exception $e) {
            return error(110000);
        }
    }
}
