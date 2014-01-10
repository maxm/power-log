(function() {

  var pollPulsesHours = 6
  var pollPulseMillis = pollPulsesHours*60*60*1000;
  var lastPulses = []
  var chartsReady = false

  google.load('visualization', '1.0', {'packages':['corechart']});
  google.setOnLoadCallback(function() { chartsReady = true; });

  $(document).ready(function() {
    pollPulses();
    updateNow();
  });

  function updateNow() {
    if (lastPulses.length > 2) {
      var lastPulse = lastPulses[lastPulses.length-1];
      var watts = deltaToWatts(lastPulse - lastPulses[lastPulses.length-2])
      $('#wattsNow').text(watts.toString())
      var now = new Date().getTime();
      var seconds = (now - lastPulse) / 1000;
      if (seconds < 45) {
        $('#wattsNowTime').text("seconds ago")
      } else {
        var minutes = Math.round(seconds / 60);
        if (minutes == 1) {
          $('#wattsNowTime').text("one minute ago")  
        } else {
          $('#wattsNowTime').text(minutes + " minutes ago")
        }
      }
    }
    setTimeout(updateNow, 5000);
  }

  function pollPulses() {
    var time = new Date().getTime();
    $.get("range", { from: time - pollPulseMillis },function(data) {
      lastPulses = data.Pulses;
      drawChart();
      updateNow();
      setTimeout(pollPulses, 2000);
    }, "json");
  }

  function drawChart() {
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
      },
      vAxis: { 
        viewWindow:{
          min: 0
        }
      },
      chartArea: {
        height: 600
      }
    };

    var chart = new google.visualization.LineChart(document.getElementById('chart'));
    chart.draw(data, options);
  }

  function deltaToWatts(delta) {
    return Math.round(60*60*1000/delta);
  }

})();