<?php
/**
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/2/28
 * Time: 16:32
 */

namespace App\Console\Commands;

use App\Models\Channel;
use App\Models\MemberExtend;
use App\Models\ContactQrCode;
use App\Models\ContactQrCodeParams;
use App\Services\CommonService;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;

class UpdateGuideQr extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'update-guide-qr {type}';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '迁移导购老数据二维码(参数传递 恰恰,蒙牛传递参数1,其他传递2)';

    public function handle()
    {

        $enterpriseIds = [
            13,
            86
        ];

        $type = $this->argument('type') ?? 0;
        if (empty($type) || !in_array($type, [1, 2])) {
            echo "恰恰,蒙牛传递参数1,其他传递2";
            exit;
        }

        $model = MemberExtend::query();
        if ($type == 1) {
            $model->whereIn('enterprise_id', $enterpriseIds);
        }
        if ($type == 2) {
            $model->whereNotIn('enterprise_id', $enterpriseIds);
        }

        // 获取导购数据
        $MemberExtend = $model
            ->where('role', '导购')
            ->get()->toArray();

        echo '共' . count($MemberExtend) . '导购数据';
        $i = 1;
        foreach ($MemberExtend as $value) {
            echo '处理进度' . $i . '/' . count($MemberExtend) . "导购数据\r\n";
            $this->upGuideQr($value);
            $i++;
        }
    }

    protected function upGuideQr($params)
    {
        $ucMemId      = $params['mem_uid'];
        $enterpriseId = $params['enterprise_id'];

        $sourceType = 3;

        $channelName = '支付即导购';
        //对应渠道是否存在
        $channelCollect = Channel::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('source_type', $sourceType)
            ->get(['id', 'channel']);
        $channelObj     = $channelCollect->where('channel', $channelName)->first();

        if (!$channelObj) {
            //同步渠道
            $channelData = [
                'enterprise_id' => $enterpriseId,
                'channel'       => $channelName,
                'source_type'   => $sourceType,
                'admin_id'      => 0,
                'admin_name'    => '系统'
            ];
            $channelObj  = Channel::query()->create($channelData);
            if (empty($channelObj)) {
                echo '同步联系人二维码渠道异常！';
            }
        }

        // 本地存储二维码
        $contactQr = ContactQrCode::query()->create([
            'enterprise_id' => $enterpriseId,
            'name'          => '',
            'is_verify'     => ContactQrCode::VERIFY_NO,
            'channel_id'    => $channelObj['id'],
            'channel'       => $channelName,
            'type'          => ContactQrCode::CODE_TYPE_ONE,
            'source_type'   => $sourceType,
        ]);

        ContactQrCodeParams::query()->create([
            'enterprise_id'      => $enterpriseId,
            'contact_qr_code_id' => $contactQr['id'],
            'type'               => ContactQrCodeParams::MEMBER,
            'object'             => $ucMemId,
            'object_extra'       => $ucMemId,
        ]);

        // 企业微信二维码链接
        $officialUrl = $params['contact_qr_code'];

        // 企业微信二维码存储到腾讯云
        $transformUrl = CommonService::uploadImage($officialUrl);
        // 更新本地数据
        $contactQr['uc_id']         = '';
        $contactQr['config_id']     = $params['config_id'];
        $contactQr['official_url']  = $officialUrl;
        $contactQr['transform_url'] = $transformUrl;

        $contactQr->save();

    }

}
