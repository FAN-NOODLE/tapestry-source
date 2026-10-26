<?php


namespace App\Http\Controllers\Service\V1\Contact;

use App\Http\Controllers\Service\V1\BaseController;
use App\Jobs\SyncStoreWorkerJob;
use App\Library\Uc\Member;
use App\Models\ContactQrCodeParams;
use App\Models\SyncOrgLog;
use Illuminate\Http\Request;
use Illuminate\Validation\ValidationException;


class SyncController extends BaseController
{

    /**
     * 同步组织结构
     * @throws ValidationException
     * @author: DLong
     */
    public function syncOrg()
    {
        $params = $this->validate($this->request, [
            'enterprise_id' => 'required|integer',
            'qy_domain' => 'required',
            'admin_id' => 'required|integer',
            'admin_name' => 'required',
            'app_identifier' => 'required',
        ], [
            'enterprise_id.required' => '[参数]enterprise_id必填',
            'qy_domain.required' => '[参数]qy_domain必填',
            'enterprise_id.integer' => '[参数]enterprise_id需为integer',
            'admin_id.required' => '[参数]admin_id必填',
            'admin_id.integer' => '[参数]admin_id需为integer',
            'admin_name.required' => '[参数]admin_name必填',
            'app_identifier.required' => '[参数]app_identifier必填'
        ]);

        //当天有同步未完成的禁止重复操作
        $count = SyncOrgLog::query()->where([
            ['enterprise_id', '=', $params['enterprise_id']],
            ['status', '<>', SyncOrgLog::STATUS_UC_OK],
        ])->count();
        if ($count) {
            notice(720000, ['同步操作还未完成，请耐心等待！']);
        }
        //调UC同步通讯录
        $ucResult = (new Member())->syncContact($params['qy_domain']);

        if (empty($ucResult['requestId'])) {
            notice(720000, ['uc同步失败！']);
        }
        //同步操作记录
        $qyDomain = $params['qy_domain'];
        unset($params['qy_domain']);
        unset($params['app_identifier']);
        $params['source_type'] = SyncOrgLog::SOURCE_TYPE_MANUAL;
        SyncOrgLog::query()->create($params);

        return success(['data' => []]);
    }

    /**
     * 获取上次同步时间
     * @throws ValidationException
     * @author: DLong
     */
    public function syncLastTime()
    {
        $params = $this->validate($this->request, [
            'enterprise_id' => 'required|integer',
        ], [
            'enterprise_id.required' => '[参数]enterprise_id必填'
        ]);

        $data = SyncOrgLog::query()
            ->where([
                ['enterprise_id', '=', $params['enterprise_id']],
                ['source_type', '=', SyncOrgLog::SOURCE_TYPE_MANUAL],
                ['status', '=', SyncOrgLog::STATUS_UC_OK],
            ])
            ->orderBy('created_at', 'desc')
            ->first();

        return success(['data' => $data ?? []]);
    }

    /**
     * 当前同步状态
     * @throws ValidationException
     * @author: DLong
     */
    public function status()
    {
        $params = $this->request->input();

        $syncStatus = SyncOrgLog::query()
            ->where([
                ['enterprise_id', '=', $params['enterprise_id']],
                ['source_type', '=', SyncOrgLog::SOURCE_TYPE_MANUAL]
            ])
            ->orderBy('created_at', 'desc')
            ->first(['id', 'status']);
        if ($syncStatus['status'] == 1) {
            self::syncGuideName($params['qy_domain'], $params['enterprise_id']);
        }
        return success(['data' => $syncStatus['status'] ?? SyncOrgLog::STATUS_UC_OK ? 'ok' : 'processing']);
    }

    public static function syncGuideName($qyDomain, $epId)
    {
        dispatch(new SyncStoreWorkerJob($qyDomain, $epId));
    }

    /**
     * 获取员工联系我二维码
     * @throws ValidationException
     * @author: DLong
     */
    public function getContactCode()
    {
        $params = $this->validate($this->request, [
            'enterprise_id' => 'required|integer',
            'object' => 'required',
        ], [
            'enterprise_id.required' => '[参数]enterprise_id必填',
            'object.required' => '[参数]object必填',
        ]);
        $data = ContactQrCodeParams::query()
            ->with('autoContactCode')
            ->where([
                ['enterprise_id', '=', $params['enterprise_id']],
                ['type', '=', ContactQrCodeParams::MEMBER],
                ['object', '=', $params['object']],
            ])
            ->first(['contact_qr_code_id', 'object', 'type', 'enterprise_id']);

        return success(['data' => $data['autoContactCode']['transform_url'] ?? []]);
    }

    public function test(Request $request)
    {

    }
}