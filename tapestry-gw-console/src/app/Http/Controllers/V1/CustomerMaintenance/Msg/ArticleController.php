<?php
/**
 * 消息推送的富文本内容输出
 * User: tony
 * Date: 20200220
 * Time: 11:20
 */

namespace App\Http\Controllers\V1\CustomerMaintenance\Msg;

use App\Models\QYwxEmsg;
use App\Models\QYwxMedia;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Cache;
use App\Http\Controllers\V1\BaseController;

class ArticleController extends BaseController
{
    protected $requireLogin = false;

    public function __construct(Request $request)
    {
        parent::__construct();
        $this->request=$request;
    }

    /**
     * 外部消息H5内容输出
     * @param string $id {素材id|推送消息id}
     * @return \Illuminate\View\View|void
     */
    public function index($id="")
    {
        if(!$id) error(404,[],404);
        if(strpos($id,"SC") === 0){
            if(!($data = QYwxMedia::query()->where(['scid'=>$id,'type'=>2])->select(['name AS link_title','link_desc','link_content'])->first())){
                error(1000007,[],404);
            }
        }else{
            if(!Cache::has('ArticleQF_'.$id)){
                if($data = QYwxEmsg::query()->where(['tsid'=>$id,'send_type'=>2])->select(['link_title','link_desc','link_content'])->first()){
                    Cache::forever('ArticleQF_'.$id,var_export($data->toArray() ,true));
                }else
                    error(1000007,[],404);
            }else{
                eval('$data = '.Cache::get('ArticleQF_'.$id).";");
            }
            //--添加pv
            QYwxEmsg::query()->where(['tsid'=>$id])->update(["pv"=>DB::raw("pv+1")]);
            //--根据条件添加uv
            if(!isset($_COOKIE['ArticleQF_'.$id])){
                setcookie('ArticleQF_'.$id,true, time()+31536000); //-一年有效期访uv
                QYwxEmsg::query()->where(['tsid'=>$id])->update(["uv"=>DB::raw("uv+1")]);
            }
        }

        if(isset($data['link_content'])){
            return view('article')
                ->with('data', $data);
            //echo '';
        }else error(1000007,[],404);
    }
}
