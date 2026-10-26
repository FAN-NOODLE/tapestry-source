<?php
/**
 * Description:用户自定义标签
 * User: lintong
 * Email: 18511245965@163.com
 * Date: 2020/8/2
 * Time: 16:06
 */

namespace App\Services;


use App\Models\UserCustomizeTags;
use App\Models\UserCustomizeTagsRelation;
use Illuminate\Support\Facades\DB;

class UserCustomizeTagsService extends BaseService
{

    /**
     * 标签名最大长度
     */
    protected const TAG_NAME_MAX_LENGTH = 10;

    /**
     *  用户最大标签数量
     */
    protected const USER_MAX_TAG_NUMBER = 10;

    /**
     * 新增用户自定义标签
     * @param $params
     * @return bool
     */
    public function addNewTags($params)
    {
        $params['tag_name'] = trim($params['name']);

        if (empty($params['tag_name'])) {
            notice(820005);
        }
        if (mb_strlen($params['tag_name']) > self::TAG_NAME_MAX_LENGTH) {
            return ['code' => 820004, 'message' => '标签文字最长10个字'];
        }
        // 查询是否重复
        $datail = $this->getTagByName($params);
        if (!empty($datail)) {
            return ['code' => 820007, 'message' => '该标签已存在'];
        }
        // 添加新的用户自定义标签
        try {
            $result = UserCustomizeTags::query()->create($params);
            return $result->getOriginal('id');
        } catch (\Exception $exception) {
            error(820001);
        }
    }

    /**
     * 根据标签名查询标签信息
     * @param $params
     * @return array
     */
    public function getTagByName($params)
    {
        $detail = UserCustomizeTags::query()
            ->select('id')
            ->where([
                'enterprise_id' => $params['enterprise_id'],
                'worker_id'     => $params['worker_id'],
                'tag_name'      => $params['tag_name']
            ])
            ->first();
        return $detail ? $detail->toArray() : [];
    }

    /**
     * 获取用户绑定的自定义标签信息
     * @param $params
     * @param bool $flag
     * @return array
     */
    public function getUserTagIds($params, $flag = true)
    {
        $tags = UserCustomizeTagsRelation::query()
            ->select(['user_tags_id', 'id', 'enterprise_id', 'worker_id', 'relation_id'])
            ->where([
                'enterprise_id' => $params['enterprise_id'],
                'worker_id'     => $params['worker_id'],
                'relation_id'   => $params['relation_id']
            ])
            ->orderBy('id', 'desc')
            ->get()
            ->toArray();
        if ($flag) {
            return $tags ? array_unique(array_column($tags, 'user_tags_id')) : [];
        } else {
            return $tags ?: [];
        }
    }

    /**
     * 获取用户自定义标签详情
     * @param $ids
     * @return array
     */
    public function getTagsInfoByIds($ids)
    {
        return UserCustomizeTags::query()
            ->select(['id', 'enterprise_id', 'worker_id', 'tag_name'])
            ->whereIn('id', $ids)
            ->get()
            ->toArray();
    }

