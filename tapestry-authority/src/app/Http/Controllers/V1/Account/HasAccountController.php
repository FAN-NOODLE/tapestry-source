<?php


namespace App\Http\Controllers\V1\Account;

use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use App\Http\Requests\AccountRequest;
use App\Models\EnterpriseAdmin;
use App\Models\EnterpriseDataRange;
use App\Services\EnterpriseAdminRoleService;
use App\Services\EnterpriseAdminService;
use App\Services\EnterpriseDataRangeService;
use App\Services\EnterpriseService;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;
use Illuminate\Validation\ValidationException;
use function Composer\Autoload\includeFile;

class HasAccountController extends BaseController
{
    /**
     * 创建门店时同步创建账号
     *
     * @return JsonResponse
     * @throws ValidationException
     * @throws Exception
     */
    public function index()
    {
        $input = $this->request->only([
            'operator_id',
            'enterprise_id',
            'username',
            'mobile',
            'before_mobile',
            'range_type',
            'range',
            'admin_available_system',
            'role_ids',
        ]);
        Log::info("zhimma-params" , $input);

        // 验证
        AccountRequest::add($input);
        // 验证手机号
        $mobileArr = explode('-', $input['mobile']);
        $mobile    = $mobileArr[1] ?? $input['mobile'];
        $preg      = "/^1[3-9]\d{9}$/";
        if (!preg_match($preg, $mobile)) {
            notice(701001); // 请输入正确的手机号码
        }

        // 验证创建账号规则
        $detail = (new EnterpriseService())->enterpriseDetail($input['enterprise_id']);
        if (empty($detail)) {
            notice(701005); // 企业不存在
        }

        $accountCount = (new EnterpriseAdminService())->getEnterpriseChildAdminCount($input['enterprise_id']);
        if ($detail['account_limit'] == $accountCount) {
            notice(701004, [$accountCount]); // 已达到允许创建子账号上限
        }

        if (count($mobileArr) > 1) {
            $input['global_roaming'] = $mobileArr[0];
        }
        $input['account'] = $input['mobile'];
        $input['mobile']  = $mobile;

        // 判断此账号是否已经存在
        $admin = (new EnterpriseAdminService())->getEnterpriseAdminByMobile([
            'mobile' => $input['mobile'],
            'enterprise_id' => $input['enterprise_id']
        ]);

        // 获取被替换账号详情
        if (!empty($input['before_mobile'])) {
            $beforeAdmin = (new EnterpriseAdminService())->getEnterpriseAdminByMobile([
                'mobile'        => $input['before_mobile'],
                'enterprise_id' => $input['enterprise_id']
            ]);
            Log::info('zhimma' , [$beforeAdmin]);

            if (!empty($beforeAdmin)) {
                // 处理替换逻辑
                $input['before_admin_id'] = $beforeAdmin['id'];
                (new EnterpriseAdminService())->dealEnterpriseAdminDataRange($input);
            }
        }

        if (empty($admin)) {
            (new AddController())->addAccountData($input);
        } else {
            // 编辑信息
            $input['id'] = $admin['id'];
            $this->editAccountData($input);
        }

        return success();
    }

    /**
     * @param $input
     * @throws Exception
     */
    public function editAccountData($input)
    {
        // 更新数据
        try {
            DB::beginTransaction();

            // 编辑范围信息
            (new EnterpriseDataRangeService())->editDataRange($input, false);

            // 编辑角色表信息
            (new EnterpriseAdminRoleService())->editEnterpriseAdminRole($input, false);

            DB::commit();
        } catch (Exception $exception) {
            DB::rollBack();
            Log::info('同步创建账号失败' . $exception->getMessage(), $input);
            throw new Error($exception->getCode(), $exception->getMessage());
        }
    }

}