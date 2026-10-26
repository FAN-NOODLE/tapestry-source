<?php
/**
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/2/28
 * Time: 16:32
 */

namespace App\Console\Commands;

use App\Library\Uc\Member;
use App\Models\MemberExtend;
use App\Models\ContactQrCode;
use App\Models\ContactQrCodeParams;
use CrsSdk\Boss\Enterprise;
use Illuminate\Console\Command;

class UpdateUCGuideQr extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'update-uc-guide-qr {type}';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '更新uc导购二维码(参数传递 恰恰,蒙牛传递参数1,其他传递2)';

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

        $qyDomains = [];
        if (!empty($MemberExtend)) {
            $guideEnterpriseIds = array_unique(array_column($MemberExtend, 'enterprise_id'));

            foreach ($guideEnterpriseIds as $enterpriseId) {
                $enterpriseData           = (new Enterprise())->detail(['enterprise_id' => $enterpriseId]);
                $enterpriseData           = $enterpriseData['data'];
                $qyDomains[$enterpriseId] = $enterpriseData['domain'];
            }
        }

        $enterpriseUids = [];
        // 组装用户请求
        foreach ($MemberExtend as $value) {
            $enterpriseUids[$value['enterprise_id']][] = $value['mem_uid'];
        }
        $members   = [];
        $config = [];
        // 请求uc获取用户userId;
        foreach ($enterpriseUids as $enterpriseId => $memUids) {
            $domain = $qyDomains[$enterpriseId];

            $params = [
                'memUids'  => $memUids,
                'pageNum'  => 1,
                'pageSize' => count($memUids),
                'memAll'   => 1
            ];

            $result = (new Member())->ignorePermissionsList($domain, '', $params);

            foreach ($result['list'] ?? [] as $member) {
                $memUid                          = $member['memUid'];
                $members[$enterpriseId][$memUid] = $member['memUserid'];
            }

            $data = ContactQrCodeParams::query()
                ->join('contact_qr_code', 'contact_qr_code_params.contact_qr_code_id', '=', 'contact_qr_code.id')
                ->where([
                    'contact_qr_code.enterprise_id'        => $enterpriseId,
                    'contact_qr_code.source_type'          => ContactQrCode::SOURCE_TYPE_AUTO_GUIDE,
                    'contact_qr_code.type'                 => ContactQrCode::CODE_TYPE_ONE,
                    'contact_qr_code_params.enterprise_id' => $enterpriseId,
                    'contact_qr_code_params.type'          => ContactQrCodeParams::MEMBER,
                ])
                ->whereIn('contact_qr_code_params.object', $memUids)
                ->get(['contact_qr_code.id', 'contact_qr_code_params.object as mem_uid','contact_qr_code.config_id'])
                ->toArray();

            foreach ($data as $datum) {
                $config[$enterpriseId][$datum['mem_uid']] = [
                    'id'=>$datum['id'],
                    'config_id'=>$datum['config_id']
                ];
            }
        }

        $i = 1;
        foreach ($MemberExtend as $value) {
            echo '处理进度' . $i . '/' . count($MemberExtend) . "导购数据\r\n";
            $domain    = $qyDomains[$value['enterprise_id']];
            $memUserid = $members[$value['enterprise_id']][$value['mem_uid']]??'';
            $ecwScene  = $config[$value['enterprise_id']][$value['mem_uid']] ?? [];

            if (empty($memUserid)) {
                echo "未找到相关user_id:企业" . $value['enterprise_id'] . ",用户mem_uid: " . $value['mem_uid'] . "\r\n";
                continue;
            }
            if (empty($ecwScene)) {
                echo "未找到相关二维码企业id:" . $value['enterprise_id'] . ",用户mem_uid: " . $value['mem_uid'] . "\r\n";
                continue;
            }

            (new Member())->updateContactQrCode($domain, [
                'id'  => $ecwScene['config_id'],
                'ecwType'   => 1,
                "ecwState"  => $ecwScene['id'],
                'memUserid' => $memUserid,
            ]);
            $i++;
        }
    }


}
