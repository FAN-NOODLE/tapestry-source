<?php

namespace App\Traits;

/**
 * 递归工具
 *
 * @package App\Traits
 */
trait RecursionTrait
{
    /**
     * 获取指定元素下的所有子级
     *
     * @param array $data
     * @param int $id
     * @param int $level
     * @return array 一维
     */
    public function getChildren($data, $id = 0, $level = 0)
    {
        $arr = [];
        foreach ($data as $key => $value) {
            if ($value['pid'] == $id) {
                $value['level'] = $level + 1;
                $arr[] = $value;
                $arr = array_merge($arr, $this->getChildren($data, $value['id'], $value['level']));
            }
        }
        return $arr;
    }

    /**
     * 获取指定元素下的所有子级树
     *
     * @param array $data
     * @param int $id
     * @param int $level
     * @return array 树形结构
     */
    public function getChildrenTree($data, $id = 0, $level = 0)
    {
        $tree = [];
        foreach ($data as $key => $value) {
            if ($value['pid'] === $id) {
                $value['level'] = $level + 1;
                $value['children'] = $this->getChildrenTree($data, $value['id'], $value['level']);
                $tree[] = $value;
                unset($data[$key]);
            }
        }

        return $tree;
    }

    /**
     * 获取所有的父级
     *
     * @param array $data
     * @param $pid
     * @return array 一维
     */
    public function getParent($data, $pid)
    {
        static $parent = [];

        foreach ($data as $key => $value) {
            if ($value['id'] === $pid) {
                $parent[] = $value;

                $this->getParent($data, $value['pid']);
            }
        }

        return $parent;
    }
}
