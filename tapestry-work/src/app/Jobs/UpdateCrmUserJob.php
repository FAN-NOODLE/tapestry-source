<?php

namespace App\Jobs;

use App\Services\CustomerScrmService;
use Illuminate\Support\Carbon;
use Illuminate\Support\Facades\Log;

class UpdateCrmUserJob extends Job
{

    /**
     * 队列名称
     *
     * @var string
     */
    public $queue = 'update_user_info_for_crm';

    protected $input;

    /**
     * Create a new job instance.
     *
     * @return void
     */
    public function __construct($input)
    {
        $this->input = $input;
    }

    /**
     * Execute the job.
     *
     * @return void
     */
    public function handle()
    {
        Log::info('更新crm用户信息--队列开始', ['$input' => $this->input, 'time' => Carbon::now()]);

        $crmParams = [
            'unionid' => $this->input['unionid'],
            'mobile'  => isset($this->input['mobile']) ? $this->input['mobile'] : '',
//            'eu_id'   => $this->input['eu_id'],
        ];
        $param     = ['enterprise_id' => $this->input['enterprise_id'],];
        $result    = CustomerScrmService::updateUser($crmParams, $param);

        Log::info('更新crm用户信息--队列结束',
            ['input' => $this->input, 'time' => Carbon::now(), 'Result' => $result ?? '失败']
        );
    }
}
