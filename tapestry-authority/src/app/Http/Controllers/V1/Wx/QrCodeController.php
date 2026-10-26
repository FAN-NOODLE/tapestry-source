<?php


namespace App\Http\Controllers\V1\Wx;


use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use App\Traits\OpenPlatformTrait;
use CrsSdk\Common\QCloudCos;
use EasyWeChat\Kernel\Http\StreamResponse;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class QrCodeController extends BaseController
{
    use OpenPlatformTrait;

    public function index(Request $request)
    {
        $params = $this->getParams($request);

        $enterpriseId = $params['enterprise_id'];
        $path         = $params['path'];
        $qrParams     = $params['params'] ?? [];

        $url = $this->qrCodeUnLimit($enterpriseId, $path, $qrParams);

        return success(['url' => $url]);
    }

    public function qrCodeUnLimit($enterpriseId, $path, $params)
    {
        try {
            $miniProgram = $this->miniProgram($enterpriseId);

            $scene = encryptShortLink($params);

            if (empty($scene)) {
                notice(800006);
            }

            $response = $miniProgram->app_code->getUnlimit($scene, [
                'page'       => $path,
                'width'      => 100,
                'line_color' => [
                    'r' => 105,
                    'g' => 166,
                    'b' => 134,
                ],
            ]);

            if (is_array($response) && isset($response['errcode']) && $response['errcode'] > 0) {
                throw new Error($response['errcode'], $response['errmsg'], [], 422);
            }

            if ($response instanceof StreamResponse) {
                $content = $response->getBodyContents();
                $url     = (new QCloudCos())->uploadContent($content);

                return $url ?? '';
            }

        } catch (\Exception $e) {
            Log::notice('生成场景二维码', []);
            notice(800005, [], ['code:' . $e->getCode(), 'message:' . $e->getMessage()]);
        }

        return '';
    }

    protected function getParams(Request $request)
    {
        // 接受参数
        $input = $request->all();
        // 参数验证规则
        $rules = [
            'enterprise_id' => 'required',
            'path'          => 'required',
            'params'        => 'array|required',
        ];

        // 提示语
        $message = [
            'enterprise_id.required' => '企业Id不能为空',
            'path.required'          => '小程序页面路径不能为空',
            'params.*'               => '参数必填且为数组格式',
        ];

        $this->validate($request, $rules, $message);

        return $input;
    }

}
