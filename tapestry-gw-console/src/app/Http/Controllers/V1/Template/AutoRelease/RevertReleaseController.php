<?php
/**
 * Created by Phpstorm.
 * User : zoulongbo
 * Date : 2019-08-01
 * Tome : 10:50
 */

namespace App\Http\Controllers\V1\Template\AutoRelease;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\AutoRelease;
use Illuminate\Support\Facades\Log;

class RevertReleaseController extends BaseController
{
    /**
     * 自动化发布 ： 回退
     * @return mixed|\Psr\Http\Message\StreamInterface|void
     */
    public function index()
    {
        try {
            return (new AutoRelease())->versionRevert(['enterprise_id' => $this->enterpriseId]);
        } catch (\Exception $e) {
            Log::info('自动化发布，版本回退失败:' . $e->getMessage());
            error(920304, [$e->getMessage()]);
        }
    }
}
