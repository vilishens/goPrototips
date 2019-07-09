var POINT_LIST_DATA = 'pointListData';
var POINT_LIST_ITEM = 'pointListItem';
var ITEM_CLASS_DEFAULT = 'btn-outline-secondary';
var ITEM_CLASS_SIGNED = 'btn-outline-success';
var ITEM_CLASS_DISCONNECTED = 'btn-outline-secondary button-blink';
//var CLASS_ITEM_FROZEN = 'outline-danger';
var ITEM_ID_PREFIX = 'ptItem';

var URL_POINT_LIST="/pointlist/data";
var URL_POINT_ITEM_CFG ="/pointlist/act/cfg/"
var URL_POINT_ITEM_RESTART ="/pointlist/act/start/"
var AJAX_DATA={};

function makeList() {
    handlePointList()
    var nbr = SetInterv(-5, "handlePointList()", 1500);   // 1.5 sec
}

function handlePointList() {
 
    AJAX_DATA = {};

    $.ajax({
        url: URL_POINT_LIST,
        type: 'post',
        data: AJAX_DATA, //JSON.stringify(d), 
        dataType: 'json',
        contentType: 'application/json;charset=utf-8',
        async: true,
        timeout: 500,   // 0.5 second
        success : function(data, status, xhr) {
            drawPointList(data);
        },
        error : function(request,error) {
            alert("Error: "+error);
        },
    });
}


function drawPointList(d) {

    clean = emptyList(d);

    var htmlStr = '';
    for(ptn in d["List"]) {

        var name = d["List"][ptn];
        var item = d["Data"][name];

        var cl = itemClass(d["Data"][name]);

        var all = $('#'+POINT_LIST_DATA);
 
        var itObj = all.find("#"+listItemId(name));

        if (clean) {
            htmlStr += drawPointListItem(d, cl, name);
            htmlStr += drawPointListItemX(d, cl, name);

htmlStr += miklo(d["Data"][name], cl, name);

        } else if (itObj.length > 0) {
            setItemClass(itObj, cl);
        }   
    }

    if(clean) {
        var obj = $('#'+POINT_LIST_DATA);
        obj.html(htmlStr);
    }    
}

function emptyList(d) {
    var obj = $('#'+POINT_LIST_DATA);

    var k = obj.find('.row');
    var v = obj.find('.row').find('.'+ POINT_LIST_ITEM);

//    var ki = obj.find('.'+ POINT_LIST_ITEM).length;
//    var ko = Object.keys(d).length;

    if(obj.find('.row').find('.'+ POINT_LIST_ITEM).length != Object.keys(d).length) {
        // the current list and data have different count of points, the list needs to be recreated 
        obj.empty();
        return true;
    }

    for(ind in d["List"]) {
        var name = d["List"][ind];
        var item = '#'+listItemId(name);

        var itemHas = obj.find(item);

        if(0 == itemHas.length) {
            // couldn't find the item in the current list
            // it means there are different items in data, the list needs to be recreated
            obj.empty();
            return true;
        } else {

            cl = itemClass(d["Data"][name]);

            if(!itemHas.hasClass(cl)) {
                // the item doesn't have the received class
                obj.empty();
                return true;
            }    
        }

    }
    return false;
}

function itemClass(item) {

    var cl = ITEM_CLASS_DEFAULT;

    if(item["Signed"] && item["Disconnected"]) {
        cl = ITEM_CLASS_DISCONNECTED;
    } else if(item["Signed"]) {
        cl = ITEM_CLASS_SIGNED;
    }

    return cl;
}

function listItemId(name) {
    return ITEM_ID_PREFIX+name;
}

//#############################################################
//#############################################################
//#############################################################

function setItemClass(obj, cl) {
    if (obj.hasClass(cl)) {
        return;
    }

    // not the right class let's set the right one
    obj.find('#'+item).removeClass(CLASS_ITEM_DEFAULT+' '+CLASS_ITEM_FROZEN+' '+CLASS_ITEM_ACTIVE);
    obj.find('#'+item).addClass(cl);
}

//@@@@@@@@@@@@@@@@@

