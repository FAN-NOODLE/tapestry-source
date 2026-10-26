<?php


namespace App\Services;


use App\Common\Constant;
use App\Library\Uc\Tag;
use App\Models\CustomerGroup;
use App\Models\TagCustomerGroupRelation;
use Illuminate\Support\Facades\Log;

class TagService
{
    /**
     * 标签组列表获取记录数
     */
    const TAG_GROUP_PAGE_SIZE = 100;

    /**
     * 添加标签
     *
     * @param $input
     * @return mixed
     * @throws \Exception
     */
    public function addTag($input)
    {
        // 请求uc添加标签
        $params = [
            'groupId'      => $input['tag_group_id'],
            'groupName'    => $input['tag_group_name'],
            'order'        => $input['order'] ?? 0,
            'tagList'      => $input['tag_list'],
            "thirdGroupId" => ""
        ];

        if (!isset($input['tag_group_id']) || empty($input['tag_group_id'])) {
            unset($params['groupId']);
        }

        return (new Tag())->add($input['qy_domain'], $params);
    }

    /**
     * 更新标签或标签组
     * @param $input
     * @return mixed
     * @throws \Exception
     */
    public function editTag($input)
    {
        $params = [
            'id'      => $input['id'],
            'name'    => $input['name'],
            "order"   => 0,
            "thirdId" => ""
        ];

        return (new Tag())->update($input['qy_domain'], $params);
    }

    /**
     * 删除tag
     * @param $input
     * @return mixed
     * @throws \Exception
     */
    public function delTag($input)
    {
        $params = [];

        if (isset($input['tag_id'])) {
            $params = [
                'tagIds' => [(int)$input['tag_id']]
            ];
        }

        if (isset($input['tag_group_id'])) {
            $params = [
                'groupIds' => [(int)$input['tag_group_id']]
            ];
        }

        if (empty($params)) {
            notice(810008);
        }

        try {
            (new Tag())->delete($input['qy_domain'], $params);
        } catch (\Exception $exception) {
            Log::info('删除标签失败');
            notice(810013);
        }

        // 删除成功后 需要删除标签组下标签或者标签对应的标签关系
        if (isset($input['tag_id'])) {
            return $this->delTagCustomerGroupRelation([$input['tag_id']]);
        } elseif (isset($input['tag_group_id'])) {
            // 删除该分组下所有关系
            return TagCustomerGroupRelation::query()
                ->where('tag_group_id', $input['tag_group_id'])
                ->delete();
        }

        return true;
    }

    /**
     * 获取所有标签列表
     * @param $domain
     * @return mixed
     * @throws \Exception
     */
    public function getAllTagList($domain)
    {
        return  (new Tag())->allList($domain);
    }

    /**
     * 获取企业标签库标签列表
     * @param $domain
     * @return mixed
     * @throws \Exception
     */
    public function tagListForSelect($input)
    {
        $params = [
            'pageNum' => $input['page'] ?? Constant::PAGE,
            'pageSize' => $input['size'] ?? Constant::SIZE
        ];
        return (new Tag())->pageList($input['qy_domain'], $params);
    }

    /**
     * 返回企业标签组列表
     * @param $domain
     * @return mixed
     * @throws \Exception
     */
    public function listTag($domain)
    {
        $params = [
            'pageSize' => self::TAG_GROUP_PAGE_SIZE
        ];

        // 处理数据适合返回
        $tagGroupList = (new Tag())->list($domain, $params);

        $backGroupList = [];
        foreach ($tagGroupList['list'] as $item) {
            if($item['tagCount'] == 0 ) {
                continue;
            }
            $backGroupList[] = [
                'tag_group_id'   => (string)$item['id'],
                'tag_group_name' => $item['name'],
                'tag_count'      => $item['tagCount'] ?? 0
            ];
        }

        return $backGroupList;
    }

    /**
     * 根据分组id返回起亚标签库
     * @param $domain
     * @param $input
     * @return mixed
     * @throws \Exception
     */
    public function pageListTag($domain, $input)
    {

        $params = [
            'groupId'  => (int)$input['tag_group_id'],
            'pageNum'  => $input['page'],
            'pageSize' => $input['size']
        ];

        if (isset($input['keyword']) && !empty($input['keyword'])) {
            // 标签名称搜索
            $params['name'] = trim($input['keyword']);
        }

        if (isset($input['custom_group_id']) && !empty($input['custom_group_id'])) {
            // 用户分组Id
            $tagIds = $this->getTagIdsByCustomGroupId($input['enterprise_id'], $input['tag_group_id'], $input['custom_group_id']);
            if (!empty($tagIds)) {
                $params['ids'] = $tagIds;
            } else {
                return [];
            }
        }

        $tagList = (new Tag())->pageList($domain, $params);

        $ucTagIds = array_column($tagList['list'], 'id');
        if (empty($ucTagIds)) {
            return [];
        }

        // 获取tag 关联 的用户分组
        $groupList = $this->getGroupIdByTagIds($ucTagIds);
        $groupList = array_column($groupList, null, 'tag_id');

        // 获取分组名称
        $cGroupStringIds = array_column($groupList, 'custom_group_ids');
        $cGroupIds = [];

        foreach ($cGroupStringIds as $cg) {
            $cGroupIds = array_merge($cGroupIds, json_decode($cg, true));
        }

        $cGroupIds = array_values(array_unique($cGroupIds));

        $customerGroupList = [];
        if (!empty($cGroupIds)) {
            $customerGroupList = CustomerGroup::query()
                ->whereIn('id', $cGroupIds)
                ->get()
                ->toArray();

            $customerGroupList = array_column($customerGroupList, null, 'id');
        }


        $backList = [];
        foreach ($tagList['list'] as $k=>$item) {

            $backList['list'][$k]['tag_id'] = (string)$item['id'];
            $backList['list'][$k]['name'] = $item['name'];
            $backList['list'][$k]['tag_group_id'] = (string)$item['groupId'];
            $backList['list'][$k]['created_at'] = $item['createTime'] ?? 0;
            $backList['list'][$k]['custom_group_ids'] = !empty($groupList[$item['id']]['custom_group_ids']) ? json_decode($groupList[$item['id']]['custom_group_ids'], true) : [];
            $backList['list'][$k]['custom_group_names'] = '';
            $groupNameArr = [];
            // 如果用户分组中存在则处理数据
            if (!empty($groupList[$item['id']]['custom_group_ids'])) {
                foreach ($customerGroupList as $v) {
                    if (in_array($v['id'], json_decode($groupList[$item['id']]['custom_group_ids']))) {
                        $groupNameArr[] = $v['name'];
                    }
                }
                $backList['list'][$k]['custom_group_names'] = implode(',', $groupNameArr);
            }

        }

        $backList['page'] = $tagList['pageNum'];
        $backList['size'] = $tagList['pageSize'];
        $backList['total'] = $tagList['total'];

        return $backList;
    }

