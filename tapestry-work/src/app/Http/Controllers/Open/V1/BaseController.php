<?php
/**
 * Created by PhpStorm.
 * User: 陈朔
 * Date: 2018/11/9
 * Time: 9:54 AM
 */

namespace App\Http\Controllers\Open\V1;

use App\Http\Controllers\Controller;

abstract class BaseController extends Controller
{
    /**
     * 企业信息
     * @var mixed
     */
    protected $enterpriseInfo;

    public function __construct()
    {
        parent::__construct();
        // 接收中间件传递的企业配置信息
        $this->enterpriseInfo = $this->request->get('enterprise_info');
    }

    protected function success(array $data = [], $httpCode = 200)
    {
        return response()->json([
            'code'       => 200,
            'message'    => 'success',
            'request_id' => $_SERVER['x_request_id'] ?? LOG_ID,
            'data'       => !empty($data) ? $data : (object)$data
        ], $httpCode);
    }
}
