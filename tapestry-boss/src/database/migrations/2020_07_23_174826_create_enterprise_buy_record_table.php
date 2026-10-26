<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateEnterpriseBuyRecordTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('enterprise_buy_record', function (Blueprint $table) {
            $table->unsignedBigInteger('id')->primary()->comment('id');
            $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业id');
            $table->string('order_id')->default('')->comment('订单id');
            $table->unsignedTinyInteger('order_status')->default(0)->comment('订单状态。0-未支付，1-已支付，2-已关闭， 3-未支付且已过期， 4-申请退款中， 5-申请退款成功， 6-退款被拒绝');
            $table->unsignedTinyInteger('order_type')->default(0)->comment('订单类型。0-普通订单，1-扩容订单，2-续期，3-版本变更');
            $table->string('paid_corpid')->default('')->comment('客户企业的corpid');
            $table->string('operator_id')->default('')->comment('下单操作人员userid。如果是服务商代下单，没有该字段。');
            $table->string('suiteid')->default('')->comment('应用id');
            $table->string('appid')->default('')->comment('应用id（仅旧套件有该字段）');
            $table->string('edition_id')->default('')->comment('购买版本ID');
            $table->string('edition_name')->default('')->comment('购买版本名字');
            $table->unsignedBigInteger('price')->default(0)->comment('实付款金额，单位分');
            $table->unsignedInteger('user_count')->default(0)->comment('购买的人数');
            $table->unsignedInteger('order_period')->default(0)->comment('购买的时间，单位天');
            $table->unsignedInteger('order_time')->default(0)->comment('下单时间');
            $table->unsignedInteger('paid_time')->default(0)->comment('付款时间');
            $table->unsignedInteger('begin_time')->default(0)->comment('购买生效期的开始时间');
            $table->unsignedInteger('end_time')->default(0)->comment('购买生效期的结束时间');
            $table->unsignedTinyInteger('order_from')->default(0)->comment('下单来源。0-客户下单；1-服务商代下单；2-代理商代下单');
            $table->string('operator_corpid')->default('')->comment('下单方corpid');
            $table->unsignedBigInteger('service_share_amount')->default(0)->comment('服务商分成金额，单位分');
            $table->unsignedBigInteger('platform_share_amount')->default(0)->comment('平台方分成金额，单位分');
            $table->unsignedBigInteger('dealer_share_amount')->default(0)->comment('代理商分成金额，单位分');
            $table->string('dealer_corp_info')->default('')->comment('渠道商信息（仅当有渠道商报备后才会有此字段）[代理商corpid,代理商名]');
            $table->unsignedInteger('created_at')->default(0);
            $table->unsignedInteger('updated_at')->default(0);
            $table->unsignedInteger('deleted_at')->default(0);
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('enterprise_buy_record');
    }
}
