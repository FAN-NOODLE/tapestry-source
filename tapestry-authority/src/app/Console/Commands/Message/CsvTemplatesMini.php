<?php

namespace App\Console\Commands\Message;

use App\Models\EnterpriseThirdPartAccount;
use Illuminate\Console\Command;

class CsvTemplatesMini extends Command
{

    protected $signature = 'csv_template_mini';

    protected $description = '导出消息到文件';

    /**
     * 导出小程序类目
     *
     * @throws \GuzzleHttp\Exception\GuzzleException
     * @throws \Psr\SimpleCache\InvalidArgumentException
     */
    public function handle()
    {
        $enterpriseCategories = EnterpriseThirdPartAccount::query()
            ->where('wx_category_ids','<>','')
            ->pluck('wx_category','name')
            ->toArray();

        $data = [];
        foreach ($enterpriseCategories as $key=>$enterpriseCategory){
            $arr = json_decode($enterpriseCategory,1);
            array_walk($arr, function (&$value) use ($key) {
                $value['mini_name']=$key;
            });

            $data = array_merge($arr,$data);
        }

        $data = collect($data)->sortBy('id',1)->toArray();

        $header = ['类目ID', '类目名称', '小程序'];

        exportToCsv('小程序类目', $header, $data);
    }
}