function drawPointListItemX(d, cl, name) {

 //   var name = d["Data"];

  
    var str = '';
    str += '<div class="dropdown">';
    str +=      '<button class="btn btn-secondary dropdown-toggle" type="button" id="dropdownMenuButton" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">'
     + name;
    str +=      '</button>';

    str +=      '<div class="dropdown-menu" aria-labelledby="dropdownMenuButton">';

//#######################

/*    str +=         '<div class="dropdown show">';
    str +=         '    <a class="btn btn-secondary dropdown-toggle" href="#" role="button" id="dropdownMenuLink" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">';
    str +=         '        Dropdown link';
    str +=         '    </a>';

    str +=         '    <div class="dropdown-menu" aria-labelledby="dropdownMenuLink">';
    str +=         '        <a class="dropdown-item" href="#">Action</a>';
    str +=         '        <a class="dropdown-item" href="#">Another action</a>';
    str +=         '        <a class="dropdown-item" href="#">Something else here</a>';
    str +=         '    </div>';
    str +=         '</div>';
*/
   // str +=      '  <a class="dropdown-item" href="#">Action</a>';



//#######################




    str +=      '  <a class="dropdown-item" href="#">Action</a>';
    str +=      '  <a class="dropdown-item" href="#">Another action</a>';
    str +=      '  <a class="dropdown-item" href="#">Something else here</a>';
    str +=      '</div>';
    str += '</div>';

    return str;
}

function bootstrapa_menu(d, cl, name) {
/*    
    <div class="container">
  	<div class="row">
        <h2>Multi level dropdown menu in Bootstrap</h2>
        <hr>
        <div class="dropdown">
            <button class="btn btn-secondary dropdown-toggle" type="button" id="dropdownMenu1" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
              Dropdown
            </button>
            <ul class="dropdown-menu multi-level" role="menu" aria-labelledby="dropdownMenu">
                <li class="dropdown-item"><a href="#">Some action</a></li>
                <li class="dropdown-item"><a href="#">Some other action</a></li>
                <li class="dropdown-divider"></li>
                <li class="dropdown-submenu">
                  <a  class="dropdown-item" tabindex="-1" href="#">Hover me for more options</a>
                  <ul class="dropdown-menu">
                    <li class="dropdown-item"><a tabindex="-1" href="#">Second level</a></li>
                    <li class="dropdown-submenu">
                      <a class="dropdown-item" href="#">Even More..</a>
                      <ul class="dropdown-menu">
                          <li class="dropdown-item"><a href="#">3rd level</a></li>
                            <li class="dropdown-submenu"><a class="dropdown-item" href="#">another level</a>
                            <ul class="dropdown-menu">
                                <li class="dropdown-item"><a href="#">4th level</a></li>
                                <li class="dropdown-item"><a href="#">4th level</a></li>
                                <li class="dropdown-item"><a href="#">4th level</a></li>
                            </ul>
                          </li>
                            <li class="dropdown-item"><a href="#">3rd level</a></li>
                      </ul>
                    </li>
                    <li class="dropdown-item"><a href="#">Second level</a></li>
                    <li class="dropdown-item"><a href="#">Second level</a></li>
                  </ul>
                </li>
              </ul>
        </div>
    </div>
</div>
*/
}
 

//@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

