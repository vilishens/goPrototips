function DoAjax(urlStr, data, timeOut) {
    $.ajax({  
        url: urlStr,
        type: 'post',
        data: data, 
        dataType: 'json',
        contentType: 'application/json;charset=utf-8',
        async: true,
        timeout: timeOut,   // 500 == 0.5 second
        success : function(data, status, xhr) {
            return;
        },
        error : function(request,error) {
            alert("Request: "+JSON.stringify(request)+", Error: "+error);
        },
    });
}