var BIT_STATE_ACTIVE = 0x0001
var BIT_STATE_FREEZE = 0x0002

function SetInterv(nbr, todo, interv) {
    if(0 > nbr) {
        nbr = setInterval(todo, interv); 
    }    

    return nbr;
}

function UnsetInterv(nbr) {
    if(0 <= nbr) {
        clearInterval(nbr);
        nbr = -1;
    }    

    return nbr;
}

function IsPointFrozen(bit) {
    return (0 != (bit & BIT_STATE_FREEZE));
}

function IsPointActive(bit) {
    return (0 != (bit & BIT_STATE_ACTIVE));
}

function ReturnData(url, d) {
    $.ajax({  
        url: url,
        type: 'post',
        data: JSON.stringify(d), 
        dataType: 'json',
        contentType: 'application/json;charset=utf-8',
        async: true,
        timeout: 500,   // 0.5 second
        success : function(data, status, xhr) {
 //           alert("Data "+ data + " STATUS " + status + " XHR " +xhr);
            return;
        },
        error : function(request,error) {
            alert("Request: "+JSON.stringify(request)+", Error: "+error);
        },
    });
}