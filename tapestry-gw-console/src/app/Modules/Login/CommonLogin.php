<?php

namespace App\Modules\Login;

use CrsSdk\Boss\EnterpriseAdmin;
use CrsSdk\Work\Admin as WorkAdmin;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Redis;
use CrsSdk\Boss\Enterprise as BossEnterprise;
use CrsSdk\Work\Enterprise as WorkEnterprise;

class CommonLogin
{
    /**
     * 授权是否有效
     *
     * @var bool
     */
    public $authEffective = true;

    /**
     * 企业是否试用到期
     *
     * @var bool
     */
    public $isExpire = false;

    /**
     * 企业信息
     *
     * @var null|array
     */
    public $enterprise;

    /**
     * UC企业信息
     *
     * @var null|array
     */
    public $ucEnterprise;

    /**
     * 管理员信息
     *
     * @var null|array
     */
    public $admin;

    /**
     * 企业微信ID
     *
     * @var null|string
     */
    protected $corpId;

    /**
     * UC员工ID
     *
     * @var null|string
     */
    public $memberId;

    /**
     * @var \Illuminate\Redis\Connections\Connection
     */
    protected $redis;

    /**
     * @var WorkEnterprise
     */
    protected $workEnterpriseSdk;

    /**
     * @var WorkAdmin
     */
    protected $workAdminSdk;

    /**
     * @var BossEnterprise
     */
    protected $bossEnterpriseSdk;

    public function __construct()
    {
        $this->redis             = Redis::connection();
        $this->workEnterpriseSdk = new WorkEnterprise();
        $this->workAdminSdk      = new WorkAdmin();
        $this->bossEnterpriseSdk = new BossEnterprise();
    }

    /**
     * 注册企业
     *
     * @param array $data
     * @throws \Exception
     */
    public function registerEnterprise($data)
    {
        // 版本 1:标准版（旧版） ,2:中心商城版(无门店), 3:中心商城(有门店), 4:直营连锁门店
        $data['version'] = (int)$data['version'];
        if (!in_array($data['version'], [2, 3, 4])) {
            throw new \Exception('版本不存在');
        }
        // 进入注册企业流程
        $params = [
            'name'              => $data['name'],
            'contact'           => $data['contact'],
            'mobile'            => $data['mobile'],
            'address'           => $data['name'],
            'corp_id'           => $this->corpId,
            'account'           => $data['account'],
            'email'             => $data['account'],
            'version'           => $data['version'],
            'version_tag'       => '试用版',
            'business_source'   => 1,// 业务来源:企业微信
            'account_expire_at' => time() + 86400 * 7,
            'account_start_at'  => time(),
            'account_limit'     => 5,
            'industry_id'       => $data['industry_id'],
            'member_id'         => $this->memberId,
            'available_system'  => [1],
            'permission_list'   => [],
        ];

        try {
            $enterprise = $this->bossEnterpriseSdk->add($params)['data'];
        } catch (\Exception $exception) {
            Log::error('注册企业异常', [
                'message' => $exception->getMessage(),
                'file'    => $exception->getFile(),
                'line'    => $exception->getLine(),
                'data'    => $data
            ]);
            throw new \Exception('注册企业异常: ' . $exception->getMessage());
        }

        $this->enterprise = $this->bossEnterpriseSdk->detail([
            'enterprise_id' => $enterprise['id']
        ])['data'];

        // 重新加载一下管理员信息
        $this->admin = $this->adminInfo($this->memberId, $enterprise['id']);
    }

