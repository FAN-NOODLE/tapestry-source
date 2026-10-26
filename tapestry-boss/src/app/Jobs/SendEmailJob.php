<?php

namespace App\Jobs;

use App\Library\Sendcloud\SendCloud;
use App\Models\EnterpriseAdmin;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Str;

class SendEmailJob extends Job
{
    protected $sendMail = [
        "apiUser" => "crs_notice",
        "apiKey" => "REDACTED_FOR_SOURCE_SHARING",
        "from" => "changxun@chufa.vchangyi.com",
        "fromName" => "畅鲟",
        "xsmtpapi" => "",
        "templateInvokeName" => "",

    ];
    protected $versionDesc = [
        2 => "零售版",
        3 => "电商版",
        4 => "连锁门店版"
    ];
    protected $data = [];
    protected $smartGuideUserCount = 0;

    // 邮件场景
    protected $scene = '';

    public function __construct($data, $smartGuideUserCount, $scene = '')
    {
        $this->data = $data;
        if (is_object($data)) {
            $this->data = $data->toArray();
        }
        $this->smartGuideUserCount = $smartGuideUserCount;

        $this->scene = $scene;

    }

    /**
     * 发送邮件
     */
    public function handle()
    {
        // 查询企业超管
        $superAdmin = EnterpriseAdmin::query()
            ->where('enterprise_id', $this->data['id'])
            ->where('is_super', EnterpriseAdmin::IS_SUPER)
            ->get();
        if (empty($superAdmin)) {
            Log::info("未查询到企业超管信息", $this->data);
            return;
        }

        switch ($this->scene) {
            case 'reset_password':
                $xsmtpapi = $this->getResPasswordEmailData($superAdmin);
                break;
            default:
                $xsmtpapi = $this->getDefaultEmailData($superAdmin);
                break;
        }
        Log::info('xsmtpapi发送数据：'. json_encode($xsmtpapi, true));

        try {
            $this->sendMail['xsmtpapi'] = json_encode($xsmtpapi, JSON_UNESCAPED_UNICODE);
            $result = (new SendCloud())->sendTemplateEmail($this->sendMail);
        } catch (\Exception $e) {
            Log::error("发送邮件错误，错误信息：{$e->getMessage()}, 错误码：{$e->getCode()}");
        }

    }
    protected function getResPasswordEmailData($superAdmin)
    {
        if ($this->data["version_tag"] == "正式版") {
            $smart = "";
            if ($this->smartGuideUserCount > 0) {
                $smart = "，同时开通{$this->smartGuideUserCount}个智慧导购账号";
            }
            $this->sendMail["templateInvokeName"] = "crs_account_reset_password_formal";
            foreach ($superAdmin as $key => $admin) {
                $xsmtpapi["to"][$key] = $admin["email"];
                $xsmtpapi["sub"]["%name%"][$key] = $this->data['name'];
                $xsmtpapi["sub"]["%start%"][$key] = date("Y年m月d日", $this->data['account_start_at']);
                $xsmtpapi["sub"]["%end%"][$key] = date("Y年m月d日", $this->data['account_expire_at']);
                $xsmtpapi["sub"]["%email%"][$key] = $admin["email"];
                $xsmtpapi["sub"]["%account_total%"][$key] = $this->data['account_limit'];
                $xsmtpapi["sub"]["%pwd%"][$key] = $this->data["originPwd"] ?? "";
                $xsmtpapi["sub"]["%smart%"][$key] = $smart;
            }
        } else if ($this->data["version_tag"] == "试用版") {
            $this->sendMail["templateInvokeName"] = "crs_account_reset_password_experience";
            foreach ($superAdmin as $key => $admin) {
                $xsmtpapi["to"][$key] = $admin["email"];
                $xsmtpapi["sub"]["%name%"][$key] = $this->data['name'];
                $xsmtpapi["sub"]["%start%"][$key] = date("Y年m月d日", $this->data['account_start_at']);
                $xsmtpapi["sub"]["%end%"][$key] = date("Y年m月d日", $this->data['account_expire_at']);
                $xsmtpapi["sub"]["%email%"][$key] = $admin["email"];
                $xsmtpapi["sub"]["%account_total%"][$key] = $this->data["account_limit"];
                $xsmtpapi["sub"]["%pwd%"][$key] = $this->data["originPwd"] ?? "";
                $xsmtpapi["sub"]["%smart%"][$key] = "";
            }

        }
        return $xsmtpapi;
    }
    protected function getDefaultEmailData($superAdmin)
    {
        if ($this->data["version_tag"] == "正式版") {
            $smart = "";
            if ($this->smartGuideUserCount > 0) {
                $smart = "，同时开通{$this->smartGuideUserCount}个智慧导购账号";
            }
            $this->sendMail["templateInvokeName"] = "crs_account_renewal";
            foreach ($superAdmin as $key => $admin) {
                $xsmtpapi["to"][$key] = $admin["email"];
                $xsmtpapi["sub"]["%version%"][$key] = $this->versionDesc[$this->data["version"]];
                $xsmtpapi["sub"]["%name%"][$key] = $this->data['name'];
                $xsmtpapi["sub"]["%start%"][$key] = date("Y年m月d日", $this->data['account_start_at']);
                $xsmtpapi["sub"]["%end%"][$key] = date("Y年m月d日", $this->data['account_expire_at']);
                $xsmtpapi["sub"]["%email%"][$key] = $admin["email"];
                $xsmtpapi["sub"]["%mobile%"][$key] = $admin["mobile"];
                $xsmtpapi["sub"]["%account_total%"][$key] = $this->data['account_limit'];
                $xsmtpapi["sub"]["%pwd%"][$key] = $this->data["originPwd"] ?? "";
                $xsmtpapi["sub"]["%smart%"][$key] = $smart;
                $xsmtpapi["sub"]["%remark%"][$key] = isset($this->data['emailRemark']) ? $this->data['emailRemark'] : "(初始密码)";
            }
        } else if ($this->data["version_tag"] == "试用版") {
            $this->sendMail["templateInvokeName"] = "crs_register_enterprise_notice";
            foreach ($superAdmin as $key => $admin) {
                $xsmtpapi["to"][$key] = $admin["email"];
                $xsmtpapi["sub"]["%version%"][$key] = $this->versionDesc[$this->data["version"]];
                $xsmtpapi["sub"]["%name%"][$key] = $this->data['name'];
                $xsmtpapi["sub"]["%start%"][$key] = date("Y年m月d日", $this->data['account_start_at']);
                $xsmtpapi["sub"]["%end%"][$key] = date("Y年m月d日", $this->data['account_expire_at']);
                $xsmtpapi["sub"]["%email%"][$key] = $admin["email"];
                $xsmtpapi["sub"]["%mobile%"][$key] = $admin["mobile"];
                $xsmtpapi["sub"]["%account_total%"][$key] = $this->data["account_limit"];
                $xsmtpapi["sub"]["%pwd%"][$key] = $this->data["originPwd"] ?? "";
                $xsmtpapi["sub"]["%smart%"][$key] = "";
            }

        }
        return $xsmtpapi;
    }
}