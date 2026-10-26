<?php

namespace App\Jobs;

use App\Library\Json;
use App\Library\WorkWechatNotice;
use CrsSdk\MicroAuthority\Role;
use CrsSdk\Open\Enterprise as OpenEnterpriseSdk;
use Exception;
use GuzzleHttp\Client;
use Illuminate\Support\Facades\App;
use Illuminate\Support\Facades\Log;

/**
 * 企业注册数据处理
 *
 * @package App\Jobs
 */
class EnterpriseRegisterJob extends Job
{
    /**
     * @var string
     */
    public $queue = 'enterprise-register';

    /**
     * 企业信息
     *
     * @var array
     */
    private $enterprise;

    /**
     * 超管信息
     *
     * @var array
     */
    private $admin;

    /**
     * 明文密码
     *
     * @var string
     */
    private $plaintextPassword;

    /**
     * 选配权限ID列表
     *
     * @var array
     */
    private $permissionList;

    /**
     * @var WorkWechatNotice
     */
    private $workWechatNotice;

    public function __construct(array $enterprise, array $admin, string $plaintextPassword, array $permissionList)
    {
        $this->enterprise = $enterprise;
        $this->admin = $admin;
        $this->plaintextPassword = $plaintextPassword;
        $this->permissionList = $permissionList;
    }

    public function handle()
    {
        Log::debug('企业注册queue开始执行');

        $this->workWechatNotice = new WorkWechatNotice(App::environment());

        $this->initRole($this->enterprise['id'], $this->admin['id'], $this->permissionList);
        // $this->syncOpen($this->enterprise['hash'], $this->enterprise['name']);
        $this->syncScrm($this->enterprise['id'], $this->enterprise['hash'], $this->enterprise['corp_id'], $this->enterprise['domain']);
        //$this->sendEmail();
    }

    private function initRole(int $enterpriseID, int $adminID, array $permissionList = [])
    {
        try {
            (new Role())->init([
                'admin_id'           => $adminID,
                'mark_list'          => ['super_admin'],
                'enterprise_id'      => $enterpriseID,
                'permission_id_list' => $permissionList,
            ]);
        } catch (Exception $exception) {
            Log::error('企业注册处理:初始化角色异常', [
                'message' => $exception->getMessage(),
                'file'    => $exception->getFile(),
                'line'    => $exception->getLine(),
                'args'    => func_get_args()
            ]);

            $this->workWechatNotice->send('企业注册处理:初始化角色异常', [
                'message' => $exception->getMessage(),
                'file'    => $exception->getFile(),
                'line'    => $exception->getLine(),
                'args'    => func_get_args()
            ]);
        }
    }

    private function syncOpen(string $hash, string $name)
    {
        try {
            (new OpenEnterpriseSdk())->add([
                'name'            => $name,
                'enterprise_hash' => $hash
            ]);
        } catch (Exception $exception) {
            Log::error('企业注册处理:同步企业信息给open异常', [
                'message' => $exception->getMessage(),
                'file'    => $exception->getFile(),
                'line'    => $exception->getLine(),
                'args'    => func_get_args()
            ]);

            $this->workWechatNotice->send('企业注册处理:同步企业信息给open异常', [
                'message' => $exception->getMessage(),
                'file'    => $exception->getFile(),
                'line'    => $exception->getLine(),
                'args'    => func_get_args()
            ]);
        }
    }

    private function syncScrm(int $enterpriseID, string $hash, string $corpId, string $domain)
    {
        if (empty($corpId) || empty($domain)) {
            return;
        }
        try {
            $client = new Client();
            $res = $client->request('POST', config('client.api_url.go-micro-scrm') . '/micro-scrm/v1/callback/enterprise-add', [
                'http_errors' => false,
                'verfiy' => false,
                'headers'     => [
                    'Content-Type'    => 'application/json',
                    'Enterprise-Hash' => $hash
                ],
                'body'        => Json::encode([
                    'enterprise_id' => $enterpriseID,
                    'hash'          => $hash,
                    'corp_id'       => $corpId,
                    'domain'        => $domain
                ]),
            ]);
            $result = $res->getBody()->getContents();
            Log::info('企业注册处理:同步scrm响应', [$result]);
        } catch (Exception $exception) {
            Log::error('企业注册处理:创建企业同步scrm异常', [
                'message' => $exception->getMessage(),
                'file'    => $exception->getFile(),
                'line'    => $exception->getLine(),
                'args'    => func_get_args()
            ]);

            $this->workWechatNotice->send('企业注册处理:创建企业同步scrm异常', [
                'message' => $exception->getMessage(),
                'file'    => $exception->getFile(),
                'line'    => $exception->getLine(),
                'args'    => func_get_args()
            ]);
        }
    }

    private function sendEmail()
    {
        $data = $this->enterprise;
        $data['admin'] = $this->admin;
        $data['originPwd'] = $this->plaintextPassword;

        Log::info('企业注册处理:发送注册信息给超管邮箱', $data);

        (new SendEmailJob($data, 0))->handle();
    }
}
