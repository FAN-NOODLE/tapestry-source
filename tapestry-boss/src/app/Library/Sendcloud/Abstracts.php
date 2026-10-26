<?php


namespace App\Library\Sendcloud;


use App\Exceptions\Error;
use GuzzleHttp\Client;
use GuzzleHttp\Exception\ClientException;

abstract class Abstracts
{
    protected $baseUri = "http://api.sendcloud.net";

    /**
     * GuzzleHttp\Client
     *
     * @var Client|null
     */
    protected $client = null;

    /**
     * 头部信息
     *
     * @var string[]
     */
    protected $header = [
        'Content-Type' => "application/x-www-form-urlencoded"
    ];

    /**
     * 初始化
     *
     * Abstracts constructor.
     */
    public function __construct()
    {
        $this->client = new Client();
    }

    /**
     * @param string $url
     * @param array $data
     * @return mixed
     * @throws \GuzzleHttp\Exception\GuzzleException
     */
    public function post(string $url, array $data = [])
    {
        $options = [
            "form_params" => $data,
        ];
        try {
            $response = $this->client->request("post", $url, $options);
            $data = $response->getBody();
        } catch (ClientException $clientException) {
            throw new Error(1000, $clientException->getMessage(), 500);
        }

        $responseBody = json_decode($data, true);
        if ($responseBody['result'] == false) {
            throw new Error($responseBody['statusCode'], $responseBody['message']);
        }

        return $responseBody['info'];
    }

    /**
     * 设置头部
     *
     * @param array $header
     *
     * @return Abstracts
     */
    public function setHeader($header = []) :Abstracts
    {
        $this->header = array_merge($this->header, $header);
        return $this;
    }
}