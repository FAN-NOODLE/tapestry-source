<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/12/10
 * Time: 16:24
 */

namespace App\Http\Controllers\Service\V1\WelcomeSpeech;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\WelcomeSpeech;
use Illuminate\Http\Request;

class DeleteController extends BaseController
{
    /**
     * 欢迎语删除
     *
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request)
    {
        $id = $request->get('id', 0);
        WelcomeSpeech::query()
            ->where('id', $id)
            ->delete();

        return success([]);
    }
}
