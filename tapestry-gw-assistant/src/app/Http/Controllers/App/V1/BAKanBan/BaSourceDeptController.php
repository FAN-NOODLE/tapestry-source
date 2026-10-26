<?php

namespace App\Http\Controllers\App\V1\BAKanBan;

use App\Http\Controllers\App\V1\BaseController as Controller;
use CrsSdk\GoMicroScrm\BaKanBan;
use Illuminate\Support\Facades\Auth;

/**
 * 小程序获取BA业绩看板的数据
 */
class BaSourceDeptController extends Controller
{
    public function index()
    {
        $param = request()->post();
        $param['enterprise_id'] = Auth::user()['enterprise_id'];

        $result = (new BaKanBan())->baSourceDepte($param);
        return success($result);
    }
}