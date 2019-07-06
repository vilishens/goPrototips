var ALL_POINT_LIST = 'allPointList';
var CLASS_LIST_ITEM = 'pointListItem';
var CLASS_ITEM_DEFAULT = 'btn-outline-secondary';
var CLASS_ITEM_ACTIVE = 'btn-outline-success';
var CLASS_ITEM_FROZEN = 'btn-outline-danger';
var ID_ITEM_PREFIX = 'ptList';

function makeList() {

    handlePointList()

    var nbr = SetInterv(-5, "handlePointList()", 1500);   // 1.5 sec

}

function htmlPointList(d) {

    var obj = $('#allPointList');

    obj.empty()


}

function handlePointList() {
 
    var urlStr =  "/pointlist/data";

//    return;

    $.ajax({
        url: urlStr,
        type: 'post',
        data: {}, 
        dataType: 'json',
        timeout: 500,
        success : function(d) {

            var data = d["Data"];

            drawPointList(data)
        },
    });
}

function emptyList(d) {
    var obj = $('#'+ALL_POINT_LIST);

    var k = obj.find('.row');
    var v = obj.find('.row').find('.'+ CLASS_LIST_ITEM);

    var ki = obj.find('.'+ CLASS_LIST_ITEM).length;
    var ko = Object.keys(d).length;

    if(obj.find('.row').find('.'+ CLASS_LIST_ITEM).length != Object.keys(d).length) {
        // the current list and data have different count of points, the list needs to be recreated 
        obj.empty();
        return true;
    }

    for(pt in d) {
        var item = '#'+listItemId(pt);
        if(0 == obj.find(item).length) {
            // couldn't find the item in the current list
            // it means there are different items in data, the list needs to be recreated
            obj.empty();
            return true;
        }
    }

    return false;
}

function listItemId(name) {
    return ID_ITEM_PREFIX+name;
}

function setItemClass(d) {
    var obj = $('#'+ALL_POINT_LIST);
    var item = listItemId(d['Point']);

    var iClass = CLASS_ITEM_DEFAULT;
    if(IsPointFrozen(d['State'])) {
        iClass = CLASS_ITEM_FROZEN;
    } else if (IsPointActive(d['State'])) {
        iClass = CLASS_ITEM_ACTIVE;
    }

    if (obj.find('#'+item).hasClass(iClass)) {
        return;
    }

    // not the right class let's set the right one
    obj.find('#'+item).removeClass(CLASS_ITEM_DEFAULT+' '+CLASS_ITEM_FROZEN+' '+CLASS_ITEM_ACTIVE);
    obj.find('#'+item).addClass(iClass);
}

function drawPointList(d) {

    clean = emptyList(d);

    var htmlStr = '';
    for(ptn in d) {
        if (clean) {
            htmlStr += drawPointListItem(d[ptn]);
        } else {
            setItemClass(d[ptn]);
        }   
    }

    if(clean) {
        var obj = $('#'+ALL_POINT_LIST);
        obj.html(htmlStr);
    }    
}

function drawPointListItem(d) {

    var pClass = CLASS_ITEM_DEFAULT;
    if(IsPointFrozen(d["State"])) {
        pClass = CLASS_ITEM_FROZEN;
    } if(IsPointActive(d["State"])) {
        pClass = CLASS_ITEM_ACTIVE;
    }

    var ptnDscr = d["Descr"];
    var ptnId = ID_ITEM_PREFIX+d['Point'];
 
    var str = '';
    str += 
    str += '<div class="row mt-2">'
    str += '    <div class="btn-group dropright '+CLASS_LIST_ITEM+'">';
    str += '        <a class="btn '+pClass+' dropdown-toggle" href="#"'; 
    str += '            role="button" id="'+ptnId+'" data-toggle="dropdown"' 
    str += '            data-toggle="tooltip" data-placement="right" title="'+ptnDscr+'"';
    str += '            aria-haspopup="true" aria-expanded="false">'; 
    str += '            '+d["Point"];
    str += '        </a>';

    str += '        <div class="dropdown-menu">';
    str += '            <a class="dropdown-item" href="point/'+ptn+'/showcfg">Configuration</a>';
    str += '            <a class="dropdown-item" href="#">Another action</a>';
    str += '            <a class="dropdown-item" href="#">Something else here</a>';
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
