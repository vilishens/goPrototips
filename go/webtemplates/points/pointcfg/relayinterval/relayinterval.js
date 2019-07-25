var URL_PAGE_HANDLER="/point/cfg/"; 

var THIS_POINT = "";
var THIS_CFG = 0x000001; //relay intervals

var STATE_EDIT = 0x0001;

var BTN_EDIT = "btnEdit";
var BTN_FREEZE = "btnFreeze";
var BTN_LOAD = "btnLoad";
var BTN_LOAD_DEFAULT = "btnLoadDefault";
var BTN_LOAD_SAVED = "btnLoadSaved";
var BTN_SAVE = "btnSave";

var BTN_EDIT_TXT = "Edit";
var BTN_FREEZE_TXT = "Freeze";
var BTN_LOAD_TXT = "Load";
var BTN_LOAD_DEFAULT_TXT = "Load default";
var BTN_LOAD_SAVED_TXT = "Load saved";
var BTN_SAVE_TXT = "Save";

var BTN_CLASS_ACTIVE = "btn-warning active";
var BTN_CLASS_INACTIVE = "btn-outline-secondary disabled";
var BTN_CLASS_IN_USE = 'btn-success active';

var TABLE_START = "tableStart";
var TABLE_BASE = "tableBase";
var TABLE_FINISH = "tableFinish";

var TABLE_START_TEXT = "Start";
var TABLE_BASE_TEXT = "Base";
var TABLE_FINISH_TEXT = "Finish";

var TD_CLASS_EDIT_OK = "tdEditOk";
var TD_CLASS_EDIT_ERROR = "tdEditError";
var TD_CLASS_EDIT_NONE = "tdEditNone";
var TD_CLASS_EDIT_ONLY = "tdEditOnly";
var TD_CLASS_EDIT_TAB_HEAD = "tdEditTabHead";
var TD_CLASS_EDIT_ADD = "tdEditAdd"
var TD_CLASS_EDIT_DELETE = "tdEditDelete"
var TD_CLASS_EDIT = 'tdEdit';

var TR_CLASS_HEADER = 'trEditHeader';
var TR_CLASS_ACTIVE_ROW = 'active-row';
var TR_CLASS_DRAGGED = "trDragged";
var TR_CLASS_EDIT = 'trEdit';
var TABLE_CLASS_ROW_NEW = "this-is-a-new-row";

var J_BUTTON_LABEL_ADD = "Add";
var J_BUTTON_LABEL_DELETE = "Del";

var COLOR_EDITED = '#C6E710';
var COLOR_ERROR = '#E9999A';
//var COLOR_DRAG = '#90EE90';
var COLOR_DRAG = '#FFE4B5';

var AllD = {};
var CfgDefault = {};
var CfgRun = {};
var CfgSaved = {};
var CfgIndex = {};
var CfgState = 0;

var ThisState = 0;

var ColorOri;
var FontWeightOri;

function makePage(name) {
    THIS_POINT = name;
    handlePointCfg()

    var nbr = SetInterv(-5, "handlePointCfg()", 1000);   // 1 sec
}




function setHandlersAndParams() {
    $('.btnMngmt').on('click', function(){btnClick($(this));}); 

   
    if(!editState()) {
        FontWeightOri = $('.' + TD_CLASS_EDIT).css('font-weight');
    }    
}

function handlePointCfg() {
 
    AllD = {};

    $.ajax({
        url: URL_PAGE_HANDLER+THIS_POINT+"/"+THIS_CFG.toString(),
        type: 'post',
        data: {}, //JSON.stringify(d), 
        dataType: 'json',
        contentType: 'application/json;charset=utf-8',
        async: true,
        timeout: 500,   // 0.5 second
        success : function(data, status, xhr) {
            setAllData(data);
            drawPage();
            setHandlersAndParams();
        },
//        error : function(request,error) {
//            alert("Error: "+error);
//        },
    });
}

function setAllData(data) {

    AllD = data;
    cfgCd = THIS_CFG.toString();

    CfgDefault = AllD["CfgDefault"][cfgCd];
    CfgRun = AllD["CfgRun"][cfgCd];
    CfgSaved = AllD["CfgSaved"][cfgCd];
    CfgIndex = AllD["CfgIndex"][cfgCd];
    CfgState = AllD["CfgState"][cfgCd];
}

function drawPage() {
    drawTitle();
    drawButtons();
    if(!editState()) {
        drawCfg();
    }    
}

function editState() {
    return !(0 == (ThisState & STATE_EDIT))
}

function drawTitle() {
    $('#pointName').text(THIS_POINT);
}

function drawButtons() {
    if( editState()) {

    } else {
        drawBtn(BTN_FREEZE, BTN_FREEZE_TXT);
        drawBtn(BTN_LOAD, BTN_LOAD_TXT);
        drawBtn(BTN_SAVE, BTN_SAVE_TXT);
        drawBtn(BTN_LOAD_DEFAULT, BTN_LOAD_DEFAULT_TXT);
        drawBtn(BTN_LOAD_SAVED, BTN_LOAD_SAVED_TXT);
        drawBtn(BTN_EDIT, BTN_EDIT_TXT);
    }
}

function drawBtn(id, str) {
    $('#'+id).text(str);
    setAvailable(id);
}

function setAvailable(id) {
    var avail = isAvailable(id);
    var btn = $('#'+id);

    btn.removeClass(BTN_CLASS_ACTIVE);
    btn.removeClass(BTN_CLASS_INACTIVE);

    if(!avail) {
        btn.addClass(BTN_CLASS_INACTIVE);
    } else {
        btn.addClass(BTN_CLASS_ACTIVE);
    }  
}

