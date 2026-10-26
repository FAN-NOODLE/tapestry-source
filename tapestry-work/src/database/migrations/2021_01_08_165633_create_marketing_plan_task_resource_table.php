<?php
/**
 * 营销计划任务素材表
 */
use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;
use Illuminate\Support\Facades\DB;

class CreateMarketingPlanTaskResourceTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        $tableName = 'marketing_plan_task_resource';
        if (!Schema::hasTable($tableName)) {
            Schema::create($tableName, function (Blueprint $table) {
                $table->unsignedBigInteger('id', true)->comment('ID');
                $table->unsignedBigInteger('enterprise_id')->default(0)->comment('企业id')->index();
                $table->unsignedBigInteger('task_id')->default(0)->comment('营销计划任务表ID');
                $table->tinyInteger('type')->default(1)->comment('1、图片(默认)，2、小程序海报 3、视频');
                $table->string('mediaid', 255)->default('')->comment('素材mediaid');
                $table->unsignedInteger('mediaid_updated_at')->default(0)->comment('mediaid上次更新时间');
                $table->string('resource_url', 255)->default('')->comment('素材地址');
                $table->string('applet_path', 255)->default('')->comment('小程序码跳转路径');
                $table->tinyInteger('resource_updated')->default(1)->comment('素材是否更新:1、否  2、是');
                $table->string('applet_card_title', 100)->default('')->comment('小程序卡片标题');
                $table->string('applet_card_cover_url', 255)->default('')->comment('小程序卡片封面地址');
                $table->tinyInteger('is_crontab')->default(1)->comment('是否需要执行定时任务  1默认是  2否');
                $table->unsignedInteger('created_at')->default(0)->comment('创建时间');
                $table->unsignedInteger('updated_at')->default(0)->comment('更新时间');
                $table->unsignedInteger('deleted_at')->default(0)->comment('删除时间');

                $table->index('task_id', 'task_id_index');
            });

            DB::statement("alter table `{$tableName}` COMMENT '营销计划任务素材表' ");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        //Schema::dropIfExists('marketing_plan');
    }
}
