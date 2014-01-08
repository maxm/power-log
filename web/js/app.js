(function() {

  var pollPulsesHours = 6
  var pollPulseMillis = pollPulsesHours*60*60*1000;
  var lastPulses = []
  var chartsReady = false

  google.load('visualization', '1.0', {'packages':['corechart']});
  google.setOnLoadCallback(function() { chartsReady = true; });

  $(document).ready(function() {
    pollPulses();
  });

  function pollPulses() {
    var time = new Date().getTime();
    $.get("range", { from: time - pollPulseMillis },function(data) {
      lastPulses = data.Pulses;
      if (lastPulses.length > 2) {
        var watts = deltaToWatts(lastPulses[lastPulses.length-1] - lastPulses[lastPulses.length-2])
        $('#wattsNow').text(watts.toString())
        drawNowChart();
      }
      setTimeout(pollPulses, 2000);
    }, "json");
  }

  function drawNowChart() {
    if(!chartsReady) return;
    var watts = [['Time', 'Power']];
    var lastT = 0;
    $.each(lastPulses, function(i,t) {
      if (i > 0) {
        watts.push([new Date(t), deltaToWatts(t-lastT)])
      }
      lastT = t;
    })
    var data = google.visualization.arrayToDataTable(watts);

    var options = {
      title: 'Last ' + pollPulsesHours + ' hours',
      hAxis: {
        format: 'HH:mm'
      }
    };

    var chart = new google.visualization.LineChart(document.getElementById('nowChart'));
    chart.draw(data, options);
  }

  function deltaToWatts(delta) {
    return Math.round(60*60*1000/delta);
  }

})();