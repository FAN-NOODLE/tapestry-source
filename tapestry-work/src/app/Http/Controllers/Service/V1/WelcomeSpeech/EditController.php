<?php

namespace App\Http\Controllers\Service\V1\WelcomeSpeech;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Models\Channel;
use App\Models\ContactQrCode;
use App\Models\WelcomeSpeech;
use App\Services\MemberService;
use Illuminate\Http\Request;

class EditController extends BaseController
{
    /**
     * 编辑欢迎语
     *
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request)
    {
        $id            = $request->post('id', 0);
        $enterpriseId  = $request->post('enterprise_id', 0);
        $qCode         = $request->post('qr_code', 0);
        $channel_id    = $request->post('channel_id', 0);
        $type          = $request->post('type', 1);
        $content       = $request->post('content', '');
        $object        = $request->post('object', '');
        $originalData  = $request->post('original_data', '');
        $isDefault     = $request->post('is_default', 0);
        $qyDomain      = $request->post('qy_domain', '');      // 企业domain
        $appIdentifier = $request->post('app_identifier', ''); // 应用标识

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

        if (WelcomeSpeech::IS_DEFAULT_YES == $isDefault) {
            // 默认欢迎语
            $defaultWelcomeSpeech = WelcomeSpeech::query()
                ->where('enterprise_id', $enterpriseId)
                ->where('is_default', WelcomeSpeech::IS_DEFAULT_YES)
                ->first();

            if ($defaultWelcomeSpeech) {
                WelcomeSpeech::query()
                    ->where('enterprise_id', $enterpriseId)
                    ->where('is_default', WelcomeSpeech::IS_DEFAULT_YES)
                    ->update([
                        'type'          => $type,
                        'content'       => $content,
                        'object'        => $object,
                        'original_data' => $originalData,
                    ]);
            } else {
                WelcomeSpeech::query()
                    ->create([
                        'enterprise_id' => $enterpriseId,
                        'type'          => $type,
                        'content'       => $content,
                        'object'        => $object,
                        'original_data' => $originalData,
                        'is_default'    => WelcomeSpeech::IS_DEFAULT_YES,
                    ]);
            }
        } else {
            //$contactQrCode = ContactQrCode::query()->find($channel_id);
//            $contactQrCode = Channel::query()->find($channel_id);
//
//            if (!$contactQrCode) {
//                throw new Error(700000, '新增欢迎语失败，选择渠道不存在');
//            }

            WelcomeSpeech::query()
                ->where('id', $id)
                ->update([
                    'channel_id'         => $channel_id,
                    'contact_qr_code_id' => $qCode,
                    'type'               => $type,
                    'content'            => $content,
                    'object'             => $object,
                    'original_data'      => $originalData,
                ]);
        }

        return success([]);
    }
}
