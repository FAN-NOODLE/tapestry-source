<?php


namespace App\Http\Controllers\App\V1\EnterpriseAdmin;


use App\Http\Controllers\App\V1\BaseController;
use App\Jobs\SendEmailJob;
use App\Models\Enterprise;
use App\Services\EnterpriseNew\EnterpriseAdminNewService;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Mail;

class ResetPasswordController extends BaseController
{
    public function index(Request $request, EnterpriseAdminNewService $enterpriseAdminNewService)
    {
        $enterpriseId = $request->input('enterprise_id', 0);
        // 企业ID不能为空
        if (!$enterpriseId) {
            notice(600001);
        }

        // 根据企业ID查询超管信息
        $admin = $enterpriseAdminNewService->getSuperAdmin($enterpriseId);
        if (empty($admin)) {
            notice(600005);
        }
        $admin = $admin->toArray();

        // 随机生成超管密码-6位
        $password = "REDACTED_FOR_SOURCE_SHARING";
        $md5Password = md5(md5($password) . $admin['salt']);

        // 更新密码
        $res = $enterpriseAdminNewService->updateSuperAdmin(
            ['id' => $admin['id']],
            ['password' => $md5Password]
        );

        if ($res) {
            // 发送邮件
//            $email   = $admin['email'];
//            $title   = '超级管理员密码重置成功';
//            $content = '您的畅移天鲟CRS超级管理员密码重置成功，密码为：' . $password;
//            Mail::raw($content, function ($message) use ($email, $title) {
//                //$message->getHeaders()->addTextHeader('Content-Type', 'text/html');
//                $message->subject($title . '（' . date('Y-m-d') . ')');
//                $message->to($email);
//            });
            $enterprise = Enterprise::query()->where('id', $enterpriseId)->first();
            // 发送邮件
            $enterpriseArr = $enterprise->toArray();
            $enterpriseArr['originPwd'] = $password;
            (new SendEmailJob($enterpriseArr, 0, 'reset_password'))->handle();
            return success(['data' => []]);
        } else {
            notice(600006);
        }
    }
}
