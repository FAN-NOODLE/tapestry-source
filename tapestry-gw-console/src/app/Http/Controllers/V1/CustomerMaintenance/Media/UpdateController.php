<?php
/**
 * 更新素材
 * User: tony
 * Date: 20200220
 * Time: 11:20
 */

namespace App\Http\Controllers\V1\CustomerMaintenance\Media;

use App\Models\QYwxMedia;
use CrsSdk\Work\Source;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Validator;
use App\Http\Controllers\V1\BaseController;

class UpdateController extends BaseController
{
    public function __construct(Request $request)
    {
        parent::__construct();
        $this->request=$request;
    }

    /**
     * 更新素材
     *
     * @param string $scid
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function index($scid="")
    {
        if(!$scid) error(1100002);
        $input = $this->request->only(
            'type', 'name', 'content', 'desc', 'img', 'user_name'
        );

        $rules = [
            'type' => 'required|in:0,1,2,3',
            'img' => 'url'
        ];
        $messages = [
            'type.required'=>'推送类型不能为空',
            'type.in' =>'推送类型数据不在范围内',
            'img' => '图像地址必须是一个合法的url'
        ];

        $validator = Validator::make($input,$rules,$messages);
        if ($validator->fails()) {
            notice(1000010,$validator->errors()->all());
        }

        $data=['type'=>$input['type']];

        if(isset($input['img']) && $input['img']){
            $data['image_media_url']    = trim($input['img']);
        }
        if(isset($input['name']) && $input['name']){
            $data['name']               = trim($input['name']);
        }

        $data['link_desc'] ="";
        $data['link_content'] ="";
        $data['minip_appid'] ="";
        $data['minip_page'] ="";
        $data['text_content'] ="";

        $data['user_name']  =   $this->enterpriseAdmin?$this->enterpriseAdmin['username']:"";
        //0文本素材|1图片素材|2图文链接消息|3小程序卡片素材
        switch ($input['type']){
            case 0:
                $data['text_content']       = trim($input['content']);
                break;
            case 1:
                break;
            case 2:
                $data['link_desc']          = trim($input['desc']);
                $data['link_content']       = trim($input['content']);
                if(isset($input['user_name']))
                    $data['user_name']       = trim($input['user_name']);
                break;
            case 3:
                $data['minip_appid']        = trim($input['desc']);
                $data['minip_page']         = trim($input['content']);
                break;
        }

        $data['user_id']    =   $this->enterpriseAdminId?$this->enterpriseAdminId:"";
        $data['updated_at'] =   time();
        $where=[['scid',$scid]];
        if($this->qyDomain){
            $where[]=['uc_eid',$this->qyDomain];
        }

        try{
            if(!QYwxMedia::query()->where($where)->update($data)){
                error(1100003);
            }
            $id=$scid;
        }catch (\Exception $e){
            Log::info("素材更新失败：[{$e->getCode()}] " . $e->getMessage());
            error(1100004);
        }

        if (!empty($data['image_media_url'])) {
            try {
                $params = [
                    'uc_eid'          => $this->qyDomain ? $this->qyDomain : "",
                    'image_media_url' => $data['image_media_url'],
                    'scid'            => $id,
                ];
                (new Source())->update_media_id($params);
            } catch (\Exception $exception) {
                Log::error('生成mediaId失败params:' . json_encode($params) . '错误提示：' . $exception->getMessage());
            }
        }

        return success(['id' => $id]);
    }
}