    /**
     * 绑定新的自定义标签
     * @param $params
     * @return bool
     */
    public function addNewRelation($params)
    {
        $install   = [];
        $tagInsert = [];
        $time      = time();
        $tagNames  = array_filter(array_column($params['tag_names'], 'tag_name'));
        // todo 添加新的关联关系
        $userTagIds = $this->getUserTagIds($params);
        // 根据传过来的标签ID查询标签信息
        $params['names'] = $tagNames;
        $tagInfos        = $this->getTagsByNames($params);
        $newTags         = [];
        foreach ($params['tag_names'] as $value) {
            if (empty($value)) continue;
            if (0 == $value['tag_id']) {
                // 新标签
                $newTags[] = $value['tag_name'];
            }
        }
        // 找新标签并去掉已绑定标签
        $tagIds   = array_column($tagInfos, 'id');
        $tagNames = array_column($tagInfos, 'tag_name');
        $tagIds   = array_diff($tagIds, $userTagIds);
        $newTags  = array_diff($newTags, $tagNames);
        array_map(function ($name) {
            if (mb_strlen(trim($name)) > self::TAG_NAME_MAX_LENGTH) {
                notice(820004);
            }
        }, $newTags);
        // 判断是否超过10个
        if (self::USER_MAX_TAG_NUMBER < (count($userTagIds) + count($newTags) + count($tagIds))) {
            return ['code' => 820003, 'message' => '顾客自定义标签不能超过10个'];
        }
        // 添加新的用户自定义标签
        try {
            DB::beginTransaction();
            // 插入新标签
            if (!empty($newTags)) {
                foreach ($newTags as $key => $value) {
                    if (empty($value)) continue;
                    $tagInsert[$key]['enterprise_id'] = $params['enterprise_id'];
                    $tagInsert[$key]['worker_id']     = $params['worker_id'];
                    $tagInsert[$key]['tag_name']      = $value;
                    $tagInsert[$key]['created_at']    = $time;

                }
                $tagInsert = array_values($tagInsert);
                UserCustomizeTags::query()->insert($tagInsert);
                // 查询新插入的标签ID
                $params['names'] = $newTags;
                $newTagsInfo     = $this->getTagsByNames($params);
                $tagIds          = array_merge($tagIds, array_column($newTagsInfo, 'id'));
            }
            foreach ($tagIds as $key => $value) {
                if (empty($value)) continue;
                $install[$key]['enterprise_id'] = $params['enterprise_id'];
                $install[$key]['worker_id']     = $params['worker_id'];
                $install[$key]['user_tags_id']  = $value;
                $install[$key]['relation_id']   = $params['relation_id'];
                $install[$key]['created_at']    = $time;
            }
            $install = array_values($install);
            UserCustomizeTagsRelation::query()->insert($install);
            DB::commit();
            return true;
        } catch (\Exception $exception) {
            DB::rollBack();
            error(820002);
        }
    }

    /**
     * 获取用户绑定的自定义标签
     * @param $params
     * @return array
     */
    public function getRelation($params)
    {
        $userTagIds = $this->getUserTagIds($params, false);
        if (empty($userTagIds)) return [];
        $tagIds       = array_unique(array_column($userTagIds, 'user_tags_id'));
        $userTagInfos = $this->getTagsInfoByIds($tagIds);
        // 把标签ID当成key
        $userTagInfos = array_combine(array_column($userTagInfos, 'id'), $userTagInfos);
        $result       = [];
        foreach ($userTagIds as $key => $value) {
            $tmp                = [];
            $tmp['relation_id'] = $value['id'];
            $tmp['tag_id']      = $value['user_tags_id'];
            $tmp['tag_name']    = isset($userTagInfos[$value['user_tags_id']]) ? $userTagInfos[$value['user_tags_id']]['tag_name'] : '';
            $result[]           = $tmp;
        }
        return $result;
    }

    /**
     * 获取未绑定的自定义标签列表
     * @param $params
     * @return array
     */
    public function getTags($params)
    {
        $userTagIds = $this->getUserTagIds($params);
        $tagsInfo   = $this->getTagInfo($params, $userTagIds, ['id', 'tag_name']);
        $result     = [];
        foreach ($tagsInfo as $key => $value) {
            if (empty($value)) continue;
            $result[$key]['tag_id']   = $value['id'];
            $result[$key]['tag_name'] = $value['tag_name'];
        }
        return array_values($result);
    }

    /**
     * 获取用户自定义标签信息
     * @param $params
     * @param array $expect
     * @param array $fields
     * @return array
     */
    public function getTagInfo($params, $expect = [], $fields = ['*'])
    {
        $detail = UserCustomizeTags::query()
            ->select($fields)
            ->when((!empty($expect)), function ($query) use ($expect) {
                $query->whereNotIn('id', $expect);
            })
            ->where([
                'enterprise_id' => $params['enterprise_id'],
                'worker_id'     => $params['worker_id'],
            ])
            ->forPage($params['page'], $params['size'])
            ->orderBy('id', 'desc')
            ->get();
        return $detail ? $detail->toArray() : [];
    }

