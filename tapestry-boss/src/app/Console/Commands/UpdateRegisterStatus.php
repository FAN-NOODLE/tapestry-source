<?php

namespace App\Console\Commands;


use EasyWeChat\Kernel\Support\XML;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Log;
use Overtrue\LaravelWeChat\Facade as LaravelWechat;
use App\Models\RegisterProgramAgent;
use App\Services\OpenPlatform\FastRegisterService;
use App\Services\OpenPlatform\EnterpriseService;

class UpdateRegisterStatus extends Command
{


    /**
     * @var string
     */
    protected $signature = 'auto:update_register_status';

    /**
     * @var string
     */
    protected $description = '定时更新代注册小程序审核状态';


    public function handle(RegisterProgramAgent $model)
    {
        // 查询审核中列表
        $list = $model->getApplyListByStatus();

        if (!empty($list)){

            $formatData = [
                'status'   => $model::STATUS_APPLY_FAIL,
            ];
            foreach ($list as $key => $item){
                //调用开放平台
                try {
                    $result = LaravelWechat::openPlatform('default')->component->getRegistrationStatus(
                        $item['name'],$item['legal_persona_wechat'],$item['legal_persona_name']
                    );
                    if (empty($result)){
                        if ($result['errcode'] == 0 ) {
                            $formatData['status'] = $model::STATUS_APPLY_SUCCESS;
                        }else{
                            $formatData['fail_reason'] = RegisterProgramAgent::REGISTER_CODE_MAP[$result['errcode']] ?? '其它错误' ;
                        }
                        $model->updateData($item['id'],$formatData);
                    }
                }catch (\Exception $e){
                    Log::info('查看注册异常ERROR:'.$e->getMessage());
                }
            }
        }

        return true;
    }
}
