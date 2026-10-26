<?php

namespace App\Modules\Rbac;

use App\Traits\RecursionTrait;
use App\Models\Permissions as PermissionsModel;

/**
 * 权限
 *
 * @package App\Modules\Rbac
 */
class Permissions
{
    use RecursionTrait;

    private $fields = [
        'id',
        'pid',
        'name',
        'path',
        'key',
        'icon',
        'type',
        'weight',
        'api_list'
    ];

    /**
     * 版本
     *
     * @var int
     */
    private $version;

    /**
     * 版本对应的权限表字段
     *
     * @var string
     */
    private $versionField;

    public function __construct($version)
    {
        if ($version === 0) {
            $this->fields = ['*'];
        } else {
            if (!isset(PermissionsModel::VERSION_RELATION[$version])) {
                throw new \Exception('未知的版本');
            }
            $this->version      = $version;
            $this->versionField = PermissionsModel::VERSION_RELATION[$version];
            $this->fields[] = $this->versionField;
        }
    }

    /**
     * 系统所有权限树形结构
     *
     * @return array
     */
    public function allTree()
    {
        $all = $this->all();

        return $this->getChildrenTree($all);
    }

    /**
     * 权限依赖关系key转成ID
     *
     * @return array
     */
    public function dependRelation()
    {
        $keys = [];
        foreach (PermissionsModel::DEPEND_RELATION as $dependKey => $functionKeys) {
            $keys[] = $dependKey;
            $keys = array_merge($keys, $functionKeys);
        }
        $keys = array_values(array_unique($keys));

        $keyIdRelation = PermissionsModel::query()->whereIn('key', $keys)->pluck('id', 'key')->toArray();

        $data = [];
        foreach (PermissionsModel::DEPEND_RELATION as $dependKey => $functionKeys) {
            if (isset($keyIdRelation[$dependKey]) === false) {
                continue;
            }

            $data[$keyIdRelation[$dependKey]] = [];
            foreach ($functionKeys as $key) {
                if (isset($keyIdRelation[$key]) === false) {
                    continue;
                }
                $data[$keyIdRelation[$dependKey]][] = $keyIdRelation[$key];
            }
        }

        return $data;
    }

    /**
     * 标配权限列表
     *
     * @param array $versionTagList
     * @param bool $needIsFullRelease 是否需要全量发布的
     * @return array
     */
    public function getStandardConfig($versionTagList = [], $needIsFullRelease = true)
    {
        $where = [
            [$this->versionField, '=', PermissionsModel::VERSION_STANDARD_CONFIG],
            ['version_tag', 'in', $versionTagList]
        ];
        if ($needIsFullRelease === true) {
            $where[] = ['is_full_release', '=', 1];
        }
        return $this->all($where);
    }

    /**
     * 选配功能列表
     *
     * @return array
     * @throws \Exception
     */
    public function additionalConfigList()
    {
        return $this->all([
            ['type', 'in', [PermissionsModel::TYPE_MENU, PermissionsModel::TYPE_MENU_AND_VIEW]],
            ['is_full_release', '=', 1],
            [$this->versionField, '=', PermissionsModel::VERSION_SELECT_CONFIG]
        ]);
    }

    /**
     * 无功能列表
     *
     * @return array
     * @throws \Exception
     */
    public function noConfigList()
    {
        return $this->all([
            ['type', '=', PermissionsModel::TYPE_MENU],
            ['is_full_release', '=', 1],
            [$this->versionField, '=', PermissionsModel::VERSION_NO_CONFIG]
        ]);
    }

    /**
     * 所有权限列表
     *
     * @param array $where 二维数组
     * @return array
     */
    public function all($where = [])
    {
        return PermissionsModel::query()
            ->select($this->fields)
            ->whereIn('system_type', array_keys(PermissionsModel::CONSOLE_SYSTEM_TYPE_RELATION))
            ->when($where !== [], function ($query) use ($where) {
                foreach ($where as $key => $value) {
                    if ($value[1] === 'in') {
                        if ($value[2] === []) {
                            continue;
                        }
                        $query->whereIn($value[0], $value[2]);
                    } else {
                        $query->where($value[0], $value[2]);
                    }
                }
            })->orderBy('weight', 'DESC')->get()->toArray();
    }

    /**
     * 获取下面所有的权限ID
     *
     * @return array
     */
    public function getChildrenIDList($id)
    {
        $all = $this->all();
        $children = $this->getChildren($all, $id);

        return array_column($children, 'id');
    }

    public function getVersion()
    {
        return $this->version;
    }

    public function getVersionField()
    {
        return $this->versionField;
    }
}