    /**
     * 获取店员所有自定义标签
     * @param $params
     * @return array
     */
    public function getTagsList($params)
    {
        // 获取所有自定义标签
        $tags = UserCustomizeTags::query()
            ->select(['id', 'tag_name'])
            ->where([
                'enterprise_id' => $params['enterprise_id'],
                'worker_id'     => $params['worker_id'],
            ])
            ->forPage($params['page'], $params['size'])
            ->orderBy('id', 'desc')
            ->get()
            ->toArray();
        if (empty($tags)) return [];
        $ids = array_column($tags, 'id');
        // 获取每个标签下绑定的人数
        $relationNum = UserCustomizeTagsRelation::query()
            ->select(DB::raw('count(id) num,user_tags_id'))
            ->whereIn('user_tags_id', $ids)
            ->groupBy('user_tags_id')
            ->get()
            ->toArray();
        $relationNum = array_pluck($relationNum, 'num', 'user_tags_id');
        $result      = [];
        foreach ($tags as $key => $value) {
            $tmp             = [];
            $tmp['tag_id']   = $value['id'];
            $tmp['tag_name'] = $value['tag_name'];
            $tmp['user_num'] = isset($relationNum[$value['id']]) ? $relationNum[$value['id']] : 0;
            $result[]        = $tmp;
        }
        return $result;
    }

    /**
     * 删除自定义标签
     * @param $params
     * @return bool
     */
    public function deleteTag($params)
    {
        $id = $params['id'];
        try {
            DB::beginTransaction();
            // 删除标签
            UserCustomizeTags::query()
                ->where('id', $id)
                ->delete();
            // 解除绑定
            UserCustomizeTagsRelation::query()
                ->where('user_tags_id', $id)
                ->delete();
            DB::commit();
            return true;
        } catch (\Exception $e) {
            DB::rollBack();
            notice(820000, '删除标签失败！');
        }
    }

    /**
     * getRelationIdsByTags
     * @desc 通过标签id获取用户relation_id列表
     * @param $enterpriseId
     * @param $workerId
     * @param $tags
     *
     * @return array
     * @date 2020-08-05
     */
    public static function getRelationIdsByTags($enterpriseId, $workerId, $tags)
    {
        if (empty($tags)) {
            return [];
        }
        $result = UserCustomizeTagsRelation::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('worker_id', $workerId)
            ->whereIn('user_tags_id', $tags)
            ->select(['relation_id'])
            ->get()
            ->toArray();
        return array_column($result, 'relation_id');
    }

    /**
     * getTagOfUsers
     *
     * @desc 获取用户的标签数据
     *
     * @param $enterpriseID
     * @param $workerId
     * @param $relationIdList
     *
     * @return array
     * @date 2020-08-05
     */
    public static function getTagOfUsers($enterpriseID, $workerId, $relationIdList)
    {
        if (empty($relationIdList)) {
            return [];
        }
        return UserCustomizeTagsRelation::query()
            ->leftJoin('user_customize_tags as a', 'a.id', '=', 'user_customize_tags_relation.user_tags_id')
            ->where('user_customize_tags_relation.enterprise_id', $enterpriseID)
            ->where('user_customize_tags_relation.worker_id', $workerId)
            ->whereIn('relation_id', $relationIdList)
            ->orderByDesc('user_customize_tags_relation.id')
            ->select(['user_tags_id', 'relation_id', 'tag_name'])
            ->get()
            ->toArray();
    }

    /**
     * 搜索用户自定义标签
     * @param $params
     * @return array
     */
    public function searchTags($params)
    {
        $detail = UserCustomizeTags::query()
            ->select(['id', 'tag_name'])
            ->where([
                'enterprise_id' => $params['enterprise_id'],
                'worker_id'     => $params['worker_id'],
            ])
            ->where('tag_name', 'like', "{$params['keyword']}" . '%')
            ->orderBy('id', 'desc')
            ->get();
        return $detail ? $detail->toArray() : [];
    }

    /**
     * 通过标签名获取标签信息
     * @param $params
     * @return array
     */
    public function getTagsByNames($params)
    {
        // 根据传过来的标签ID查询标签信息
        $tagInfos = UserCustomizeTags::query()
            ->select(['id', 'tag_name'])
            ->where([
                'enterprise_id' => $params['enterprise_id'],
                'worker_id'     => $params['worker_id'],
            ])
            ->whereIn('tag_name', $params['names'])
            ->get()
            ->toArray();
        return $tagInfos;
    }
}