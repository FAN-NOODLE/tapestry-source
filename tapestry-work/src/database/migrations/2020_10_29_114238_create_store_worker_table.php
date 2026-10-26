<?php

use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateStoreWorkerTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('store_worker')) {
            Schema::create('store_worker', function (Blueprint $table) {
                $table->unsignedBigInteger('id')->primary()->comment('主键ID');
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业ID');
                $table->unsignedTinyInteger('type')->default(1)->comment('导购类型：1、门店导购，2、社群导购，3、企业导购');
                $table->unsignedBigInteger('user_id')->default(0)->comment('用户id');
                $table->unsignedBigInteger('store_id')->default(0)->comment('门店ID');
                $table->string('worker_no', 120)->default('')->comment('员工编号');
                $table->string('name', 120)->default('')->comment('名称');
                $table->string('mobile', 150)->default('')->comment('电话');
                $table->string('mem_uid', 150)->default('')->comment('员工uid');
                $table->unsignedBigInteger('contact_qr_code_id')->default(0)->comment('联系我二维码id');
                $table->string('nickname', 120)->default('')->comment('员工昵称');
                $table->string('avatar')->default('')->comment('头像');
                $table->string('sms_mobile', 20)->default('')->comment('接收短信手机号');
                $table->unsignedTinyInteger('is_send_sms')->default(2)->comment('是否接收提醒短信1:发送 2：不发送');
                $table->string('qr_code')->default('')->comment('企业微信二维码链接');
                $table->string('store_sun_code')->default('')->comment('太阳码链接');
                $table->unsignedTinyInteger('is_boss')->default(0)->comment('是否是店长：1=是, 2=否');
                $table->unsignedTinyInteger('source')->default(0)->comment('来源：1=企业微信');
                $table->unsignedTinyInteger('worker_status')->default(0)->comment('员工状态 1=在职，2=离职');
                $table->unsignedInteger('entry_time')->default(0)->comment('入职时间');
                $table->unsignedInteger('out_time')->default(0)->comment('离职时间');
                $table->unsignedBigInteger('recently_logged_store_id')->default(0)->comment('最近登录门店id');
                $table->unsignedInteger('recently_logged_time')->default(0)->comment('最近登录门店时间');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
            });

            DB::statement("alter table `store_worker` comment '员工表'");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('store_worker');
    }
}
