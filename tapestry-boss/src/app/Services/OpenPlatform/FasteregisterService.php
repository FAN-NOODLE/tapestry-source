<?php
/*
 * This source file is subject to the MIT license that is bundled
 * with this source code in the file LICENSE.
 */

namespace App\Services\OpenPlatform;

use App\Exceptions\Error;
use App\Models\RegisterProgramAgent;
use App\Services\BaseService;
use Illuminate\Support\Facades\Log;

class FastRegisterService extends BaseService
{
    /**
     * @var registerProgramAgent
     */
    protected $registerProgramAgent;


    public function __construct(RegisterProgramAgent $registerProgramAgent)
    {
        $this->registerProgramAgent = $registerProgramAgent;
    }

    function updateRegisterAgent($data)
    {

        Log::info('微信授权回调进来了');

        if (empty($data) || !isset($data['info'])) return false;

        try {

            $where = [
                'legal_persona_wechat' => $data['info']['legal_persona_wechat'] ?? '',
                'legal_persona_name'   => $data['info']['legal_persona_name'] ?? '',
                'name'                 => $data['info']['name'] ?? '',
            ];
            if ($data['status'] == 0) {
                $updateData = [
                    'app_id'     => $data['appid'],
                    'status'    => $this->registerProgramAgent::STATUS_APPLY_SUCCESS,
                    'auth_code' => $data['auth_code'],
                ];
            } else {
                $updateData = [
                    'status'      => $this->registerProgramAgent::STATUS_APPLY_FAIL,
                    'fail_reason' => $this->registerProgramAgent::REGISTER_CODE_MAP[$data['status']] ?? $data['msg'],
                ];
            }

            // 审核数据更新
            $this->registerProgramAgent->updateByData($where, $updateData);

        } catch (\Exception $e) {
            Log::error('审核数据更新失败 data' . json_encode($data, JSON_UNESCAPED_UNICODE) . 'ERROR ' . $e->getMessage());
        }
        return true;
    }
}
