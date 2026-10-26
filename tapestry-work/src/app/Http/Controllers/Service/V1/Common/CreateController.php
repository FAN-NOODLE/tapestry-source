<?php
/**
 * Created by :PhpStorm
 * User: lvshuaishuai
 * Date: 2020/11/18
 * Time: 4:00 下午
 */

namespace App\Http\Controllers\Service\V1\Common;


use App\Http\Controllers\Service\V1\BaseController;
use App\Library\EasyWechatCacheRedis\MiniProgram;
use App\Services\EnterpriseService;
use App\Services\CommonService;
use EasyWeChat\Kernel\Http\StreamResponse;
use Exception;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class CreateController extends BaseController
{
    public function index(Request $request)
    {
        $this->validateData($request);

        /**
         * @param string    path            小程序页面路径
         * @param int       enterprise_id   企业ID
         * @param array     params          太阳码携带的参数
         */
        $data         = $request->only('path', 'enterprise_id', 'params');
        $path         = $data['path'];
        $enterpriseId = $data['enterprise_id'];
        $params       = $data['params'];
        try {
//            $workSunCode = WorkSunCode::query()
//                ->where('enterprise_id', '=', $enterpriseId)
//                ->where('worker_id', '=', $params['worker_id'])
//                ->first();
//            if (!empty($workSunCode)) {
//                $url = $workSunCode->sun_code;
//            } else {
//                $url = $this->getQrcode($path, $params, $enterpriseId);
//            }
            $url = $this->getQrcode($path, $params, $enterpriseId);
            return success([
                'data' => ['code_url' => $url],
                'meta' => []
            ]);
        } catch (Exception $exception) {
            Log::info($exception);
        }
        return [];
    }

    public function validateData($request)
    {
        $rules    = [
            'path'          => 'required',
            'enterprise_id' => 'required|integer',
        ];
        $messages = [
            'path.required'          => 'Sorry,二维码路径不能为空',
            'enterprise_id.required' => 'Sorry,企业id不能为空',
            'enterprise_id.integer'  => 'Sorry,企业id必须为整数',
        ];
        $this->validate($request, $rules, $messages);
    }

    /**
     * 生成二维码
     *
     * @param $path
     * @param $scene
     * @param $enterpriseId
     *
     * @return string
     * @throws Exception
     */
    public function getQrcode($path, $scene, $enterpriseId, $handle = true)
    {
        // 确认scene 数据类型，如果是数组，则接收的是二维码参数，需要生成参数ID
        if ($handle) {
            if (is_array($scene)) {
                $scene = encryptShortLink($scene);
            }
            if (empty($scene)) {
                throw new Exception('scene server error');
            }
        }

        return $this->qrcodeUnlimit($scene, $path, $enterpriseId, $handle);
    }

    /**
     * 场景二维码
     *
     * @param $scene
     * @param $path
     * @param $enterpriseId
     *
     * @return string
     */
    public function qrcodeUnlimit($scene, $path, $enterpriseId, $handle = true)
    {
        if ($handle) {
            $workerId = decryptShortLink($scene);
            // 获取小程序id和小程序secret
            $appId     = $workerId['args'][4]['value'];
            $appSecret = $workerId['args'][5]['value'];
        } else {
            $pluginList = EnterpriseService::getEnterpriseIdentifierList($enterpriseId);
            $appId      = '';
            $appSecret  = '';
            foreach ($pluginList as $item) {
                if ($item['identifier'] == "crs_store_assistant_self") {
                    $appId     = $item['mini_app_id'];
                    $appSecret = $item['mini_app_secret'];
                    break;
                }
            }
        }

        if (empty($appId) || empty($appSecret)) {
            return '';
        }
        $config = [
            'app_id'        => $appId,
            'secret'        => $appSecret,
            'response_type' => 'array'
        ];
        try {
            $app      = MiniProgram::init($config);
            $response = $app->app_code->getUnlimit($scene, [
                'page'       => $path,
                'width'      => 100,
                'line_color' => [
                    'r' => 105,
                    'g' => 166,
                    'b' => 134,
                ],
            ]);
            if ($response instanceof StreamResponse) {
                $content = $response->getBodyContents();
                $url     = CommonService::uploadContent($content);
//                $url = (new AzureBlobService())->uploadByContent($content);
//                $createDpData[] = [
//                    'enterprise_id' => $enterpriseId,
//                    'worker_id'     => $workerId['args'][3]['value'],
//                    'sun_code'      => $url,
//                    'created_at'    => time()
//                ];
//                DB::table('work_sun_code')->insert($createDpData);
                return $url ?? '';
            } else if (is_array($response)) {
                Log::error('生成太阳码错误信息 error', $response);
            }
        } catch (Exception $e) {
            writeLog('生成二维码失败', $e);
        }

        return '';
    }


    /**
     * @param Request $request
     *
     * @return array|void
     * @Time: 2022/11/28
     */
    public function createQrcode(Request $request)
    {
        $this->validateData($request);

        $data         = $request->only('path', 'enterprise_id', 'params');
        $path         = $data['path'];
        $enterpriseId = $data['enterprise_id'];
        $params       = $data['params'];
        try {
            $url = $this->getQrcode($path, $params, $enterpriseId, false);
            return success(['code_url' => $url]);
        } catch (Exception $exception) {
            Log::info($exception);
        }
    }


    /**
     * @param Request $request
     *
     * @return array|string
     * @throws \GuzzleHttp\Exception\GuzzleException
     * @Time: 2022/11/28
     */
    public function getQRCodeAndLink(Request $request)
    {
        $this->validateData($request);
        $data         = $request->only('path', 'enterprise_id', 'scene', 'width','env_version');
        $path         = $data['path'];
        $enterpriseId = $data['enterprise_id'];
        $scene       = $data['scene'];
        $envVersion       = !empty($data['env_version']) ? $data['env_version'] : 'release';
        $sunWidth       = $data['width'] ?? 1417;

        $pluginList = EnterpriseService::getEnterpriseIdentifierList($enterpriseId);
        $appId = '';
        $appSecret = '';
        foreach ($pluginList as $item) {
            if ($item['identifier'] == "crs_store_assistant_self") {
                $appId = $item['mini_app_id'];
                $appSecret = $item['mini_app_secret'];
                break;
            }
        }

        if (empty($appId) || empty($appSecret)) {
            return '';
        }
        $config = [
            'app_id' => $appId,
            'secret' => $appSecret,
            'response_type' => 'array'
        ];
        $app = MiniProgram::init($config);
        $url = '';
        $link = '';
        try {
            $response = $app->app_code->getUnlimit($scene, [
                'page' => ltrim($path,'/'),
                'width' => $sunWidth,
                'env_version' => $envVersion
            ]);
            if ($response instanceof StreamResponse) {
                $content = $response->getBodyContents();
                $url = CommonService::uploadContent($content);
            }
            $response = $app->app_code->httpPostJson('/wxa/generatescheme', [
                "jump_wxa" => [
                    'path'=>$path,
                    'query'=>'scene='.$scene,
                    'env_version' => $envVersion
                ],
                'expire_type' => 1,
                'expire_interval' => 30
            ]);
            $link = $response['openlink'];
        } catch (Exception $exception) {
            Log::info($exception);
        }
        return success(['code_url' => $url,'link'=>$link]);

    }

}
