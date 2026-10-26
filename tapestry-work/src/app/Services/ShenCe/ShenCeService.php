<?php


namespace App\Services;

namespace App\Services\ShenCe;

use App\Sdk\Active\CrmHelper;
use App\Sdk\ShenCe\Client;
use Exception;
use SensorsAnalytics;

class ShenCeService
{
    /**
     * @var SensorsAnalytics
     */
    private $client;

    private $appName = [
        'coach' => 'Coach_Wecom',
        'sw'    => 'SW_Wecom',
        'ks'    => 'KS_Wecom',
    ];

    /**
     * @param $enterpriseId
     * @throws Exception
     */
    public function __construct($enterpriseId, string $consumerType = Client::DebugConsumer)
    {
        $enterpriseConfig = CrmHelper::getEnterpriseConfig($enterpriseId);
        $serverUrl        = env('SHEN_CE_DOMAIN');
        switch ($enterpriseConfig['name']) {
            case "coach":
                $projectName = env('SHEN_CE_COACH_PROJECT');
                break;
            case "sw":
                $projectName = env('SHEN_CE_SW_PROJECT');
                break;
            case "ks":
                $projectName = env('SHEN_CE_KS_PROJECT');
                break;
            default:
                throw new Exception("Unknown enterprise");
        }

        $sdk = new Client($consumerType, $serverUrl, $projectName);

        $instance = $sdk->instance();
        $instance->register_super_properties(['app_name' => $this->appName[$enterpriseConfig['name']]]);

        $this->client = $instance;
    }


    /**
     * @return SensorsAnalytics
     */
    public function getClient(): SensorsAnalytics
    {
        return $this->client;
    }

    /**
     * 调用结束后发送数据
     */
    public function __destruct()
    {
        $this->client->flush();
        $this->client->close();
    }
}