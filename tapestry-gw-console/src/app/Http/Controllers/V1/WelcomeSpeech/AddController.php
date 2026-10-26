<?php

namespace App\Http\Controllers\V1\WelcomeSpeech;

use Illuminate\Http\Request;
use CrsSdk\Work\WelcomeSpeech;
use Illuminate\Http\JsonResponse;
use Illuminate\Support\Facades\Validator;
use App\Http\Controllers\V1\BaseController;

class AddController extends BaseController
{
    /**
     * 新增欢迎语
     * @param Request $request
     * @return JsonResponse
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $requests = $request->only([
            'qr_code', 'type', 'content', 'mediaId', 'picurl', 'title', 'desc', 'url', 'picMediaId',
            'appid', 'page', 'original_data', 'channel_id'
        ]);

        $rules = [
            'channel_id' => 'required',
            'qr_code'    => 'required',
            'type'       => 'required',
            'content'    => 'required',
            'mediaId'    => 'required_if:type,2',
            'picurl'     => 'required_if:type,3',
            'title'      => 'required_if:type,3',
            'desc'       => 'required_if:type,3',
            'url'        => 'required_if:type,3',
            'picMediaId' => 'required_if:type,4',
            'appid'      => 'required_if:type,4',
            'page'       => 'required_if:type,4',
        ];

        $message = [
            'channel_id.required'    => '渠道id不能为空',
            'qr_code.required'       => '联系我二维码不能为空',
            'type.required'          => '欢迎语类型不能为空',
            'content.required'       => '消息内容不能为空',
            'mediaId.required_if'    => '图片的media_id不能为空',
            'picurl.required_if'     => '图文消息封面的url不能为空',
            'title.required_if'      => '图文消息标题不能为空',
            'desc.required_if'       => '图文消息的描述不能为空',
            'url.required_if'        => '图文消息的链接不能为空',
            'picMediaId.required_if' => '小程序消息封面的media_id不能为空',
            'appid.required_if'      => '小程序appid不能为空',
            'page.required_if'       => '小程序page路径不能为空',
        ];

        Validator::make($requests, $rules, $message)->validate();

        switch ($requests['type']) {
            case 2: // 图片
                $object = json_encode([
                    'mediaId' => $requests['mediaId'],
                ], JSON_UNESCAPED_UNICODE);
                break;
            case 3: // 网页
                $object = json_encode([
                    'picurl' => $requests['picurl'],
                    'title'  => $requests['title'],
                    'desc'   => $requests['desc'],
                    'url'    => $requests['url'],
                ], JSON_UNESCAPED_UNICODE);
                break;
            case 4: // 小程序
                $object = json_encode([
                    'picMediaId' => $requests['picMediaId'],
                    'title'      => $requests['title'],
                    'appid'      => $requests['appid'],
                    'page'       => $requests['page'],
                ], JSON_UNESCAPED_UNICODE);
                break;
            default:
                $object = '';
                break;
        }

        $params = [
            'enterprise_id'  => $this->enterpriseId,
            'qr_code'        => $requests['qr_code'],
            'type'           => $requests['type'],
            'content'        => $requests['content'],
            'channel_id'     => $requests['channel_id'],
            'object'         => $object,
            'original_data'  => json_encode(($requests['original_data'] ?? ''), JSON_UNESCAPED_UNICODE),
            'qy_domain'      => $this->qyDomain,
            'app_identifier' => $this->appIdentifier,
        ];

        (new WelcomeSpeech())->add($params);

        return success();
    }
}
