(function() {

  var pollPulseDelta = 2*60*60*1000; // 2 hours
  var lastPulses = []

  $(document).ready(function() {
    pollPulses();
  });

  function pollPulses() {
    var time = new Date().getTime();
    $.get("range", { from: time - pollPulseDelta },function(data) {
      lastPulses = data.Pulses;
      if (lastPulses.length > 2) {
        var delta1 = new Date().getTime() - lastPulses[lastPulses.length-1]
        var delta2 = lastPulses[lastPulses.length-1] - lastPulses[lastPulses.length-2];
        var delta = delta1 > 60*60*10 ? delta1 : delta;
        var watts = Math.round(60*60*1000/delta);
        $('#currentWatts').text(watts.toString())
      }
      setTimeout(pollPulses, 2000);
    }, "json");
  }

})();