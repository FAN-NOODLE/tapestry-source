<?php
/**
 * 自动化发布
 * 定时发布-撤销发布
 * Created by Phpstorm.
 * User : zoulongbo
 * Date : 2019-08-05
 * Tome : 10:40
 */

namespace App\Http\Controllers\V1\Template\AutoRelease;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\AutoRelease;
use Illuminate\Http\Request;

class CancelReleaseController extends BaseController
{
    public function index(Request $request)
    {
        $input = $request->only(['version_id']);
        if (empty($input)) {
            error(920208);
        }

        return (new AutoRelease())->cancelRelease($input);
    }
}
