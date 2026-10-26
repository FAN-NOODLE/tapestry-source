<?php
/**
 * Created by Phpstorm.
 * User : zoulongbo
 * Date : 2019-08-01
 * Tome : 10:38
 */
namespace App\Http\Controllers\V1\Template\AutoRelease;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Boss\MiniProgram;
use CrsSdk\Shop\AutoRelease;
use Illuminate\Support\Facades\Log;

class SubmitAuditController extends BaseController
{
    /**
     * 自动化打包 ： 提审
     * @return mixed|\Psr\Http\Message\StreamInterface|void
     */
    public function index()
    {
        $input['enterprise_id'] = $this->enterpriseId;

        $this->writeOperatorLog("提交审核");

        try {
            $templateData = (new AutoRelease())->submitAudit($input);
            $result       = (new MiniProgram())->autoPack($templateData['data']);
        } catch (\Exception $e) {
            Log::info('自动化打包提审失败：' . $e->getMessage());
            error(920301, [$e->getMessage()]);
        }

        return $result;
    }
}
