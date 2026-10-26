<?php

namespace App\Http\Controllers\V1\CustomerMaintenance\Media;

use CrsSdk\Common\QCloudCos;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Validator;
use App\Http\Controllers\V1\BaseController;

class UploadController extends BaseController
{
    /**
     * 不需要登录
     */
    protected $requireLogin = false;
    protected $QcloudCos = null;

    public function __construct(Request $request,QCloudCos $QCloudCos)
    {
        parent::__construct();
        $this->request=$request;
        $this->QcloudCos=$QCloudCos;
    }

    /**
     * @return array 上传文件到腾讯云端
     */
    public function UP()
    {
        $url = $this->QcloudCos->upload($this->request->file);
        return ['url' => $url];
    }

    /**
     * 上传符合企业微信临时素材上传规范的图片素材文件
     * @return \Illuminate\Http\JsonResponse 图片url,name
     */
    public function upMediaIMG(){
        $type=$this->request->input('type','jpeg,jpg,png');
        $validator = Validator::make($this->request->all(),[
            'file' => 'required|image:'.$type
        ], [
            'file.required' => '必须上传图片',
            'file.image'    => '图片格式不正确,只支持'.$type
        ]);
        if ($validator->fails()) {
            notice(1000010,$validator->errors()->all());
        }

        // 3.判断大小是否符合 1
        $tmpFile = $this->request->file->getRealPath();
        if (filesize($tmpFile) >= 2097152) {
            notice(1100001,['2M']);
        }
        $name = $this->request->file->getClientOriginalName();

        $data=$this->UP();
        if(isset($data['url'])){
            $data['name']=$name;
        }
        return success($data);
    }

    /**
     * 上传富文本内容
     * @return \Illuminate\Http\JsonResponse 用于富文本图文内容-图片url
     */
    public function upIMG(){
        $validator = Validator::make($this->request->all(),[
            'file' => 'required|image:jpeg,jpg,png,gif,svg'
        ], [
            'file.required' => '必须上传图片',
            'file.image'    => '图片格式不正确,只支持jpg、png、gif、svg'
        ]);
        if ($validator->fails()) {
            notice(1000010,$validator->errors()->all());
        }

        // 3.判断大小是否符合 1
        $tmpFile = $this->request->file->getRealPath();
        if (filesize($tmpFile) >= 2097152) {
            notice(1100001,['2M']);
        }
        return success($this->UP());
    }
}
