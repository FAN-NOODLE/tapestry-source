<?php
/**
 * 导购提现认证表
 */
use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateStoreWorkerCertificationTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
//        if (!Schema::hasTable('store_worker_certification')) {
//            Schema::create('store_worker_certification', function (Blueprint $table) {
//                $table->unsignedBigInteger('id', true);
//                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业id')->index();
//                $table->unsignedBigInteger('worker_id')->default(0)->comment('导购ID');
//                $table->string('real_name', 100)->default('')->comment('真实姓名');
//                $table->string('mobile', 20)->default('')->comment('手机号');
//                $table->tinyInteger('id_type')->default(0)
//                    ->comment('证件类型(0：默认无 1:居民身份证, 2:港澳通行证 3：港澳居住证 4：台湾通行证 5：台湾居住证)');
//                $table->string('id_code', 50)->default('')->comment('证件号码');
//                $table->string('bank_card_username', 100)->default('')->comment('银行卡户名');
//                $table->string('bank_card_code', 50)->default('')->comment('银行卡号');
//                $table->string('bank_name', 50)->default('')->comment('银行名称');
//                $table->tinyInteger('agree_deal')->default(0)->comment('同意协议(0:默认, 1:同意)');
//                $table->tinyInteger('auth_status')->default(0)->comment('状态(0-默认 1:认证通过 2:认证失败)');
//                $table->tinyInteger('status')->default(1)->comment('状态(1:启用,2:禁用)');
//                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
//                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
//                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');
//
//                $table->index('mobile');
//                $table->index('id_code');
//            });
//
//            DB::statement("alter table `store_worker_certification comment` COMMENT '导购提现认证信息表' ");
//        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('store_worker_certification');
    }
}
