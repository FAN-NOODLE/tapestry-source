<?php
/**
 * Created by PhpStorm.
 * User: 陈朔
 * Date: 2018/11/9
 * Time: 9:54 AM
 */

namespace App\Http\Controllers\App\V1\MiniProgram;

use App\Http\Controllers\App\V1\BaseController;
use App\Services\OpenPlatform\MiniProgramService;
use App\Models\EnterpriseThirdPartAccount;

abstract class MiniBaseController extends BaseController
{
    /***
     * @var MiniProgramService
     */
    protected $miniProgramService;

    /**
     * @var EnterpriseThirdPartAccount
     */
    protected $thirdPartAccount;

    public function __construct(MiniProgramService $miniProgramService, EnterpriseThirdPartAccount $thirdPartAccount)
    {
        parent::__construct();
        
        $this->miniProgramService = $miniProgramService;
        $this->thirdPartAccount = $thirdPartAccount;
    }
}
