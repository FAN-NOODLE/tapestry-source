<?php
/**
 * Created by Phpstorm.
 * User : zoulongbo
 * Date : 2019-08-01
 * Tome : 10:52
 */

namespace App\Http\Controllers\V1\Template\AutoRelease;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\AutoRelease;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class ReleaseController extends BaseController
{
    /**
     * 发布类型： 立即发布
     */
    const RELEASE_TYPE_NOW = 1;

    /**
     * 发布类型： 定时发布
     */
    const RELEASE_TYPE_TIME = 2;

    /**
     * 自动化发布： 发布
     *
     * @param Request $request
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function index(Request $request)
    {
        $input = $request->only(['release_type', 'pre_release_time', 'version_id']);
        if (empty($input['version_id'])) {
            error(920308);
        }
        if (empty($input) || !in_array($input['release_type'], [self::RELEASE_TYPE_NOW, self::RELEASE_TYPE_TIME])) {
            error(920305);
        }
        if (($input['release_type'] == self::RELEASE_TYPE_TIME) && empty($input['pre_release_time'])) {
            error(920306);
        }

        try {
           // todo 直接shop发布
            $input['enterprise_id'] = $this->enterpriseId;
            (new AutoRelease())->release($input);
            return success([]);
        } catch (\Exception $e) {
            Log::info('发布失败，发布类型:'.$input['release_type'].' reason:' . $e->getMessage());
            error(920307, [$e->getMessage()]);
        }
    }
}
