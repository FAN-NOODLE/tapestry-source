<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/4/11 0011
 * Time: 11:32
 */

namespace App\Http\Controllers\App\V1\Order;


use App\Http\Controllers\App\V1\BaseController as Controller;
use Illuminate\Support\Facades\Auth;


class OrderBaseController extends Controller
{
    use \App\Traits\OrderTrait;

    /**
     * 登录人信息
     *
     * @var array
     */
    protected $user = [];

    /**
     * 构造方法 校验店长信息并加入类属性
     *
     * OrderBaseController constructor.
     */
    public function __construct()
    {
        parent:: __construct();
        $this->user = Auth::user();
    }


    /**
     * 获取商品对数量
     * @param array $goodsList
     *
     * @return int
     */
    public function getGoodsNum($goodsList = [])
    {
        $num = 0;
        if (!empty($goodsList)) {
            $num = array_sum(array_column($goodsList, 'quantity'));
        }

        return $num;
    }
}
