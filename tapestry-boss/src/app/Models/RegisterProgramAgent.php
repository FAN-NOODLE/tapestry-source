<?php


namespace App\Models;


class RegisterProgramAgent extends Base
{

    // 审核失败
    public const STATUS_APPLY_FAIL = 3;
    // 审核中
    public const STATUS_APPLY = 2 ;
    // 审核成功
    public const STATUS_APPLY_SUCCESS = 1 ;

    //注册错误状态码
    public const REGISTER_CODE_MAP = [
        '89249'  => '该主体已有任务执行中，距上次任务 24h 后再试',
        '89247'  => '内部错误',
        '86004'  => '无效微信号',
        '61070'  => '法人姓名与微信号不一致',
        '89248'  => '企业代码类型无效，请选择正确类型填写',
        '89250'  => '未找到该任务',
        '89251'  => '待法人人脸核身校验',
        '89252'  => '法人&企业信息一致性校验中',
        '89253'  => '缺少参数',
        '89254'  => '第三方权限集不全，请补充权限集后重试',
        '89255'  => 'code参数无效，请检查code长度以及内容是否正确 ；注意code_type的值不同需要传的code长度不一样',
        '99999'  => '其他错误',
        '100001' => '已下发的模板消息法人并未确认且已超时（24h），未进行身份证校验',
        '100002' => '已下发的模板消息法人并未确认且已超时（24h），未进行人脸识别校验',
        '100003' => '已下发的模板消息法人并未确认且已超时（24h）',
        '101'    => '工商数据返回：“企业已注销”',
        '102'    => '工商数据返回：“企业不存在或企业信息未更新”',
        '103'    => '工商数据返回：“企业法定代表人姓名不一致”',
        '104'    => '工商数据返回：“企业法定代表人身份证号码不一致”',
        '105'    => '法定代表人身份证号码，工商数据未更新，请 5-15 个工作日之后尝试',
        '1000'   => '工商数据返回：“企业信息或法定代表人信息不一致”',
        '1001'   => '主体创建小程序数量达到上限',
        '1002'   => '主体违规命中黑名单',
        '1003'   => '管理员绑定账号数量达到上限',
        '1004'   => '管理员违规命中黑名单',
        '1005'   => '管理员手机绑定账号数量达到上限',
        '1006'   => '管理员手机号违规命中黑名单',
        '1007'   => '管理员身份证创建账号数量达到上限',
        '1008'   => '管理员身份证违规命中黑名单',
        '-1'     => '企业与法人姓名不一致',
    ];

    /**
     * 表名称
     *
     * @var string
     */
    protected $table = 'register_program_agent';

    protected $guarded = [];

    function insertData($data)
    {
        self::query()->create($data);
    }

    function getInfoById($id)
    {
        return self::query()->where('id',$id)->get()->first();
    }

    function updateData($id,$update = [])
    {
        return self::query()->where('id',$id)->update($update);
    }

    function getInfoByEnterpriseId($enterprise_id = 0)
    {
        return self::query()->where('enterprise_id',$enterprise_id)->first();
    }

    function getApplyListByStatus($status = self::STATUS_APPLY)
    {
        return self::query()->where('status',$status)->get()->toArray();
    }

    function updateByData($where,$update)
    {
        return self::query()->where($where)->update($update);
    }
}