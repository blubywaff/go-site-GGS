var result;
function getIns() {
    var fail = document.getElementById("fail").value;
    var fails = fail.split(",");
    for(var i = 0; i < fails.length; i++) {
        fails[i].replace(" ", "");
        fails[i] = parseFloat(fails[i]);
    }
    var susp = document.getElementById("susp").value;
    var susps = susp.split(",");
    for(var i = 0; i < susps.length; i++) {
        susps[i].replace(" ", "");
        susps[i] = parseFloat(fails[i]);
    }
    result = sepIn(fails, susps, 0);
}
function setOut() {
    var eta = document.getElementById("eta");
    var beta = document.getElementById("beta");
    var r2 = document.getElementById("r2");
    eta.innerHTML = "eta: " + ETA;
    beta.innerHTML = "beta: " + BETA;
    r2.innerHTML = "r-squared: " + R2;
    if(R2 > 1) {
        r2.innerHTML = r2.innerHTML + " (sample size too small)";
    }
}

function generateLine() {
    dpsLine = [];
    for(var i = 1; i <= 10000; i++) {
        dpsLine.push({x: i, y: (1-calcY(i, ETA, BETA))*100});
    }
}

function generateData() {
    dps = [];
    for(var i = 0; i < result[0].length; i++) {
        dps.push({x: result[0][i], y: result[1][i]*100});
    }

}

function changeLog() {
    isLog = !isLog;
    if(!isLog) {
        document.getElementById("logbut").innerHTML = "Lin";
        unit = 1000;
    } else {
        document.getElementById("logbut").innerHTML = "Log";
        unit = 0.5;
    }
    makeGraph();
}