    /**
     * 绑定管理员
     *
     * @param string $mobile 手机号
     * @throws \Exception
     */
    public function bindAccount($mobile)
    {
        $enterpriseAdminSdk = new EnterpriseAdmin();

        // 根据当前企业+手机号检索是否存在管理员,存在则可以绑定此member_id
        try {
            $this->admin = $enterpriseAdminSdk->detail([
                'enterprise_id' => $this->enterprise['id'],
                'mobile'        => $mobile
            ])['data'];
        } catch (\Exception $exception) {
            // 管理员不存在 需要管理后台先创建
            if ($exception->getCode() !== 202104) {
                throw $exception;
            }
        }
        $this->admin['admin_id']    = $this->admin['id'];
        $this->admin['operator_id'] = $this->admin['id'];
        $this->admin['role_id']     = 1;

        // 绑定member_id
        $this->admin['member_id'] = $this->memberId;
        $enterpriseAdminSdk->save($this->admin);
    }

    /**
     * 从UC查询员工ID
     *
     * @param string $epId
     * @param string $eaId
     * @return mixed
     * @throws \Exception
     */
    protected function queryMemberId($epId, $eaId)
    {
        $adminInfo = $this->workAdminSdk->detail([
            'qy_domain' => $epId,
            'admin_id'  => $eaId
        ]);

        Log::info('从UC查到的管理员信息', $adminInfo);

        return $adminInfo['memUid'] ?? '';
    }

    /**
     * 查询boss管理员信息
     *
     * @param        $memberId
     * @param string $enterpriseId
     *
     * @return |null
     * @throws \Exception
     */
    protected function adminInfo($memberId, $enterpriseId = '')
    {
        $data               = null;
        $enterpriseAdminSdk = new EnterpriseAdmin();
        if(!empty($enterpriseId)){
            // 检查boss是否存在此管理员
            try {
                $data = $enterpriseAdminSdk->detail([
                    'member_id'     => $memberId,
                    'enterprise_id' => $enterpriseId,
                ])['data'];
            } catch (\Exception $exception) {
                // 管理员不存在 需要先绑定
                if ($exception->getCode() !== 202007) {
                    throw $exception;
                }
            }

            return $data;
        }
    }

    /**
     * 查询boss企业信息
     *
     * @param string $corpId 企业微信ID
     * @return array|null
     * @throws \Exception
     */
    protected function enterpriseInfo($corpId)
    {
        $data = null;
        try {
            $data = $this->bossEnterpriseSdk->detail([
                'corp_id' => $corpId
            ])['data'];
        } catch (\Exception $exception) {
            // 此code码说明企业不存在 可以创建企业
            if ($exception->getCode() !== 201002) {
                throw $exception;
            }
        }
        return $data;
    }

    /**
     * 检查企业是否过期 true:过期 false:没过期
     *
     * @return bool
     */
    protected function checkExpire()
    {
        if ($this->enterprise === null) {
            return false;
        }
        if ($this->enterprise['version_tag'] === '试用版' && ($this->enterprise['account_start_at'] >= time() || $this->enterprise['account_expire_at'] <= time())) {
            return true;
        }

        return false;
    }

    /**
     * 生成暂存类token
     *
     * @return string
     * @throws \ReflectionException
     */
    public function __toString()
    {
        $args   = [];
        $refObj = new \ReflectionObject($this);
        foreach ($refObj->getProperties() as $key => $value) {
            $property = $value->getName();

            // 私有属性跳过
            if ($value->isPrivate()) {
                continue;
            }

            // 无须的属性跳过
            if (in_array($property, ['redis'])) {
                continue;
            }

            // 属性是Sdk结尾跳过
            if (substr($property, strlen($property) - 3) === 'Sdk') {
                continue;
            }
            $args[$property] = $this->{$property};
        }

        $token = md5($this->corpId);
        $this->redis->setex("gw-console:work-wechat-token:{$token}", 600, json_encode($args, JSON_UNESCAPED_UNICODE));

        return $token;
    }

    /**
     * 还原类
     *
     * @param string $token
     * @throws \Exception
     */
    public function reduction($token)
    {
        $data = $this->redis->get("gw-console:work-wechat-token:{$token}");
        if ($data === null) {
            throw new \Exception('缓存token数据不存在');
        }
        $args = json_decode($data, true);
        foreach ($args as $key => $value) {
            $this->{$key} = $value;
        }
    }
}
