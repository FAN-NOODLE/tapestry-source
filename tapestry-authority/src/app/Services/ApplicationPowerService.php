<?php
/**
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/2/20
 * Time: 17:42
 */
namespace App\Services;

use App\Models\ApplicationPower;

class ApplicationPowerService extends BaseService
{
    /**
     *  服务类型： 商城
     */
    const SERVE_TYPE_SHOP = 1;
    /**
     *  服务类型： crs
     */
    const SERVE_TYPE_CRS = 2;
    /**
     *  服务类型： crs+商城
     */
    const SERVE_TYPE_ALL = 0;
    /**
     * 应用权限表
     * @var ApplicationPower
     */
    protected $applicationPower;

    public function __construct(ApplicationPower $applicationPower)
    {
        $this->applicationPower = $applicationPower;
    }

    /**
     * 同步企业应用权限
     * {service}
     * @param $data
     * @return bool
     */
    public function syncApplicationToConsole($data)
    {
        $res = $this->applicationPower::query()->insert($data);
        return $res;
    }
}