function isAvailable(id) {
    var edit = editState();

    if(edit &&(id != BTN_EDIT)) {
        return false;
    }

    return true;
}

function drawCfg() {
    drawCfgTable(CfgRun["Start"],TABLE_START, TABLE_START_TEXT, CfgIndex["Start"]);
    drawCfgTable(CfgRun["Base"],TABLE_BASE, TABLE_BASE_TEXT, CfgIndex["Base"]);
    drawCfgTable(CfgRun["Finish"],TABLE_FINISH, TABLE_FINISH_TEXT, CfgIndex["Finish"]);
}

function drawCfgTable(data, table, title, ind) {
    var obj = $('#' + table);
    var rowCount = ((null == data) || (0 == data.length)) ? 0 : data.length;

    obj.empty()
    var str = "";

    str += tableTitle(data, title);
    
    str += '<table id="editable-def" dropzone="move" class="pure-table pure-table-bordered">';

    str += tableTabHead();

    var i = 0;
    for (i = 0; i < rowCount; i++) {
        str += tableTabRow(data[i], i, ind, false);
    }

    str += tableTabRowNew();

    str += '</table>';
    str += '</br>';

//    if(!editState()) {
//        obj.find('td.' + TD_CLASS_EDIT_ONLY).hide();
//    }    

    obj.html(str);

    createButtonAdd(obj.find('.' + TD_CLASS_EDIT_ADD));
    createButtonDelete(obj.find('.' + TD_CLASS_EDIT_DELETE));

   obj.find('.' + TD_CLASS_EDIT_ONLY).hide();
}

function createButtonAdd(o) {
    o.button({
        label:J_BUTTON_LABEL_ADD, 
        icons:{primary:'ui-icon-plusthick'}
    });

    o.button().on('click', function() {
        jButtonClick($(this));
    })
}

function tableTitle(data, title) {

    var str = "";

    str += '<h2>' + title + '</h2>';

    return str;
}


function tableTabHead() {
    
    var str = "";

    str += '<thead>';
    str += '    <tr class="'+TR_CLASS_HEADER+'">';
    str += '        <th>GPIO</th>';
    str += '        <th>STATE</th>';
    str += '        <th>INTERVAL</th>';
    str += '        <th class="' + TD_CLASS_EDIT_ONLY + ' ' + TD_CLASS_EDIT_TAB_HEAD + '"></th>';
    str += '     </tr>';
    str += '</thead>';

    return str;
}

function tableTabRow(data, i, ind, isNew) {;

    var str = "";

    var trClass = "trEdit" + (isNew ? (" " + TABLE_CLASS_ROW_NEW) : "");
    if(!isNew && (i == ind)) {
        trClass += ' ' + TR_CLASS_ACTIVE_ROW;
    }

    str += '<tr draggable="';
    str += isNew ? "false" : "true";
 
    str += '" class="' + trClass + '">';

    str += partTabCols(data);

    if(isNew) {
        // add button "add"
        str += '<td class="' + TD_CLASS_EDIT_ADD + ' ' + TD_CLASS_EDIT_ONLY + '"></td>';
    } else {
        // add button "delete"
        str += '<td class="' + TD_CLASS_EDIT_DELETE + ' ' + TD_CLASS_EDIT_ONLY + '"></td>';
    }    
    
    return str;
}

function seconds2Interval(secsFull) {

    if(isNaN(secsFull)) {
        return "new:new:new";
    }

    var secs = Math.trunc(secsFull / 1000000000);

    var hr = Math.trunc(secs / 3600); // Hours
    var min = Math.trunc((secs % 3600) / 60); // Minutes
    var sec = Math.trunc(secs % (3600 * 60)); // secoonds

    var str = ((hr < 10) ? "0" : "") + hr.toString() + ":" +
              ((min < 10) ? "0" : "") + min.toString() + ":" + 
              ((sec < 10) ? "0" : "") + sec.toString();

    return str;
}


function partTabCols(data) {

    var str = "";

    var go = data["Seconds"];
    var f = typeof go;
    var sis = 5;

/*    
    var intervHr = intervN / 3600; // Hours
    intervN -= 3600 *intervHr;
    var intervMin = intervN / 60; // Minutes
    var intervSecs = interVN - (60 * intervMin);
*/
    if("Interval" in data) {
        // this row is a template for the new row
        interv = data["Interval"];
    } else {
        var interv = seconds2Interval(data["Seconds"]);
    }    
 
    str += '<td class="tdEdit tdEditGpio"     data-ori="' + data["Gpio"] +     '">' + data["Gpio"] + '</td>';
    str += '<td class="tdEdit tdEditState"    data-ori="' + data["State"] +    '">' + data["State"] + '</td>';
    str += '<td class="tdEdit tdEditInterval" data-ori="' + interv + '">' + interv + '</td>';

    return str;
}

function tableTabRowNew() {;

    var data = {Gpio:"new", State:"new", Interval:"new:new:new"};
    var str = tableTabRow(data, -1, -2, true);

    return str;
}

function btnClick(btn) {

    var which = btn.prop('id');

    if(isButtonInactive(btn)) {
        return 
    }

    switch(which) {
        case BTN_EDIT:
            btnEditPressed(btn);  
            break;
        case BTN_LOAD:
            btnLoadPressed(btn);
            break; 
        case BTN_FREEZE:
            btnFreezePressed(btn);
            break;
        case BTN_SAVE:
            btnSavePressed(btn);
            break;
        case BTN_LOAD_DEFAULT:
            btnLoadDefaultPressed(btn);
            break;
        case BTN_LOAD_SAVED:
            btnLoadSavedPressed(btn);
            break;
 
        default:
            alert("Button "+which+" pressed which doesn't have logic");
            break;    
    }
}

