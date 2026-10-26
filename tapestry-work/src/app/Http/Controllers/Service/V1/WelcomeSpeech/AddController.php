<?php

namespace App\Http\Controllers\Service\V1\WelcomeSpeech;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Models\ContactQrCode;
use App\Models\WelcomeSpeech;
use App\Services\MemberService;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\DB;

class AddController extends BaseController
{
    /**
     * 新增欢迎语
     *
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request)
    {
        $enterpriseId  = $request->post('enterprise_id', 0);
        $qrCode        = $request->post('qr_code', 0);
        $channel_id    = $request->post('channel_id', 0);
        $type          = $request->post('type', 1);
        $content       = $request->post('content', '');
        $object        = $request->post('object', '');
        $originalData  = $request->post('original_data', '');
        $qyDomain      = $request->post('qy_domain', '');      // 企业domain
        $appIdentifier = $request->post('app_identifier', ''); // 应用标识

        if (!$enterpriseId) {
            throw new Error(700000, '企业ID不能为空');
        }

//        if (!$qrCode) {
//            throw new Error(700001, '联系我二维码不能为空');
//        }
        if (!$channel_id) {
            throw new Error(700001, '渠道id不能为空');
        }

        if (!$type) {
            throw new Error(700002, '欢迎语类型不能为空');
        }

        if (!$content) {
            throw new Error(700003, '消息内容不能为空');
        }

//        $contactQrCode = ContactQrCode::query()->find($channel_id);
//
//        if (!$contactQrCode) {
//            throw new Error(700004, '新增欢迎语失败，选择渠道不存在');
//        }

        $welcomeSpeech = WelcomeSpeech::query()
            ->where('channel_id', $channel_id)
            //->where('contact_qr_code_id', $qrCode)
            ->where('enterprise_id', $enterpriseId)
            ->first();

        if ($welcomeSpeech) {
            throw new Error(700005, '该渠道的欢迎语已存在');
        }

        if (WelcomeSpeech::IMAGE == $type) {
            $object = json_decode($object, true);

            $object['mediaId'] = MemberService::getWorkWeChartMediaId($qyDomain, $appIdentifier, $object['mediaId']);

            $object = json_encode($object, JSON_UNESCAPED_UNICODE);
        }

        if (WelcomeSpeech::MINI_PROGRAM == $type) {
            $object = json_decode($object, true);

            $object['picMediaId'] = MemberService::getWorkWeChartMediaId($qyDomain, $appIdentifier, $object['picMediaId']);

            $object = json_encode($object, JSON_UNESCAPED_UNICODE);
        }
        WelcomeSpeech::query()
            ->create([
                'enterprise_id'      => $enterpriseId,
                'contact_qr_code_id' => $qrCode,
                'type'               => $type,
                'channel_id'         => $channel_id,
                'content'            => $content,
                'object'             => $object,
                'original_data'      => $originalData,
            ]);
        return success([]);
    }
}
