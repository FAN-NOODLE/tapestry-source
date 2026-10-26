<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddServeTypeToTask extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('task')) {
            Schema::table('task', function (Blueprint $table) {
                $table->unsignedInteger('serve_type')
                    ->default(0)
                    ->after('task_type')
                    ->comment('任务类型偏向:0=crs+shop;1=shop;2=crs');
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
        if (Schema::hasTable('task')) {
            Schema::table('task', function (Blueprint $table) {
                $table->dropColumn('serve_type');
            });

        }
    }
}
