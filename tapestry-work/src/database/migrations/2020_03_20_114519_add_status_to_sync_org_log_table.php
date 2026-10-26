<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddStatusToSyncOrgLogTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('sync_org_log')) {
            Schema::table('sync_org_log', function (Blueprint $table) {
                $table->unsignedTinyInteger('status')->default(0)->after('admin_name')->comment('同步状态，0=uc同步中，1=uc同步完成,2=自动创建联系人二维码成功(真正同步成功)');
            });
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        if (Schema::hasTable('sync_org_log')) {
            Schema::table('sync_org_log', function (Blueprint $table) {
                $table->dropColumn('status');
            });
        }
    }
}
