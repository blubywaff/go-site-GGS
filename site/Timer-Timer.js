var startTime;
var finish;
var duration = 0;
var timerID;
var Oduration = 0;
var canSet = true;
var isGo = false;
//var STARTTIMELOG;
//var FINISHTIMELOG;
function setLabel() {
    var label = document.getElementById("timeLabel");
    if(duration <= 0 && !canSet) {
        label.innerHTML = "DONE";
        label.style.background = "#ff4444";
        return;
    } else if(duration <= 0 && canSet) {
        label.innerHTML = "00:00:00";
        label.style.background = "#00ccff";
        return;
    }
    var time = "";
    var hTime = 0;
    var mTime = 0;
    var sTime = 0;
    var remTime = duration;
    hTime = Math.floor(remTime/3600000);
    remTime = remTime%3600000;
    mTime = Math.floor(remTime/60000);
    remTime = remTime%60000;
    sTime = Math.floor(remTime/1000);
    time = hTime + ":" + mTime + ":" + sTime;
    label.innerHTML = time;
    if(canSet) {
        label.style.background = "#00ccff";
        return;
    }
    if(!isGo) {
        label.style.background = "#ffff00";
        return;
    }
    label.style.background = "#00cc00";
}
function setCanSet(csSet, igSet) {
    canSet = csSet;
    isGo = igSet;
        for(var i = 0; i < document.getElementsByClassName("timeChanger").length; i++) {
            document.getElementsByClassName("timeChanger")[i].disabled = !canSet;
        }
        document.getElementById("stopButton").disabled = !isGo;
        document.getElementById("resetButton").disabled = canSet;
        document.getElementById("startButton").disabled = isGo;
        setLabel();
}
function DoTimer() {
    duration = finish - new Date().getTime();
    if(duration <= 0) {
        stop();
        document.getElementById("aaAAhhdio").play();
    }
    setLabel();
}
function start() {
    if(isGo) {
        return;
    }
    startTime = new Date().getTime();
    finish = startTime + duration;
    timerID = setInterval(DoTimer, 100);
    setCanSet(false, true);
}
function stop() {
    if(canSet) {
        return;
    }
    duration = finish - new Date().getTime();
    clearInterval(timerID);
    document.getElementById("stopButton").disabled = true;
    setCanSet(false, false);
}
function reset() {
    if(canSet) {
       return;
    }
    stop();
    duration = Oduration;
    setCanSet(true, false);
}
function parseIn(inputString) {
    var hTime = 0;
    var mTime = 0;
    var sTime = 0;
    var inputStrings = inputString.split(" ");
    for (var i = 0; i < inputStrings.length; i++) {
        var temptime;
        if (inputStrings[i].endsWith('h')) {
            temptime = inputStrings[i].substring(0, inputStrings[i].indexOf("h"));
            hTime += parseInt(temptime);
        } else if(inputStrings[i].endsWith("m")) {
            temptime = inputStrings[i].substring(0, inputStrings[i].indexOf("m"));
            mTime += parseInt(temptime);
        } else if(inputStrings[i].endsWith('s')) {
            temptime = inputStrings[i].substring(0, inputStrings[i].indexOf("s"));
            sTime += parseInt(temptime);
        }
    }
    return 3600000 * hTime + 60000 * mTime + 1000 * sTime;
}
function set() {
    if(!canSet) {
        return;
    }
    Oduration = parseIn(document.getElementById("setField").value);
    if(Oduration < 0) {
        Oduration = 0;
    }
    duration = Oduration;
    setLabel();
}
function add() {
    if(!canSet) {
        return;
    }
    Oduration += parseIn(document.getElementById("addField").value);
    if(Oduration < 0) {
        Oduration = 0;
    }
    duration = Oduration;
    setLabel();
}
function remove() {
    if(!canSet) {
        return;
    }
    Oduration -= parseIn(document.getElementById("removeField").value);
    if(Oduration < 0) {
        Oduration = 0;
    }
    duration = Oduration;
    setLabel();
}