function miklo(d, cl, name) {

    var str = '';

    str += '<div class="container">';
  	str += '    <div class="row">';
    str += '        <div class="dropdown">';
    str += '            <button class="btn dropdown-toggle '+cl+'" type="button" id="dropdownMenu1" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">';
    str += '                '+d["Point"];
    str += '            </button>';
    str += '            <ul class="dropdown-menu multi-level" role="menu" aria-labelledby="dropdownMenu">';
    str += '                <li class="dropdown-item"><a href="#">Some action</a></li>';
    str += '                <li class="dropdown-item"><a href="#">Some other action</a></li>';
    str += '                <li class="dropdown-divider"></li>';
    str += '                <li class="dropdown-submenu">';
    str += '                    <a  class="dropdown-item" tabindex="-1" href="#">Hover me for more options</a>';
    str += '                    <ul class="dropdown-menu">';
    str += '                        <li class="dropdown-item"><a tabindex="-1" href="#">Second level</a></li>';
    str += '                        <li class="dropdown-submenu">';
    str += '                            <a class="dropdown-item" href="#">Even More..</a>';
    str += '                            <ul class="dropdown-menu">';
    str += '                                <li class="dropdown-item"><a href="#">3rd level</a></li>';
    str += '                                <li class="dropdown-submenu"><a class="dropdown-item" href="#">another level</a>';
    str += '                                    <ul class="dropdown-menu">';
    str += '                                        <li class="dropdown-item"><a href="#">4th level</a></li>';
    str += '                                        <li class="dropdown-item"><a href="#">4th level</a></li>';
    str += '                                        <li class="dropdown-item"><a href="#">4th level</a></li>';
    str += '                                    </ul>';
    str += '                                </li>';
    str += '                                <li class="dropdown-item"><a href="#">3rd level</a></li>';
    str += '                           </ul>';
    str += '                        </li>';
    str += '                        <li class="dropdown-item"><a href="#">Second level</a></li>';
    str += '                        <li class="dropdown-item"><a href="#">Second level</a></li>';
    str += '                    </ul>';
    str += '                </li>';
    str += '            </ul>';           
    str += '        </div>';
    str += '    </div>';
    str += '</div>';

    return str;

}

//@@@@@@@@@@@@@@@@@@@@@@@@@@@@@



function drawPointListItem(d, cl, name) {

    var ptnDscr = "kiril"; //d["Descr"];
    var ptnId = ITEM_ID_PREFIX+name;
 
    var str = '';
    str += 
    str += '<div class="row mt-2">'
    str += '    <div class="btn-group dropright '+POINT_LIST_ITEM+'">';
    str += '        <a class="btn '+cl+' dropdown-toggle" href="#"'; 
    str += '            role="button" id="'+ptnId+'" data-toggle="dropdown"' 
    str += '            data-toggle="tooltip" data-placement="right" title="'+ptnDscr+'"';
    str += '            aria-haspopup="true" aria-expanded="false">'; 
    str += '            '+name;
    str += '        </a>';

//    var URL_POINT_ITEM_CFG ="/pointlist/act/cfg/"
 //   var URL_POINT_ITEM_RESTART ="/pointlist/act/start/"
    


    str += '        <div class="dropdown-menu">';
    str += '            <a class="dropdown-item" href="'+URL_POINT_ITEM_CFG+name+'">Configuration</a>';


    var isDisconn = d["Data"][name]["Disconnected"];
                    if(isDisconn)
                    {
    str += '            <a class="dropdown-item" href="'+URL_POINT_ITEM_RESTART+name+'">Restart</a>';
                    }    
    str += '        </div>';

    str += '    </div>';
    str += '</div>';

    return str;
}

/*
<div class="row mt-2">
<div class="btn-group dropright">

    {{ $kika := index pointList.Data $x}}

    {{ $mika := $kika.Descr }} 

    {{ if $kika.Frozen }}
        <a class="btn btn-outline-danger dropdown-toggle" href="#" 
        role="button" id="pointChoice" data-toggle="dropdown" 
        data-toggle="tooltip" data-placement="right" title={{ $mika }}
        aria-haspopup="true" aria-expanded="false">                    
    {{ else if $kika.Active }}
        <a class="btn btn-outline-success dropdown-toggle" href="#" 
        role="button" id="pointChoice" data-toggle="dropdown" 
        data-toggle="tooltip" data-placement="right" title={{ $mika }}
        aria-haspopup="true" aria-expanded="false">                    
    {{ else }}             
        <a class="btn btn-outline-secondary dropdown-toggle" href="#" 
        role="button" id="pointChoice" data-toggle="dropdown" 
        data-toggle="tooltip" data-placement="right" title={{ $mika }}
        aria-haspopup="true" aria-expanded="false">                    
    {{ end }}
        {{ $x }}
    </a>
    <div class="dropdown-menu">
        <a class="dropdown-item" href="point/{{ $x }}/showcfg">Configuration</a>
        <a class="dropdown-item" href="#">Another action</a>
        <a class="dropdown-item" href="#">Something else here</a>
    </div>
</div>
</div>
*/
