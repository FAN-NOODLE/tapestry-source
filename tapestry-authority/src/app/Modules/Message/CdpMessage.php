<?php

namespace App\Modules\Message;

use App\Modules\Message\Content\CdpMessageContent;
use GuzzleHttp\Client;
use GuzzleHttp\Exception\GuzzleException;
use Illuminate\Support\Facades\App;
use Illuminate\Support\Facades\Log;

class CdpMessage extends AbstractMessage
{
    /**
     * 消息内容与方法关系
     *
     * @var array
     */
    protected $messageContentRelation = [
        '1000' => 'helpFail',
        '1001' => 'helpOverPayExpressFee',
        '1002' => 'helpFail',
        '1003' => 'helpOver',
        '1004' => 'helpStart',
        '1005' => 'helpPeopleAdd',
        '1006' => 'luckResult',
        '1007' => 'groupPeopleAdd',
        '1008' => 'groupFail',
        '1009' => 'hagglePeopleAdd',
        '1010' => 'haggleResult',
        '1011' => 'liveStart',
        '2000' => 'orderWaitPay',
        '2001' => 'orderExpress',
        '2002' => 'orderRefund',
        '3000' => 'scoreChange',
    ];

    /**
     * 第三方平台基本参数
     *
     * @var array
     */
    private $configs = [
        'prod' => [
            'domain' => 'https://wx-gzh-prod.vchangyi.com',
            'key'    => 'CDP-GZH-PROD-21269932-ME'
        ],
        'test' => [
            'domain' => 'http://wx-gzh-test.vchangyi.com',
            'key'    => '123456'
        ]
    ];
    private $apis = [
        'tmSend' => '/app/weixin/server/tm_send'
    ];
    private $config;

    private $user;

    public function __construct()
    {
        $this->config = $this->configs['test'];

        if ('prod' == App::environment() && 0 == env('WEIXIN_BETA', 1)) {
            $this->config = $this->configs['prod'];
        }
    }

    /**
     * @return array|mixed
     */
    public function send()
    {
        $result = [
            'success_list' => [],
            'fail_list'    => []
        ];

        $client = new Client(['verify' => false]);
        foreach ($this->users as $key => $user) {
            if (empty($user['union_id'])) {
                $result['fail_list'][] = [
                    'user'     => $user['id'],
                    'response' => '用户union_id不存在'
                ];
                continue;
            }

            $this->user = $user;
            $content    = $this->getContent();

            $url = $this->config['domain'] . $this->apis['tmSend'];
            try {
                $response = $client->request('POST', $url, [
                    'http_errors' => false,
                    'headers'     => [
                        'Accept'       => 'application/json',
                        'charset'      => 'utf-8',
                        'Content-Type' => 'application/json',
                        'Sign'         => $this->getSign($content)
                    ],
                    'body'        => json_encode($content, JSON_UNESCAPED_UNICODE + JSON_UNESCAPED_SLASHES)
                ]);

                Log::info('公众号推送消息-给cdp', $content);

                $result['success_list'][] = [
                    'user'     => $user['id'],
                    'response' => json_decode($response->getBody()->getContents(), true)
                ];
            } catch (GuzzleException $exception) {
                Log::error('公众号推送消息-给cdp失败', [
                    'message' => $exception->getMessage(),
                    'file'    => $exception->getFile(),
                    'line'    => $exception->getLine()
                ]);

                $result['fail_list'][] = [
                    'user'     => $user['id'],
                    'response' => $exception->getMessage()
                ];
            }
        }

        return $result;
    }

    /**
     * @return array
     */
    protected function getContent(): array
    {
        $messageContent = new CdpMessageContent($this->message);

        return [
            'enterprise_id' => $this->enterpriseId,
            'type'          => $this->sceneType,
            'page_query'    => $this->option['page_query'] ?? '',
            'unionid'       => $this->user['union_id'],
            'cdp_codes'     => $this->user['cdp_codes'] ?? [],
            'keys'          => $messageContent->{$this->messageContentRelation[(string)$this->sceneType]}()
        ];
    }

    /**
     * @param $data
     * @return string
     */
    private function getSign($data)
    {
        ksort($data);           //--排序
        $data = serialize($data); //--序列化
        $data .= '_key=' . $this->config['key']; //--拼接key，测试系统key=123456，生产环境待约定

        return strtoupper(substr(md5($data), 6, 8)); //--截断编码最大化
    }
}
