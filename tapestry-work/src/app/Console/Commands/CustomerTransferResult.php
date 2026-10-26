<?php

namespace App\Console\Commands;

use App\Library\Uc\QyWechat;
use App\Models\Scrm\CustomerOnJobInherit;
use App\Models\Scrm\CustomerOnJobInheritRecord;
use Exception;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Log;

class CustomerTransferResult extends Command
{
    /**
     * @var string
     */
    protected $signature = 'auto:customer-transfer_result';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '查询客户接替状态';

    /**
     * @return void
     */
    public function handle()
    {
        $inherit = CustomerOnJobInherit::query()
            ->where('status', CustomerOnJobInherit::STATUS_ON)
            ->where('created_at', '<', time() - 60)
            ->with('record')
            ->get();

        foreach ($inherit as $row) {

            $record = $row->record;
            $record = collect($record->toArray())->unique('handover_member_id')->all();

            //企微接口查询
            foreach ($record as $value) {
                $this->transferResult($value, $row->enterprise_id, '');
            }

            //统计处理数据
            $this->statistics($row->id);
        }
    }

    /**
     * @param $value
     * @param $enterprise_id
     * @param $cursor
     *
     * @return bool
     */
    public function transferResult($value, $enterprise_id, $cursor): bool
    {
        try {
            $res = (new QyWechat())->getQyResult([
                'handover_userid' => $value['handover_userid'],
                'takeover_userid' => $value['takeover_userid'],
                'cursor'          => $cursor
            ], [
                'path'          => 'cgi-bin/externalcontact/transfer_result',
                'method'        => 'POST',
                'enterprise_id' => $enterprise_id
            ]);

            Log::info('在职继承查询客户接替状态查询入参enterprise_id' . $enterprise_id);
            Log::info('在职继承查询客户接替状态查询入参' . json_encode($value));
            Log::info('在职继承查询客户接替状态返回信息' . json_encode($res));

            if ($res['errcode'] == 0) {
                foreach ($res['customer'] as $item) {
                    CustomerOnJobInheritRecord::query()->where([
                        'inherit_id'      => $value['inherit_id'],
                        'external_userid' => $item['external_userid'],
                    ])->where('updated_at','<',time() - 60)->update([
                        'status'        => $item['status'],
                        'errmsg'        => $this->errMsg($item['status']),
                        'takeover_time' => $item['takeover_time'] ?? 0
                    ]);
                }
            }
            if ($res['next_cursor']) {
                Log::info('在职继承查询客户接替状态返回了下一页数据' . $res['next_cursor']);
                $this->transferResult($value, $enterprise_id, $res['next_cursor']);
            }

        } catch (Exception $e) {
            Log::info('用户在职继承结果获取失败,错误原因' . $e->getMessage());
        }
        //如果该分配人和接收人的记录中存在25小时仍然没有更新过的记录则为分配异常情况，默认为分配失败（企微默认24小时分配结束)
        CustomerOnJobInheritRecord::query()->where([
            'inherit_id' => $value['inherit_id'],
            'handover_userid' => $value['handover_userid'],
            'takeover_userid' => $value['takeover_userid'],
            'status' => CustomerOnJobInheritRecord::STATUS_WAIT,
        ])->where('updated_at', '<', time() - 3600 * 25)
            ->update([
                'status' => CustomerOnJobInheritRecord::STATUS_RECORD_ERROR,
                'errmsg' => "超过24小时未查询到企微分配结果",
                'takeover_time' => $item['takeover_time'] ?? 0
            ]);
        return true;
    }

    /**
     * @param $status
     *
     * @return string
     */
    public function errMsg($status)
    {
        $msg = [
            '1' => '接替完毕',
            '2' => '等待接替',
            '3' => '客户拒绝',
            '4' => '接替成员客户达到上限'
        ];

        return $msg[$status];
    }

    /**
     * @return void
     * 数据统计
     */
    public function statistics($inherit_id)
    {
        $waitSign     = false;    //分配中标识
        $errNum     = 0;    //分配失败的条数
        $inherit = CustomerOnJobInheritRecord::query()->where('inherit_id', $inherit_id)->get();
        foreach ($inherit as $v) {

            //2-等待接替
            if ($v->status == CustomerOnJobInheritRecord::STATUS_WAIT) {
                //有一个等待接替则标识整条数据为分配中直接跳出
                $waitSign = true;
                break;
            }
            //有处理失败的 3客户拒绝 4接替成员客户达到上限
            if ($v->status == CustomerOnJobInheritRecord::STATUS_REJECT
                || $v->status == CustomerOnJobInheritRecord::STATUS_UP
                || $v->status == CustomerOnJobInheritRecord::STATUS_WAIT_ERROR
                || $v->status == CustomerOnJobInheritRecord::STATUS_RECORD_ERROR) {
                $errNum++;
            }
        }

        if ($waitSign) {
            $msg = ['status' => 0, 'msg' => '分配中'];
        } else {
            if ($errNum == 0) {
                $msg = ['status' => 1, 'msg' => '分配成功'];
            } else {
                $msg = ['status' => 2, 'msg' => $errNum . '名分配失败'];
            }
        }

        CustomerOnJobInherit::query()->where('id', $inherit_id)->update([
            'status' => $msg['status'],
            'errmsg' => $msg['msg'],
        ]);
    }
}