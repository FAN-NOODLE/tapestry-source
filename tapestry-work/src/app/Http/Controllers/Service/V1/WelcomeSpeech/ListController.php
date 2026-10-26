<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/12/10
 * Time: 16:24
 */

namespace App\Http\Controllers\Service\V1\WelcomeSpeech;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\WelcomeSpeech;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\DB;

class ListController extends BaseController
{
    /**
     * 欢迎语列表
     *
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request)
    {
        $size         = $request->get('size', 10);
        $channel      = $request->get('channel', '');
        $enterpriseId = $request->get('enterprise_id', 0);

        $list         = WelcomeSpeech::query()
            ->join('channel', function ($join) {
                $join->on([
                    ['welcome_speech.channel_id', '=', 'channel.id']
                ]);
            })
            ->where('welcome_speech.enterprise_id', $enterpriseId)
            ->where('welcome_speech.is_default', WelcomeSpeech::IS_DEFAULT_NO)
            ->when($channel != '', function ($query) use ($channel) {
                $query->where('channel.channel', 'like', '%' . $channel . '%');
            })
            ->where('welcome_speech.deleted_at', 0)
            ->where('channel.deleted_at', 0)
            ->orderByDesc('welcome_speech.id')
            ->paginate($size, [
                'welcome_speech.id',
                'welcome_speech.created_at',
                'welcome_speech.content',
                'channel.channel',
            ])
            ->toArray();

        // 默认欢迎语数据
        $default = WelcomeSpeech::query()
            ->select(DB::raw('id, content, created_at'))
            ->where('enterprise_id', $enterpriseId)
            ->where('is_default', WelcomeSpeech::IS_DEFAULT_YES)
            ->first();

        return success([
            'default' => $default ? $default->toArray() : WelcomeSpeech::DEFAULT_CONTENT,
            'data'    => $list['data'],
            'meta'    => [
                'total' => $list['total'],
                'size'  => (int)$list['per_page'],
                'page'  => $list['current_page']
            ]
        ]);
    }
}