function isButtonInactive(btn) {

    var cls = BTN_CLASS_INACTIVE.split(" ");
    for (ind in cls) {
        if(!btn.hasClass(cls[ind])) {
            return false;
        }
    }

    return true;
}

function isButtonInUse(btn) {

    var cls = BTN_CLASS_IN_USE.split(" ");
    for (ind in cls) {
        if(!btn.hasClass(cls[ind])) {
            return false;
        }
    }

    return true;
}

function isButtonActive(btn) {

    var cls = BTN_CLASS_ACTIVE.split(" ");
    for (ind in cls) {
        if(!btn.hasClass(cls[ind])) {
            return false;
        }
    }

    return true;
}

function btnEditPressed(btn) {
    if(isButtonInactive(btn)) {
        return 
    }

    if(isButtonActive(btn)) {
        // activeonly, set in use
        ThisState |= STATE_EDIT;
        setButtonInUse(btn);
        setAllTableEditOptions();
   }
    else if(isButtonInUse(btn)) {
        // in use, set active only
        ThisState &= ~STATE_EDIT;
        setButtonActive(btn);
//        unsetAllTableEditOptions();
    }
}

function setAllTableEditOptions() {

    var trcl = '.' + TR_CLASS_EDIT;

    $(trcl).removeClass(TR_CLASS_ACTIVE_ROW);

    $(trcl).attr('contenteditable', 'true');
    $(trcl).attr('oninput', 'checkInput($(this))');
    setAllTablesDraggable();

    var tdcl = '.' + TD_CLASS_EDIT_DELETE;
    $(tdcl).show();
    var tdcl = '.' + TD_CLASS_EDIT_TAB_HEAD;
    $(tdcl).show();
  
//    setTablesAddButton()

    ThisState |= STATE_EDIT;

  //  setEditButtons();
}    


function setButtonInUse(btn) {
    var cls = BTN_CLASS_INACTIVE.split(" ");
    for (ind in cls) {
        btn.removeClass(cls[ind]);
    }
    cls = BTN_CLASS_ACTIVE.split(" ");
    for (ind in cls) {
        btn.removeClass(cls[ind]);
    }

    cls = BTN_CLASS_IN_USE.split(" ");
    for (ind in cls) {
        btn.addClass(cls[ind]);
    }



//    btn.removeClass('btn-outline-secondary').removeClass('btn-warning').removeClass('disabled');
//    btn.addClass('btn-success').addClass('active');
}

function setButtonActive(btn) {
    var cls = BTN_CLASS_INACTIVE.split(" ");
    for (ind in cls) {
        btn.removeClass(cls[ind]);
    }
    cls = BTN_CLASS_IN_USE.split(" ");
    for (ind in cls) {
        btn.removeClass(cls[ind]);
    }

    cls = BTN_CLASS_ACTIVE.split(" ");
    for (ind in cls) {
        btn.addClass(cls[ind]);
    }



//    btn.removeClass('btn-outline-secondary').removeClass('btn-warning').removeClass('disabled');
//    btn.addClass('btn-success').addClass('active');
}

function jButtonClick(btn) {

    var label = btn.button('option', 'label');

    if (J_BUTTON_LABEL_DELETE == label) {
  //      htmlRemoveTdRow(btn);
    }

    if (J_BUTTON_LABEL_ADD == label) {
//        htmlAddNewRow(btn);
    }
}

function setAllTablesDraggable() {
    setTableDraggableOption(TABLE_START);
    setTableDraggableOption(TABLE_BASE);
    setTableDraggableOption(TABLE_FINISH);
}

function setTableDraggableOption(tbl) {
    var tb = $('#' + tbl);
    tb.find('.'+ TR_CLASS_EDIT).attr('draggable', 'true');

    var removeDragg = $(tb.find('.' + TABLE_CLASS_ROW_NEW).last());
    removeDragg.attr('draggable', 'false');

    setTableSortedOption(tbl);
}

function setTableSortedOption(tab) {
    var tbSort = $('#' + tab).find('tbody');

    tbSort.sortable({
        items: "tr[draggable='true']",
        update:function(event, ui){
            $(ui.item).css('background-color', COLOR_DRAG);
            $(ui.item).addClass(TR_CLASS_DRAGGED);
            //inputReady2Use();
        }
    });

   // tbSort.on('dblclick', function(){toggleSort($(this));});       
   tbSort.on('dblclick', function(){toggleSortAll();});       
   tbSort.sortable('disable');
}

function toggleSortAll() {
    //    var tbSort = $('#' + TABLE_START).find('tbody');
    toggleSort($('#' + TABLE_START).find('tbody'));       
    toggleSort($('#' + TABLE_BASE).find('tbody'));       
    toggleSort($('#' + TABLE_FINISH).find('tbody'));       
    
    
    
    //    tbSort.sortable('disable');
}
    
function toggleSort(obj) {
    var opts = obj.sortable('option');

    if(opts["disabled"]) {
        // disabled, set active
        obj.sortable( "enable" );
        $("tr[draggable='true']").find('.' + TD_CLASS_EDIT).css("font-weight","bold");
    } else {
        // active, set disabled
        $("tr[draggable='true']").find('.' + TD_CLASS_EDIT).css("font-weight", FontWeightOri);
        obj.sortable( "disable" );
    }
};

