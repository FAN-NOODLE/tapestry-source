<?php

namespace App\Http\Controllers\Service\V1\Application;

use App\Http\Controllers\Service\V1\BaseController;
use Illuminate\Support\Facades\Redis;
use Illuminate\Http\Request;

class CloseController extends BaseController
{

	/**
     * BOSS关闭应用后通知
     *
     * @param Request $request
     * @return mixed
     */
	public function index(Request $request)
	{
		$enterpriseId = $request->post('enterprise_id');

		// 存redis在BaseController判断是否存在，存在直接返回，登陆接口删除缓存
		$key = 'enterprise:application:close:'.$enterpriseId;
		Redis::hmset($key, ['enterprise_id' => $enterpriseId, 'time' => time()]);

		return success(['data' => []]);
	}

}
