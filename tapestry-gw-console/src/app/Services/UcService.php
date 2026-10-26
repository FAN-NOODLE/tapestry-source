<?php
/**
 * Uc接口
 * User: 王鑫阳
 * Date: 2020-02-23
 * Time: 10:34
 */
namespace App\Services;

class UcService extends BaseService
{
    /**
     * @notes 域名
     */
    protected $host = [
        'http://t-rst.vchangyi.com',//uc
        'http://cdp-servicetest.vchangyi.com',//转uc 李旭接口
    ];
    /**
     * @notes 企业标示
     * @enterprise [string] path
     */
    protected $enterprise=null;

    /**
     * @notes 接口地址映射
     * @标示[地址] [array]
     */
    protected $urlMap = [
        'groupList'=>[0,'/b/{enterprise}/qy/external/group/list'],                      //群列表
        'groupDetail'=>[0,'/b/{enterprise}/qy/external/group/detail'],                  //群详情
        'groupMember'=>[0,'/b/{enterprise}/qy/external/group/member/list'],             //群成员列表
        'userExternal'=>[0,'/b/{enterprise}/qy/external/detail'],                       //外部联系人详情
        'userDetail'=>[0,'/a/{enterprise}/contact/qy/member/get'],                      //企业人员详情
        'userList'=>[0,'/a/{enterprise}/contact/qy/member/ignore-permissions/list'],    //-用户列表
        'tagList'=>[0,'/b/{enterprise}/qy/external/enterprise-tag/list'],               //-外部联系人标签
        'numEst'=>[1,'/qywx/notify/calPushCnt/{enterprise}'],                           //预计触达
        'send'=>[1,'/qywx/notify/startPush/{enterprise}'],                              //推送接口
        'sendRet'=>[1,'/qywx/notify/getPushDetail/{enterprise}'],                                     //推送数据结果接口
    ];

    public function __construct($domain)
    {
        if(env('APP_ENV')=='test'){
            $this->host = [
                'http://t-rst.vchangyi.com',
                'http://cdp-servicetest.vchangyi.com'
            ];
        }
        if(env('APP_ENV')=='prod'){
            $this->host = [
                'http://restapi.99hr.com',
                'cdp-service.vchangyi.com'
            ];
        }
        $this->enterprise = $domain;
    }

    protected function getQueryUrl($urlKey)
    {
        $url = $this->host[$this->urlMap[$urlKey][0]].str_replace('{enterprise}',$this->enterprise,$this->urlMap[$urlKey][1]);
        return $url;
    }

    private function toData(&$arr){
        $post=[];
        foreach ($arr AS $key=>$value){
            array_push($post,"{$key}=".urlencode($value));
        }
        $arr = join('&',$post);
    }

    public function send($urlKey,$data="{}",$method = 'POST',$type='json',$header=[])
    {
        $url = $this->getQueryUrl($urlKey);

        //初始化
        $ch = curl_init();
        $headers = [
            'form-data' => ['Content-Type: multipart/form-data'],
            'json'      => ['Content-Type: application/json'],
        ];

        $Header=[];
        if(!$data) $method='GET';
        if($method == 'GET'){
            if($data){
                if(is_array($data)||is_object($data)) $this->toData($data);
                $querystring = http_build_query($data);
                $url = $url.(strpos($url,'?') !== false?'&':'?').$querystring;
            }
        }else{
            if(is_array($data)||is_object($data)){
                $type='form-data';
                $this->toData($data);
            }
            isset($headers[$type]) && $Header=$headers[$type];
        }

        curl_setopt($ch, CURLOPT_URL,$url);
        curl_setopt($ch, CURLOPT_HEADER, false);
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);         // 执行后不直接打印出来
        if($method == 'POST'){
            curl_setopt($ch, CURLOPT_CUSTOMREQUEST,'POST');     // 请求方式
            curl_setopt($ch, CURLOPT_POST, true);               // post提交
            curl_setopt($ch, CURLOPT_POSTFIELDS, $data);              // post的变量
        }
        if($method == 'PUT'){
            curl_setopt ($ch, CURLOPT_CUSTOMREQUEST, "PUT");
            curl_setopt($ch, CURLOPT_POSTFIELDS,$data);
        }
        if($method == 'DELETE'){
            curl_setopt ($ch, CURLOPT_CUSTOMREQUEST, "DELETE");
            curl_setopt($ch, CURLOPT_POSTFIELDS,$data);
        }
        if(stripos($url,"https://")!==FALSE){
            curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, false); // 跳过证书检查
            curl_setopt($ch, CURLOPT_SSL_VERIFYHOST, false); // 不从证书中检查SSL加密算法是否存在
        }
        curl_setopt($ch, CURLOPT_HTTPHEADER, array_merge($Header,$header));
        $output = curl_exec($ch); //执行并获取HTML文档内容
        $aStatus = curl_getinfo($ch);
        curl_close($ch);
        if(intval($aStatus["http_code"])==200){
            return $output;
        }else{
            error_log("UcService-ERR:[{$aStatus["http_code"]}] {$urlKey} - (".$data.")");
            return false;
        }
    }
}