function createButtonDelete(o) {
    o.button({
        label:J_BUTTON_LABEL_DELETE, 
        icons:{secondary:' ui-icon-closethick'}
    })

    o.button().on('click', function() {
        jButtonClick($(this));
    })
}

function checkInput(td) {
/*
    var str = td.html();
    var ori = td.attr('data-ori');
    var color = colorOri;
    if(td.closest('tr').hasClass(TR_CLASS_DRAGGED)) {
        color = COLOR_DRAG;
    } 


    var tdClass = TD_CLASS_EDIT_NONE;
    td.removeClass(TD_CLASS_EDIT_OK);
    td.removeClass(TD_CLASS_EDIT_ERROR);

    var ok = false;

    if(ori != str) {
 
        if(td.hasClass("tdEditGpio")) {
            ok = checkInputGpio(str.trim());
        }
        if(td.hasClass("tdEditState")) {
            ok = checkInputState(str.trim());
        }
        if(td.hasClass("tdEditInterval")) {
            ok = checkInputInterval(str.trim());
        }

        if(!ok) {
            color = COLOR_ERROR;
            tdClass = TD_CLASS_EDIT_ERROR;
        } else {
            color = COLOR_EDITED;
            tdClass = TD_CLASS_EDIT_OK;
        }
    }

    td.attr("style", 'background-color:' + color);
    td.addClass(tdClass);

    var tr = td.closest('tr');
    if (tr.hasClass(TABLE_CLASS_ROW_NEW)) {
        setTableAddTr(tr);
    }

    inputReady2Use();

    return ok;
*/    
}

