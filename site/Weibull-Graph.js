var dpsLine = [];
var dps = [];
var isLog = true;
var unit = 0.5;
var lineColor = "#ff4b4b";
var pointColor = "#4b4bff";

function makeGraph() {
    if(!window.navigator.onLine) {
        alert("Internet connection failed. Restore connection and reload page to restore graph.");
    }
    var chart = new CanvasJS.Chart("chartContainer", {
        animationEnabled: true,
        zoomEnabled: true,
        theme: "light1",
        title:{
            text: "Weibull Plot"
        },
        axisX:{
            logarithmic: isLog,
            title: "Time to Failure",
            valueFormatString: "####",
            interval: unit
        },
        axisY:{
            logarithmic: isLog,
            title: "Probability of Failure %",
            titleFontColor: lineColor,
            lineColor: lineColor,
            gridThickness: 0,
            lineThickness: 1,
            includeZero: false,
            labelFormatter: addSymbols
        },
        legend:{
            verticalAlign: "top",
            fontSize: 16,
            dockInsidePlotArea: true
        },
        data: [{
            type: "spline",
            xValueFormatString: "####",
            showInLegend: true,
            name: "Weibull Line",
            color: lineColor,
            dataPoints: dpsLine,
            toolTipContent: "<font color=\"" + lineColor + "\">{x}: </font>{y}%"
        },
            {
                type: "scatter",
                xValueFormatString: "####",
                showInLegend: true,
                name: "Data Points",
                color: pointColor,
                toolTipContent: "<font color=\"" + pointColor + "\">{x}: </font>{y}%",
                dataPoints: dps
            }
        ]
    });
    chart.render();
    function addSymbols(e){
        var suffixes = ["", "K", "M", "B"];
        var order = Math.max(Math.floor(Math.log(e.value) / Math.log(1000)), 0);
        if(order > suffixes.length - 1)
            order = suffixes.length - 1;
        var suffix = suffixes[order];
        return CanvasJS.formatNumber(e.value / Math.pow(1000, order)) + suffix;
    }
}