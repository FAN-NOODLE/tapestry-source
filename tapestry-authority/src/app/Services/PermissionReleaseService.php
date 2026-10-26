<?php

namespace App\Services;

use App\Models\Enterprise;
use App\Modules\Rbac\Role;
use Illuminate\Support\Facades\Log;

class PermissionReleaseService
{
    const ALL_ENTERPRISE = 0;

    private $systemType;
    private $versionTag;
    private $enterpriseID;

    /**
     * 企业管理员关系列表
     *
     * @var array
     */
    public $enterpriseAdminRelationList;

    public function __construct(int $systemType, string $versionTag, int $enterpriseID)
    {
        $this->systemType = $systemType;
        $this->versionTag = $versionTag;
        $this->enterpriseID = $enterpriseID;
    }

    public function enterpriseAdminRelationList(): array
    {
        // 查出需要处理的企业ID与超管ID
        $this->enterpriseAdminRelationList = Enterprise::query()
            ->select([
                'enterprise.id AS enterprise_id',
                'enterprise_admin.id AS admin_id'
            ])
            ->leftJoin('enterprise_admin', 'enterprise.id', '=', 'enterprise_admin.enterprise_id')
            ->where('enterprise.version', '>', 1)
            ->where('enterprise_admin.is_super', '=', 1)
            ->when($this->enterpriseID !== self::ALL_ENTERPRISE, function ($query) {
                return $query->where('enterprise.id', '=', $this->enterpriseID);
            })
            ->orderBy('enterprise.id', 'ASC')
            ->get()
            ->toArray();

        return $this->enterpriseAdminRelationList;
    }

    public function release(): string
    {
        $result = '已发布' . count($this->enterpriseAdminRelationList) . '个企业';
        foreach ($this->enterpriseAdminRelationList as $enterpriseAdmin) {
            try {
                $role = new Role($enterpriseAdmin['admin_id'], $enterpriseAdmin['enterprise_id']);
                $res = $role->permissions->superAdminRelease($this->systemType, $this->versionTag);
                if ($res === '没有权限需要处理') {
                    // 如果进到这里说明找不到权限数据，无需再遍历执行了over
                    $result = $res;
                    break;
                }
            } catch (\Exception $exception) {
                Log::error('发布超管权限异常', [
                    'message' => $exception->getMessage(),
                    'file'    => $exception->getFile(),
                    'line'    => $exception->getLine(),
                    'args'    => [
                        'system_type' => $this->systemType,
                        'version_tag' => $this->versionTag
                    ]
                ]);
            }
        }

        return $result;
    }
}