    /**
     * 插入标签和用户分组关系
     * @param $input
     * @return int|mixed
     */
    public function addTagCustomerGroupRelation($input)
    {
        $input['custom_group_ids'] = json_encode($input['custom_group_ids']);

        $data = $this->detailTagCustomerGroupRelation($input['tag_id']);
        // 已经存在数据更新
        if (!empty($data)) {
            //$this->delTagCustomerGroupRelation($input['tag_id']);
            $this->updateTagCustomerGroupRelation($input['tag_id'], $input['custom_group_ids']);
            return $data['id'];
        }

        if (empty($input['custom_group_ids'])) {
            notice(810012);
        }

        $tagRelation = TagCustomerGroupRelation::query()->create($input);
        return $tagRelation ? $tagRelation['id'] : 0;
    }

    private function updateTagCustomerGroupRelation($tagId, $update)
    {
        return TagCustomerGroupRelation::query()->where([
            ['tag_id', $tagId],
        ])->update(['custom_group_ids' => $update]);
    }

    /**
     * 删除标签和用户分组关联关系
     * @param $tagIds
     * @return \Illuminate\Database\Eloquent\Builder|\Illuminate\Database\Eloquent\Model|object|null
     */
    private function delTagCustomerGroupRelation($tagIds)
    {

        return TagCustomerGroupRelation::query()->whereIn('tag_id', $tagIds)->delete();

    }

    /**
     * 获取标签和用户分组关联详情
     * @param $tagId
     * @return array
     */
    private function detailTagCustomerGroupRelation($tagId)
    {
        $detail = TagCustomerGroupRelation::query()->where([
            ['tag_id', $tagId],
        ])->first();

        return $detail ? $detail->toArray() : [];
    }

    /**
     * 根据tagIds获取groupId
     * @param $tagIds
     * @return array
     */
    public function getGroupIdByTagIds($tagIds)
    {
        $list = TagCustomerGroupRelation::query()
            ->select(['tag_id', 'custom_group_ids'])
            ->whereIn('tag_id', $tagIds)
            ->get();

        return  $list ?  $list->toArray() : [];
    }

    /**l
     * 根据用户分组Id获取标签id集合
     * @param $tagGroupId
     * @param $customGroupId
     * @return array
     */
    public function getTagIdsByCustomGroupId($enterpriseId, $tagGroupId, $customGroupId)
    {
        $tagList = TagCustomerGroupRelation::query()
            ->where([
                ['tag_group_id', $tagGroupId],
                ['enterprise_id', $enterpriseId]
            ])
            ->get()
            ->toArray();

        if (!$tagList) {
            return [];
        }

        $tagIds = [];
        foreach ($tagList as $value) {
            $customGroupIds = json_decode($value['custom_group_ids'], true);
            if (in_array($customGroupId, $customGroupIds)) {
                $tagIds[] = $value['tag_id'];
            }
        }

        return $tagIds;
    }

    /**
     * 判断标签名称是否重复
     * @param $qyDomain
     * @param $groupId
     * @param $tagList
     * @throws \Exception
     */
    public function isRepeatTag($qyDomain, $groupId, $tagList, $type = 'add')
    {
        // 获取分组下所有标签
        $params = [
            'groupId'  => $groupId,
            'pageNum'  => 1,
            'pageSize' => 1000
        ];

        $list = (new Tag())->pageList($qyDomain, $params);

        if ($type == 'add') {

            $listName = array_column($list['list'], 'name');
            $requestListName = array_column($tagList, 'name');
            $ins = array_intersect($listName, $requestListName);
            if (!empty($ins)) {
                notice(810011);
            }

        } elseif ($type == 'edit') {

            //编辑判断需要去掉自身
            $list = array_column($list['list'] , null,'id');
            unset($list[$tagList['id']]);
            $listName = array_column($list, 'name');
            if (in_array($tagList['name'], $listName)) {
                notice(810011);
            }

        }

    }

    /**
     * getUserOfTagList
     * @desc 调用UC接口获取设置标签的所有用户
     * @param $domain
     * @param $tagList array tagid列表
     *
     * @return array
     * @date 2020-07-15
     * @author XuHuitao
     * @throws \Exception
     */
    public function getUserOfTagList($domain, $tagList) {
        if (empty($tagList)) {
            return [];
        }
        $tag = new Tag();
        $result =  $tag->tagDetail($domain, ['tagIds' => $tagList]);
        if (isset($result['requestId'])) {
            unset($result['requestId']);
        }
        return $result;
    }


}