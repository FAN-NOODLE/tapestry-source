<?php

namespace App\Http\Controllers\V1\CustomerMaintenance\Msg;

use App\Models\QYwxEmsg;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Validator;
use App\Http\Controllers\V1\BaseController;

class UpdateStatusController extends BaseController
{
    protected $requireLogin = false;

    public function __construct(Request $request)
    {
        parent::__construct();
        $this->request = $request;
    }

    /**
     * 修改群发消息 - 只有草稿状态下的消息才能被修改
     *
     * @param string $tsid
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function index($tsid = "")
    {
        if (!$tsid) error(1000005);
        $input = $this->request->only(
            'status', "num_fact"
        );

        $rules    = [
            'status'   => 'integer|in:1,2,3',
            'num_fact' => 'integer',
        ];
        $messages = [
            'status.integer'   => '状态值为数值类型仅允许123',
            'num_fact.integer' => '实际触达人数为数值类型'
        ];

        $validator = Validator::make($input, $rules, $messages);
        if ($validator->fails()) {
            notice(1000010, $validator->errors()->all());
        }

        $data = $input;

        $where = [['tsid', $tsid]];
        DB::beginTransaction();
        try {
            QYwxEmsg::query()->where($where)->update($data);
            $id = $tsid;
            DB::commit();
        } catch (\Exception $e) {
            DB::rollBack();
            Log::info("状态更新失败：[{$e->getCode()}] " . $e->getMessage());
            notice($e->getCode());
        }
        return success(['id' => $id]);
    }
}
