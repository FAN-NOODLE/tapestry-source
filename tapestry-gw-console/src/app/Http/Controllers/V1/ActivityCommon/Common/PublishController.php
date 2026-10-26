<?php


namespace App\Http\Controllers\V1\ActivityCommon\Common;


use App\Http\Controllers\V1\ActivityCommon\ActivityBaseController;

class PublishController extends ActivityBaseController
{
    function index()
    {
        $params = $this->params;

        try {
            $activityInfo = $this->activityApi->getActivity($params);
        } catch (\Exception $e) {
            $this->notice($e, $params);
        }

        return success($activityInfo);
    }
}