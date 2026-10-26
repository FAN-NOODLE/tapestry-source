<?php

namespace App\Http\Controllers\V1\WelcomeSpeech;

use Illuminate\Http\Request;
use CrsSdk\Work\WelcomeSpeech;
use Illuminate\Http\JsonResponse;
use Illuminate\Support\Facades\Validator;
use App\Http\Controllers\V1\BaseController;

class EditController extends BaseController
{
    /**
     * 编辑欢迎语
     * @param Request $request
     * @return JsonResponse
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $requests = $request->only(['id', 'qr_code', 'channel_id', 'type', 'content', 'mediaId', 'picurl', 'title',
                                    'desc', 'url', 'picMediaId', 'appid', 'page', 'is_default', 'original_data']);

        $rules = [
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
            'id'             => $requests['id'] ?? 0,
            'enterprise_id'  => $this->enterpriseId,
            'qr_code'        => $requests['qr_code'] ?? 0,
            'channel_id'     => $requests['channel_id'] ?? 0,
            'type'           => $requests['type'],
            'content'        => $requests['content'],
            'object'         => $object,
            'original_data'  => json_encode(($requests['original_data'] ?? ''), JSON_UNESCAPED_UNICODE),
            'is_default'     => $requests['is_default'] ?? 0,
            'qy_domain'      => $this->qyDomain,
            'app_identifier' => $this->appIdentifier,
        ];

        (new WelcomeSpeech())->edit($params);

        return success([]);
    }
}