/*

var cfgThis = {};
var cfgHave = {};
var cfgSaved = {};
var indexAll = {};
var pointName = "";

var start = true;

var intervN = -1;
var thisState = 0;
var isActive = false;
var isFreeze = false;

var colorOri;
var COLOR_EDITED = '#C6E710';
var COLOR_ERROR = '#E9999A';
//var COLOR_DRAG = '#90EE90';
var COLOR_DRAG = '#FFE4B5';

var fontWeightOri;

var STATE_EDIT = 0x000001;

var BTN_EDIT = "btnEdit";
var BTN_FREEZE = "btnFreeze";
var BTN_LOAD = "btnLoad";
var BTN_LOAD_DEFAULT = "btnLoadDefault";
var BTN_LOAD_SAVED = "btnLoadSaved";
var BTN_SAVE = "btnSave";

var BTN_EDIT_TXT = "Edit";
var BTN_FREEZE_TXT = "Freeze";
var BTN_LOAD_TXT = "Load";
var BTN_LOAD_DEFAULT_TXT = "Load default";
var BTN_LOAD_SAVED_TXT = "Load saved";
var BTN_SAVE_TXT = "Save";

var TD_CLASS_EDIT_OK = "tdEditOk";
var TD_CLASS_EDIT_ERROR = "tdEditError";
var TD_CLASS_EDIT_NONE = "tdEditNone";

var TR_CLASS_HEADER = 'trEditHeader';
var TR_CLASS_ACTIVE_ROW = 'active-row';
var TR_CLASS_DRAGGED = "trDragged";

var TABLE_START = "tableStart";
var TABLE_BASE = "tableBase";
var TABLE_FINISH = "tableFinish";

var TABLE_CLASS_ROW_NEW = "this-is-a-new-row";

var J_BUTTON_LABEL_ADD = "Add";
var J_BUTTON_LABEL_DELETE = "Del";

var POINT_STATE_ACTIVE = 0x0001
var POINT_STATE_FREEZE = 0x0002

var semjo = 0;

function loadCfg(name) {

    pointName = name;
    $('#pointName').text(pointName);

    handleCfg()

    setInterv();

    $('.btnMngmt').on('click', function(){btnClick($(this));});       
}

function checkInput(td) {

    var str = td.html();
    var ori = td.attr('data-ori');
    var color = colorOri;
    if(td.closest('tr').hasClass(TR_CLASS_DRAGGED)) {
        color = COLOR_DRAG;
    } 


    var tdClass = TD_CLASS_EDIT_NONE;
    td.removeClass(TD_CLASS_EDIT_OK);
    td.removeClass(TD_CLASS_EDIT_ERROR);

    var ok = false;

    if(ori != str) {
 
        if(td.hasClass("tdEditGpio")) {
            ok = checkInputGpio(str.trim());
        }
        if(td.hasClass("tdEditState")) {
            ok = checkInputState(str.trim());
        }
        if(td.hasClass("tdEditInterval")) {
            ok = checkInputInterval(str.trim());
        }

        if(!ok) {
            color = COLOR_ERROR;
            tdClass = TD_CLASS_EDIT_ERROR;
        } else {
            color = COLOR_EDITED;
            tdClass = TD_CLASS_EDIT_OK;
        }
    }

    td.attr("style", 'background-color:' + color);
    td.addClass(tdClass);

    var tr = td.closest('tr');
    if (tr.hasClass(TABLE_CLASS_ROW_NEW)) {
        setTableAddTr(tr);
    }

    inputReady2Use();

    return ok;
}

function checkInputGpio(str) {

    var dec = Number(str);

    if(isNaN(dec)) {
        // the string is not a correct numeric 
        return false;
    }

    if(!((0 < dec) && (100 > dec))) {
        // the number is not in a range (1 ... 99)
        return false;
    }    

    return true;
}

function checkInputState(str) {
    // allowed values 0 (off) and 1 (on)
    return (str == '0' || str == '1');
}

function checkInputInterval(str) {

    var decs = str.split(":"); 

    if(3 != decs.length) {
        // couldn't get 3 required parts
        return false;
    }

    var maxs = [23, 59, 59]; // max values hours, minutes, seconds
    var i = -1;
    var val = -1;
    var total = 0;
    var str = "";
    for (i = 0; i < decs.length; i++) { 

        str = decs[i].trim();

        if(2 != str.length) {
            // there aren't 2 chars in a part
            return false;
        }

        val = Number(str);
        if(isNaN(val)) {
            // the value string is not a correct numeric
            return false;
        }

        if (!((0 <= val) && (maxs[i] >= val))) {
            // value is either less than 0 or bigger than the max value
            return false;
        }

        total += val;
    } 

    if(!total) {
        // all parts are zero, an interval can't be a zero
        return false;
    }

    return true;
}


function btnClick(btn) {

    var which = btn.prop('id');

    if(isButtonInactive(btn)) {
        return 
    }

    switch(which) {
        case BTN_EDIT:
            btnEditPressed(btn);  
            break;
        case BTN_LOAD:
            btnLoadPressed(btn);
            break; 
        case BTN_FREEZE:
            btnFreezePressed(btn);
            break;
        case BTN_SAVE:
            btnSavePressed(btn);
            break;
        case BTN_LOAD_DEFAULT:
            btnLoadDefaultPressed(btn);
            break;
        case BTN_LOAD_SAVED:
            btnLoadSavedPressed(btn);
            break;
 
        default:
            alert("Button "+which+" pressed which doesn't have logic");
            break;    
    }
}

function btnLoadDefaultPressed(btn) {

    if(isButtonInactive(btn)) {
        return 
    }

    setButtonActive(btn);

    unsetAllTableEditOptions();

    loadDefault();

    setButtonInactive(btn);
}

function btnLoadSavedPressed(btn) {

    if(isButtonInactive(btn)) {
        return 
    }

    setButtonActive(btn);

    unsetAllTableEditOptions();

    loadSaved();

    setButtonInactive(btn);
}

function loadSaved() {
    var d = {};

    var urlStr =  "/point/handlecfg/" + pointName + "/loadsavedcfg";

    returnData(urlStr, d);
}

function loadDefault() {
    //var d = getInputData();

    var d = {};

  //  var url = $("#editReturnURL").text();

    //url += "/loadcfg/"+ JSON.stringify(d);

    var urlStr =  "/point/handlecfg/" + pointName + "/loaddefaultcfg";

    returnData(urlStr, d);
}


function btnSavePressed(btn) {

    if(isButtonInactive(btn)) {
        return 
    }

    setButtonActive(btn);

    unsetAllTableEditOptions();

    saveInputData();

    setButtonInactive(btn);
}

function btnLoadPressed(btn) {

    if(isButtonInactive(btn)) {
        return 
    }

    setButtonActive(btn);

    unsetAllTableEditOptions();

    loadInputData();
}

function btnFreezePressed(btn) {

    if(isButtonInactive(btn)) {
        return 
    }

    if(isButtonAvailable(btn)) {
        // the freeze button is pressed to freeze sequence
        sendFreezeFlag("freeze");
        setButtonActive(btn);
    } else {
        sendFreezeFlag("unfreeze");
        setButtonAvailable(btn);
    }
}

function loadInputData() {
    var d = getInputData();

  //  var url = $("#editReturnURL").text();

    //url += "/loadcfg/"+ JSON.stringify(d);

    var urlStr =  "/point/handlecfg/" + pointName + "/loadcfg";

    returnData(urlStr, d);
}

function saveInputData() {
    var d = getInputData();

  //  var url = $("#editReturnURL").text();

    //url += "/loadcfg/"+ JSON.stringify(d);

    var urlStr =  "/point/handlecfg/" + pointName + "/savecfg";

    returnData(urlStr, d);
}

function sendFreezeFlag(freeze) {
    var d = {};

  //  var url = $("#editReturnURL").text();

    //url += "/loadcfg/"+ JSON.stringify(d);

    var urlStr =  "/point/handlecfg/" + pointName + "/"+ freeze;

    returnData(urlStr, d);
}

function returnData(url, d) {

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
function getInputData() {

    var d = {};

    d["Start"] = getInputTableData($('#' + TABLE_START));
    d["Base"] = getInputTableData($('#' + TABLE_BASE));
    d["Finish"] = getInputTableData($('#' + TABLE_FINISH));

    return d;
}

function getInputTableData(tbl) {
    var d = [];

    tbl.find('tr:not(.'+TABLE_CLASS_ROW_NEW+')').each(function() {
        if(!$(this).hasClass(TR_CLASS_HEADER)) {
            d.push(getInputTrData($(this)));
        }    
    })        

    return d;
}

function getInputTrData(tr) {
    var d = {};

    d["Gpio"] = tr.find('.tdEditGpio').html(); 
    d["State"] = tr.find('.tdEditState').html(); 
    d["Interval"] = tr.find('.tdEditInterval').html(); 

    return d
}



function btnEditPressed(btn) {

    if(isButtonInactive(btn)) {
        return 
    }

    if(btn.hasClass('btn-warning')) {
        // available only, set active
        thisState |= STATE_EDIT;
        setButtonActive(btn);
        setAllTableEditOptions();
   }
    else if(btn.hasClass('btn-success')) {
        // active, set available only
        thisState &= ~STATE_EDIT;
        unsetAllTableEditOptions();
    }
}

function unsetAllTableEditOptions() {
    $('.tdEdit').attr('contenteditable', 'false');
    $('.tdEdit').attr('oninput', '');

    $('.trEdit').attr('draggable', 'false');
    $('.trEdit').removeClass(TR_CLASS_DRAGGED);

    $('.tdEditOnly').hide();
    $('.tdEditDelete').hide();
    $('.tdEditAdd').hide();
    
    thisState &= ~STATE_EDIT;
    unsetEditButtons();
}

function unsetEditButtons() {
    var count = 0;
    count += (null == cfgRun["Start"]) ? 0 : cfgRun["Start"].length;
    count += (null == cfgRun["Base"]) ? 0 : cfgRun["Base"].length;
    count += (null == cfgRun["Finish"]) ? 0 : cfgRun["Finish"].length;

    setButtonAvailable($('#'+BTN_EDIT));
    setButtonFreeze();
}

function setButtonFreeze() {
    var btn = $('#'+BTN_FREEZE);

    if(thisState & STATE_EDIT) {
        setButtonInactive(btn);
        return;
    }

    if(isFreeze) {
        setButtonActive(btn);
    } else if(isActive) {
        setButtonAvailable(btn);
    } else {
        setButtonInactive(btn);
    }    
}

function setEditButtons() {
    var count = 0;
    count += (null == cfgRun["Start"]) ? 0 : cfgRun["Start"].length;
    count += (null == cfgRun["Base"]) ? 0 : cfgRun["Base"].length;
    count += (null == cfgRun["Finish"]) ? 0 : cfgRun["Finish"].length;

    setButtonActive($('#'+BTN_EDIT));
    setButtonFreeze();
}

function setAllTableEditOptions() {
    $('.trEdit').removeClass(TR_CLASS_ACTIVE_ROW);

    $('.tdEdit').attr('contenteditable', 'true');
    $('.tdEdit').attr('oninput', 'checkInput($(this))');
    setTablesDraggable();

    $('.tdEditDelete').show();
    $('.tdEditTabHead').show();
    setTablesAddButton()

    thisState |= STATE_EDIT;
    setEditButtons();
}    

function setTablesAddButton() {
    setTableAddButtonOption(TABLE_START);
    setTableAddButtonOption(TABLE_BASE);
    setTableAddButtonOption(TABLE_FINISH);
}

function setTableAddButtonOption(tab) {
    var tr = $('#' + tab).find('.' + TABLE_CLASS_ROW_NEW);
    return setTableAddTr(tr);
}

function setTableAddTr(tr) {
    var td = tr.find('.tdEditAdd');
    if(checkInputRow(tr)) {
        tr.find('.tdEditAdd').show();
        return;
    } 
    tr.find('.tdEditAdd').hide();
}

function setTablesDraggable() {
    setTableDraggableOption(TABLE_START);
    setTableDraggableOption(TABLE_BASE);
    setTableDraggableOption(TABLE_FINISH);
}

function setTableDraggableOption(tab) {
    var tb = $('#' + tab);
    tb.find('.trEdit').attr('draggable', 'true');

    var removeDragg = $(tb.find('.' + TABLE_CLASS_ROW_NEW).last());
    removeDragg.attr('draggable', 'false');

    setTableSortedOption(tab);
}

function setTableSortedOption(tab) {
    var tbSort = $('#' + tab).find('tbody');

    tbSort.sortable({
        items: "tr[draggable='true']",
        update:function(event, ui){
            $(ui.item).css('background-color', COLOR_DRAG);
            $(ui.item).addClass(TR_CLASS_DRAGGED);
            inputReady2Use();
        }
    });

   // tbSort.on('dblclick', function(){toggleSort($(this));});       
   tbSort.on('dblclick', function(){toggleSortAll();});       
   tbSort.sortable('disable');
}

function toggleSortAll() {
//    var tbSort = $('#' + TABLE_START).find('tbody');
    toggleSort($('#' + TABLE_START).find('tbody'));       
    toggleSort($('#' + TABLE_BASE).find('tbody'));       
    toggleSort($('#' + TABLE_FINISH).find('tbody'));       



    //    tbSort.sortable('disable');
}


function toggleSort(obj) {
    var opts = obj.sortable('option');

    if(opts["disabled"]) {
        // disabled, set active
        obj.sortable( "enable" );
        $("tr[draggable='true']").find('.tdEdit').css("font-weight","bold");
    } else {
        // active, set disabled
        $("tr[draggable='true']").find('.tdEdit').css("font-weight",fontWeightOri);
        obj.sortable( "disable" );
    }
};

function setInterv() {
    if(0 > intervN) {
        intervN = setInterval("handleCfg()",1000);   // 1 sec
    }    
}

function unsetInterv() {
    if(0 <= intervN) {
        clearInterval(intervN);
        intervN = -1;
    }    
}

function drawCfg(data) {

    drawCfgTable(cfgThis["Start"], TABLE_START, "Start", data["Index"]["Start"]);

    drawCfgTable(cfgThis["Base"], TABLE_BASE, "Base", data["Index"]["Base"]);

    drawCfgTable(cfgThis["Finish"], TABLE_FINISH, "Finish", data["Index"]["Finish"]);
}

function drawCfgTable(data, table, title, ind) {

    var obj = $('#' + table);
    var rowCount = ((null == data) || (0 == data.length)) ? 0 : data.length;

    obj.empty()
    var str = "";

    str += partTitle(data, title);
    
    str += '<table id="editable-def" dropzone="move" class="pure-table pure-table-bordered">';

    str += tableTabHead();

    var i = 0;
    for (i = 0; i < rowCount; i++) {
        str += tableTabRow(data[i], i, ind, false);
    }

    str += tableTabRowNew();

    str += '</table>';
    str += '</br>';

    obj.html(str);

    createButtonAdd(obj.find('.tdEditAdd'));
    createButtonDelete(obj.find('.tdEditDelete'));

    obj.find('.tdEditOnly').hide();
}

function createButtonDelete(o) {
    o.button({
        label:J_BUTTON_LABEL_DELETE, 
        icons:{secondary:' ui-icon-closethick'}
    })

    o.button().on('click', function() {
        htmlRemoveTdRow($(this));
    })
}

function htmlRemoveTdRow(btn) {
    var row = btn.closest('tr');
    row.remove();
}

function jButtonClick(btn) {

    var label = btn.button('option', 'label');

    if (J_BUTTON_LABEL_DELETE == label) {
        htmlRemoveTdRow(btn);
    }

    if (J_BUTTON_LABEL_ADD == label) {
        htmlAddNewRow(btn);
    }
}

function htmlAddNewRow(btn) {
    // find the button table
    var tbl = btn.closest('table');

    // find the button row in the table
    var row = btn.closest('tr');

    // remove classes specific to the 'NEW' row
    btn.removeClass("tdEditAdd");
    btn.addClass("tdEditDelete");
    row.removeClass(TABLE_CLASS_ROW_NEW);
  
    // destroy the 'NEW' button of the current 'NEW' row
    btn.button('destroy');
    // substitute the current 'ADD' button to 'DELETE' button 
    // which is required for table data rows
    createButtonDelete(btn);

    // set the row draggable
    row.attr("draggable", "true");

    // prepare a new 'NEW' row html code to substitute the current 'NEW' row 
    // which is ready to add to the table data rows
    var str = tableTabRowNew()
    // add the new row html row code after the last row
    tbl.find('tr:last').after(str);

    // find the last row after adding html code
    row = tbl.find('tr:last');

    row.find('.tdEdit').attr('contenteditable', 'true');
    row.find('.tdEdit').attr('oninput', 'checkInput($(this))');
    createButtonAdd(row.find('.tdEditAdd'));

    setTableAddTr(row);
    
    inputReady2Use();
}

function createButtonAdd(o) {
    o.button({
        label:J_BUTTON_LABEL_ADD, 
        icons:{primary:'ui-icon-plusthick'}
    });

    o.button().on('click', function() {
        jButtonClick($(this));
    })
}

function checkInputRow(tr) {

    var str = tr.find('.tdEditGpio').html(); 
    if (!checkInputGpio(str)) {
        return false;
    }

    str = tr.find('.tdEditState').html(); 
    if (!checkInputState(str)) {
        return false;
    }

    str = tr.find('.tdEditInterval').html(); 
    if (!checkInputInterval(str)) {
        return false;
    }

    return true
}

function tableTabRow(data, i, ind, isNew) {;

    var str = "";

    var trClass = "trEdit" + (isNew ? (" " + TABLE_CLASS_ROW_NEW) : "");
    if(!isNew && (i == ind)) {
        trClass += ' ' + TR_CLASS_ACTIVE_ROW;
    }

    str += '<tr draggable="';
    str += isNew ? "false" : "true";
 
    str += '" class="' + trClass + '">';

    str += partTabCols(data);

    if(isNew) {
        // add button "add"
        str += '<td class="tdEditAdd tdEditOnly"></td>';
    } else {
        // add button "delete"
        str += '<td class="tdEditDelete tdEditOnly"></td>';
    }    
    
    return str;
}

function tableTabRowNew() {;

    var data = {Gpio:"new", State:"new", Interval:"new:new:new"};
    var str = tableTabRow(data, -1, -2, true);

    return str;
}

function partTabCols(data) {

    var str = "";

    str += '<td class="tdEdit tdEditGpio"     data-ori="' + data["Gpio"] +     '">' + data["Gpio"] + '</td>';
    str += '<td class="tdEdit tdEditState"    data-ori="' + data["State"] +    '">' + data["State"] + '</td>';
    str += '<td class="tdEdit tdEditInterval" data-ori="' + data["Interval"] + '">' + data["Interval"] + '</td>';

    return str;
}

function tableTabHead() {
    
    var str = "";

    str += '<thead>';
    str += '    <tr class="'+TR_CLASS_HEADER+'">';
    str += '        <th>GPIO</th>';
    str += '        <th>STATE</th>';
    str += '        <th>INTERVAL</th>';
    str += '        <th class="tdEditOnly tdEditTabHead"></th>';
    str += '     </tr>';
    str += '</thead>';

    return str;
}

function partTitle(data, title) {

    var str = "";

    str += '<h2>' + title + '</h2>';

    return str;
}

function setButtonsNonEdit() {
    $('#'+BTN_LOAD).text(BTN_LOAD_TXT);
    $('#'+BTN_LOAD_DEFAULT).text(BTN_LOAD_DEFAULT_TXT);
    $('#'+BTN_LOAD_SAVED).text(BTN_LOAD_SAVED_TXT);
    $('#'+BTN_SAVE).text(BTN_SAVE_TXT);
    $('#'+BTN_FREEZE).text(BTN_FREEZE_TXT );
    $('#'+BTN_EDIT).text(BTN_EDIT_TXT);

    // unset all buttons if the final part is used (it means the exit or restart has been pressed)
    if(isFinishActive()) {
        setButtonInactive($('#'+BTN_FREEZE));
        setButtonInactive($('#'+BTN_LOAD));
        setButtonInactive($('#'+BTN_SAVE));
        setButtonInactive($('#'+BTN_LOAD_DEFAULT));
        setButtonInactive($('#'+BTN_LOAD_SAVED));
        setButtonInactive($('#'+BTN_EDIT));

        return;
    }

    setButtonInactive($('#'+BTN_LOAD));
    setButtonInactive($('#'+BTN_LOAD_DEFAULT));
    setButtonInactive($('#'+BTN_SAVE));
    if(!isActive) {
        setButtonInactive($('#'+BTN_FREEZE));
    }
    setButtonAvailable($('#'+BTN_EDIT));
    setButtonAvailable($('#'+BTN_LOAD_DEFAULT));
    setButtonAvailable($('#'+BTN_LOAD_SAVED));
}

function setButtonInactive(btn) {
    btn.removeClass('btn-warning btn-success active');
    btn.addClass('btn-outline-secondary disabled');
}

function setButtonActive(btn) {
    btn.removeClass('btn-outline-secondary').removeClass('btn-warning').removeClass('disabled');
    btn.addClass('btn-success').addClass('active');
}

function setButtonAvailable(btn) {
    btn.removeClass('btn-outline-secondary').removeClass('btn-success').removeClass('disabled');
    btn.addClass('btn-warning').addClass('active');
}

function isButtonInactive(btn) {
    return btn.hasClass('btn-outline-secondary') && btn.hasClass('disabled');
}

function isButtonActive(btn) {
    return btn.hasClass('btn-success') && btn.hasClass('active');
}

function isButtonAvailable(btn) {
    return btn.hasClass('btn-warning') && btn.hasClass('active');
}

function areSetsEqual(d1, d2) {
    return JSON.stringify(d1) === JSON.stringify(d2);
}    

function inputHasError() {
    if(inputTableHasError($('#' + TABLE_START))) {
        return true;
    }
    if(inputTableHasError($('#' + TABLE_BASE))) {
        return true;
    }
    if(inputTableHasError($('#' + TABLE_FINISH))) {
        return true;
    }

    return false;
}

function inputHasChanges() {
    if(inputTableHasChanges($('#'+TABLE_START))) {
        return true;
    }
    if(inputTableHasChanges($('#'+TABLE_BASE))) {
        return true;
    }
    if(inputTableHasChanges($('#'+TABLE_FINISH))) {
        return true;
    }

    return false;
}


function inputTableHasError(tbl) {

//    tr:not(.gridTitleRow, .gridSpan)TABLE_CLASS_ROW_NEW

    if(tbl.find('tr:not(.'+TABLE_CLASS_ROW_NEW+') td').hasClass(TD_CLASS_EDIT_ERROR)) {


//    if(tbl.find('tr td').hasClass(TD_CLASS_EDIT_ERROR)) {
        return true;
    }

    return false;
}

function inputTableHasChanges(tbl) {
    if(tbl.find('tr:not(.'+TABLE_CLASS_ROW_NEW+') td').hasClass(TD_CLASS_EDIT_OK)) {
 //   if(tbl.find('tr td').hasClass(TD_CLASS_EDIT_OK)) {
        return true;
    }

    if(0 <tbl.find('tr.'+TR_CLASS_DRAGGED).length) {
        return true;
    }

//    tbl.find('tr').each(function () {

      //  var moha = $(this).css("background-color");

    //    if($(this).css("background-color") == COLOR_DRAG) {
  //          return true;
//    }});
    
    return false;
}


function inputReady2Use() {

    setButtonInactive($('#'+BTN_LOAD));
    setButtonInactive($('#'+BTN_SAVE));

    setButtonAvailable($('#'+BTN_LOAD_DEFAULT));
    setButtonAvailable($('#'+BTN_LOAD_SAVED));

    if(inputHasError()) {
        return;
    }

    if(!areSetsEqual(cfgRun, cfgSaved)) {
        setButtonAvailable($('#'+BTN_SAVE));  
    }

    if(!inputHasChanges()) {
        return;
    }

    setButtonAvailable($('#'+BTN_LOAD));
}

function isFinishActive() {
    if (0 <= indexAll['Finish']) {
        return true;
    }

    return false;
}


function handleCfg() {
 
    var urlStr =  "/point/" + pointName + "/getpointcfg";

    $.ajax({
        url: urlStr,
        type: 'post',
        data: {}, 
        dataType: 'json',
        timeout: 500,
        success : function(data) {
            cfgSaved = data["CfgSaved"];      // the point configuration saved on disk
            cfgRun = data["CfgRun"];          // the current point configuration of the app need to keep separately

            indexAll = data['Index'];

            cfgThis = data["CfgRun"];              // the current point configuration      
            isActive = data["Active"];
            isFreeze = 0 < (data["State"] & POINT_STATE_FREEZE); 

            if(isFinishActive()) {
                // now active is the finish part (exit or restart)
                thisState &= ~STATE_EDIT;
            }

            if(!(thisState & STATE_EDIT)) {
                drawCfg(data);

                if(start) {
                    colorOri = $('.tdEdit').css('background-color');
                    fontWeightOri = $('.tdEdit').css('font-weight');
                }
                unsetAllTableEditOptions();

                setButtonsNonEdit();
                if(!isFinishActive() && !areSetsEqual(cfgRun, cfgSaved)) {
                    setButtonAvailable($('#'+BTN_SAVE));  
                }
            }  

            if(!isFinishActive()) {
                setButtonAvailable($('#'+BTN_LOAD_DEFAULT)); 
                setButtonAvailable($('#'+BTN_LOAD_SAVED));  
            }

            start = false;
        },
    });
}
*/