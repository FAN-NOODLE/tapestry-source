<?php
/**
 * Created by Phpstorm.
 * User : zoulongbo
 * Date : 2019-08-01
 * Tome : 10:38
 */
namespace App\Http\Controllers\V1\Template\AutoRelease;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\AutoRelease;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class CancelAuditController extends BaseController
{
    /**
     * 自动化打包： 撤审
     * @param Request $request
     *
     * @return mixed|\Psr\Http\Message\StreamInterface|void
     */
    public function index(Request $request)
    {
        $versionId = $request->input('version_id');
        if (empty($versionId)) {
            error(920208);
        }

        try {
            return (new AutoRelease())->cancelAudit([
                'enterprise_id' => $this->enterpriseId,
                'version_id' => $versionId
            ]);
        } catch (\Exception $e) {
            Log::info('自动化发布，撤审失败:' . $e->getMessage());
            error(920302, [$e->getMessage()]);
        }
    }
}
