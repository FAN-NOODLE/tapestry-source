<?php
/**
 * Created by PhpStorm.
 * User: 陈朔
 * Date: 2018/11/9
 * Time: 9:54 AM
 */

namespace App\Http\Controllers\Service\V1\MiniProgram;

use App\Http\Controllers\Service\V1\BaseController;
use App\Services\OpenPlatform\MiniProgramService;
use App\Models\EnterpriseThirdPartAccount;
use Illuminate\Http\Request;

abstract class MiniBaseController extends BaseController
{
    /***
     * @var MiniProgramService
     */
    protected $miniProgramService;

    /**
     * @var Request
     */
    protected $request;

    /**
     * @var array|mixed|string|null
     */
    protected $appId;

    public function __construct(
        MiniProgramService $miniProgramService,
        Request $request
    ) {
        parent::__construct();
        $this->request = $request;
        $this->miniProgramService = $miniProgramService;
        $this->appId = $request->input('app_id') ?? $this->getEnterpriseAppId($request->input(['enterprise_id']));
    }

    /**
     * 企业id查询授权小程序appId 默认商城
     * @param     $enterpriseId
     * @param int $type
     *
     * @return mixed|string
     */
    public function getEnterpriseAppId($enterpriseId, $type = 1)
    {
        $auth = EnterpriseThirdPartAccount::query()
                ->where('enterprise_id', $enterpriseId)
                ->where('type', $type)
                ->where('status', EnterpriseThirdPartAccount::STATUS_ACTIVE)
                ->first(['app_id']);
        if (empty($auth)) {
            return '';
        }
        return $auth->app_id;
    }
}
