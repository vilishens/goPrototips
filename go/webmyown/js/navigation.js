$("#stationrestart").on('click', function(){

    var urlStr = '/station/restart';

    DoAjax(urlStr, {}, 500);
})

$("#stationscanip").on('click', function(){

    var urlStr = '/station/scanip';
    
    DoAjax(urlStr, {}, 500);
})

$("#stationexit").on('click', function(){

    var urlStr = '/station/exit';
    
    DoAjax(urlStr, {}, 500);
})