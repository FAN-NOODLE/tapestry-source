<?php

namespace App\Services;

use App\Models\Announcement;

class AnnouncementService
{
    /**
     * 公告列表
     *
     * @param $page
     * @param $size
     * @param array $search
     * @param array $fields
     * @return array
     */
    public function lists($page, $size, $search = [], $fields = ['id', 'title', 'type', 'released_at', 'created_at'], $notPage = false)
    {
        $orm = Announcement::query()->select($fields)->when($search !== [], function ($query) use ($search) {
            if (isset($search['title']) && $search['title'] !== '') {
                $query->where('title', 'LIKE', "%{$search['title']}%");
            }
            if (isset($search['type']) && $search['type'] !== '') {
                $query->where('type', '=', $search['type']);
            }
            if (isset($search['time_type']) && in_array($search['time_type'], ['released_at', 'created_at'])) {
                if (isset($search['start_at']) && isset($search['end_at']) && $search['end_at'] > 0) {
                    $query->whereBetween($search['time_type'], [
                        $search['start_at'],
                        $search['end_at']
                    ]);
                }
            }
        });

        $total = $orm->count('id');

        // $notPage为false表示要分页
        $orm = $orm->orderBy('released_at', 'DESC');
        if ($notPage === false) {
            $orm = $orm->forPage($page, $size);
        }

        $list = $orm->get()
            ->toArray();

        return [$total, $list];
    }

    /**
     * 公告详情
     *
     * @param $id
     * @return Announcement|\Illuminate\Database\Eloquent\Builder|\Illuminate\Database\Eloquent\Model|object|null
     */
    public function detail($id)
    {
        return Announcement::query()->where('id', '=', $id)->first();
    }
}
