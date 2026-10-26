<?php
/**
 * 小程序 手动提审
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/1/10
 * Time: 15:36
 */

namespace App\Http\Controllers\App\V1\MiniProgram;

use Illuminate\Http\Request;

class TrigSubmitController extends MiniBaseController
{
    /**
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request)
    {
        $appId = $request->input('app_id');

        if (!$appId) {
            return error(600001);
        }
        $thirdPartAccount = $this->thirdPartAccount->where('app_id', $appId)->where('status', 1)->first();
        if (!$thirdPartAccount) {
            return error(201002);
        }

        //开始提交审核,手动提审 ，默认为true
        try {
            $result = $this->miniProgramService->submitVersion($thirdPartAccount, true);
        } catch (\Exception $e) {
            app('log')->error($e);
            $result = false;
        }

        return [
            'data' => [
                'result' => $result,
            ],
        ];
    }
}
