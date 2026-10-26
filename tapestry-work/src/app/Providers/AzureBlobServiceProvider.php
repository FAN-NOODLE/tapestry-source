<?php
/**
 * azure blob 服务提供者
 * User: wangle@vchangyi.com
 * Date: 2021/03/08
 */

namespace App\Providers;


use Illuminate\Support\Facades\Storage;
use Illuminate\Support\ServiceProvider;
use League\Flysystem\AzureBlobStorage\AzureBlobStorageAdapter;
use League\Flysystem\Filesystem;
use MicrosoftAzure\Storage\Blob\BlobRestProxy;

class AzureBlobServiceProvider extends ServiceProvider
{

    /**
     * Register the application service
     *
     * @return void
     */
    public function register()
    {


    }

    /**
     * 启动应用服务
     * @return void
     */
    public function boot()
    {
        Storage::extend('azure', function ($app, $config) {
            $endpoint = sprintf('DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=core.chinacloudapi.cn',
                $config['name'], $config['key']);
            $client = BlobRestProxy::createBlobService($endpoint);
            $adapter = new AzureBlobStorageAdapter($client, $config['container'], $config['prefix'] ?? null);
            return new Filesystem($adapter);
        });
    }


}