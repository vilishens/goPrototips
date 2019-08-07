$("#stationrescanwhole").on('click', function(){

    var urlStr = 'station/act/rescanwhole';

    DoAjax(urlStr, {}, 500);
})

$("#stationreboot").on('click', function(){

    var urlStr = '/station/act/reboot';
    
    DoAjax(urlStr, {}, 500);
})

$("#stationexit").on('click', function(){

    var urlStr = 'station/act/exit';
    
    DoAjax(urlStr, {}, 500);
})

$("#stationrestart").on('click', function(){

    var urlStr = 'station/act/restart';
    
    DoAjax(urlStr, {}, 500);
})

$("#stationshutdown").on('click', function(){

    var urlStr = 'station/act/exit';
    
    DoAjax(urlStr, {}, 500);
})


//<a class="dropdown-item" id="" href="station/act/restart">Restart</a>
//<a class="dropdown-item" id="stationshutdown" href="station/act/shutdown">Shutdown</a>
//****  */<a class="dropdown-item" id="stationreboot" href="station/act/reboot">Reboot</a>
