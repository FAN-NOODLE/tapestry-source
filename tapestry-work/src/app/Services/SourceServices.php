<?php

namespace App\Services;

use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Redis;

class SourceServices extends BaseService
{

    /**
     * @throws \Exception
     */
    public function sysSourceMediaId($list)
    {
        // 获取企业素材信息
        if (!empty($list)) {
            foreach ($list as $key => $value) {
                $this->addMediaIdRedis($value['uc_eid'], $value['image_media_url'], $value['scid']);
            }
        }
        return true;
    }

    /**
     * 将素材图片mediaId存入缓存
     *
     * @param $qyDomain
     * @param $pic
     * @param $scid
     *
     * @throws \Exception
     */
    public function addMediaIdRedis($qyDomain, $pic, $scid)
    {
        $key           = 'source_map_cache:' . $scid;
        $appIdentifier = 'contact';

        try {
            // 更新员工员工内购部门
            $mediaId = MemberService::getWorkWeChartMediaId($qyDomain, $appIdentifier, $pic);

        } catch (\Exception $exception) {
            Log::error($key . '出错：' . $exception->getMessage());
        }

        if (!empty($mediaId)) {
            Redis::set($key, $mediaId);
            Redis::expire($key, 3 * 24 * 60 * 60);
            //$statisticsMap = Redis::get($key) ?? '';
        }
        return true;
    }
}
