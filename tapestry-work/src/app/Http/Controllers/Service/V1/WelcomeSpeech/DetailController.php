<?php

namespace App\Http\Controllers\Service\V1\WelcomeSpeech;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Models\Channel;
use App\Models\ContactQrCode;
use App\Models\WelcomeSpeech;
use Illuminate\Http\Request;

class DetailController extends BaseController
{
    /**
     * 欢迎语详情
     *
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request)
    {
        $enterpriseId = $request->get('enterprise_id', 0);
        $id           = $request->get('id', 0);
        $isDefault    = $request->get('is_default', 0);

        if (!$enterpriseId) {
            throw new Error(700000, '企业ID不能为空');
        }

        if (WelcomeSpeech::IS_DEFAULT_YES == $isDefault) {
            $welcomeSpeech = WelcomeSpeech::query()
                ->where('enterprise_id', $enterpriseId)
                ->where('is_default', WelcomeSpeech::IS_DEFAULT_YES)
                ->first();

            if (!$welcomeSpeech) {
                $welcomeSpeech = WelcomeSpeech::DEFAULT_CONTENT;
                $welcomeSpeech['type']          = WelcomeSpeech::NONE_LINK;
                $welcomeSpeech['object']        = [];
                $welcomeSpeech['is_default']    = WelcomeSpeech::IS_DEFAULT_YES;
                $welcomeSpeech['enterprise_id'] = $enterpriseId;
            } else {
                $welcomeSpeech = $welcomeSpeech->toArray();

                $welcomeSpeech['object'] = json_decode($welcomeSpeech['object'], true);
            }
        } else {
            if (!$id) {
                throw new Error(700001, '欢迎语ID不能为空');
            }

            $welcomeSpeech = WelcomeSpeech::query()
                ->where('id', $id)
                ->where('enterprise_id', $enterpriseId)
                ->where('is_default', WelcomeSpeech::IS_DEFAULT_NO)
                ->first();

            if (!$welcomeSpeech) {
                throw new Error(700002, '该欢迎语不存在');
            }

            $welcomeSpeech = $welcomeSpeech->toArray();

            $welcomeSpeech['object'] = json_decode($welcomeSpeech['object'], true);

            $contactQrCode = Channel::query()->where('id', $welcomeSpeech['channel_id'])->first();

            $welcomeSpeech['channel'] = $contactQrCode['channel'];
        }

        return success(['data' => $welcomeSpeech]);
    }
}
