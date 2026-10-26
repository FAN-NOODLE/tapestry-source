<?php

namespace App\Http\Controllers\Service\V1\SubscribeMessage;

use App\Http\Controllers\Controller;
use App\Models\SubscribeMessageRecord;
use App\Models\SubscribeMessageTemplate;
use App\Sdk\Active\CrmHelper;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class SubscribeMessageController extends Controller
{

    /**
     * 获取订阅消息配置
     * @param Request $request
     * @return array
     */
    public function getConfig(Request $request): array
    {
        $input = $request->only(['enterprise_id']);

        $list      = SubscribeMessageTemplate::query()->where('enterprise_id', $input['enterprise_id'])->get();
        $returnArr = [];
        if ($list->count() > 0) {
            foreach ($list as $val) {
                $touchStrArr = explode(',', $val['touch_str']);
                foreach ($touchStrArr as $touch) {
                    if (!isset($returnArr[$touch])) {
                        $returnArr[$touch] = [$val['pri_tmpl_id']];
                    } else {
                        $returnArr[$touch][] = $val['pri_tmpl_id'];
                    }
                    $returnArr[$touch] = array_values(array_unique($returnArr[$touch]));
                }
            }
        }

        return success($returnArr);
    }


    /**
     * 获取订阅消息配置
     * @param Request $request
     * @return JsonResponse
     */
    public function subscribe(Request $request): JsonResponse
    {
        $input = $request->only(['enterprise_id', 'tmpl_list', 'union_id', 'open_id', 'rel_id', 'is_always']);

        $urlByTmplArr = array_column($input['tmpl_list'], 'url', 'pri_tmpl_id');
        $list         = SubscribeMessageTemplate::query()->whereIn('pri_tmpl_id', array_keys($urlByTmplArr))
            ->get(['pri_tmpl_id', 'page_url', 'type'])->toArray();
        foreach ($list as $info) {
            // 要处理兼容上报的路径地址
            if (strpos($urlByTmplArr[$info['pri_tmpl_id']], '?') === 0){
                $pageUri = $info['page_url'] . $urlByTmplArr[$info['pri_tmpl_id']];
            }else{
                $pageUri = $urlByTmplArr[$info['pri_tmpl_id']];
            }
            $addArr[] = [
                'enterprise_id' => $input['enterprise_id'],
                'unionid'       => $input['union_id'],
                'openid'        => $input['open_id'],
                'pri_tmpl_id'   => $info['pri_tmpl_id'],
                'is_always'     => $input['is_always'] ?? 0,
                'type'          => $info['type'],
                'page'          => $pageUri,
                'rel_id'        => $input['rel_id'],
                'created_at'    => time(),
                'updated_at'    => time()

            ];
        }
        if (!empty($addArr)) SubscribeMessageRecord::query()->insert($addArr);

        return successJson();
    }
}