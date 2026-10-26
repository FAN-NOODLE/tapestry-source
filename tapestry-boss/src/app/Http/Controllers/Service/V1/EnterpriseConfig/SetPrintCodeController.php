<?php

namespace App\Http\Controllers\Service\V1\EnterpriseConfig;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\EnterpriseConfig;
use App\Models\EnterpriseThirdPartAccount;
use CrsSdk\Common\QCloudCos;
use Intervention\Image\Facades\Image;

class SetPrintCodeController extends BaseController
{
    public function setPrintCode()
    {
        $this->validate($this->request, [
            'enterprise_id' => 'required',
        ], [
            'enterprise_id.required' => '企业ID必传',
        ]);

        $enterpriseID = $this->request->input('enterprise_id');

        $qrcode = EnterpriseConfig::query()
            ->where('enterprise_id', '=', $enterpriseID)
            ->where('type', '=', EnterpriseConfig::PRINT_QRCODE)
            ->first();

        $account = EnterpriseThirdPartAccount::query()
            ->where('enterprise_id', '=', $enterpriseID)
            ->where('authorize_type', '=', EnterpriseThirdPartAccount::AUTHORIZE_TYPE_THIRD)
            ->where('platform', '=', EnterpriseThirdPartAccount::MINI_PROGRAM_ACCOUNT)
            ->where('status', '=', EnterpriseThirdPartAccount::STATUS_ACTIVE)
            ->where('type', '=', EnterpriseThirdPartAccount::TYPE_SHOP)
            ->first();

        if (empty($account)) {
            error(201115);
        }

        if (empty($qrcode)) {
            // 背景图
            $bgImage = Image::make(resource_path('image/background.png'));
            // logo图处理
            $logoImage    = Image::make($account->logo)->resize(70, 70);
            $r            = $logoImage->width() / 2;
            $newLogoImage = Image::canvas(70, 70);
            for ($x = 0; $x < $logoImage->width(); $x++) {
                for ($y = 0; $y < $logoImage->height(); $y++) {
                    if (((($x - $r) * ($x - $r) + ($y - $r) * ($y - $r)) < ($r * $r))) {
                        $c = $logoImage->pickColor($x, $y, 'array');
                        $newLogoImage->pixel($c, $x, $y);
                    }
                }
            }
            $bgImage->insert($newLogoImage, 'top-left', 100, 254);
            $bgImage->text($account->name, 190, 290, function ($font) {
                $font->file(resource_path('font/pingfang.ttf'));
                $font->size(34);
                $font->valign('middle');
            });

            // 组装小程序码
            $codeImage = Image::make($account->main_qrcode)->resize(400, 400);
            $bgImage->insert($codeImage, 'bottom', 0, 1130);

            $url = (new QCloudCos())->uploadContent($bgImage->stream('png', 60)->getContents(), $account->name);

            $qrcode = EnterpriseConfig::query()->create([
                'key'           => EnterpriseConfig::PRINT_QRCODE_KEY,
                'type'          => EnterpriseConfig::PRINT_QRCODE,
                'name'          => $url,
                'desc'          => '打印小票小程序码',
                'value'         => 1,
                'enterprise_id' => $enterpriseID,
            ]);
        }

        return success([
            'data' => ['qr_code' => $qrcode->name]
        ]);
    }
}
