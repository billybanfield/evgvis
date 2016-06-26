        var fetchStatePromise = $.ajax({
                type: "GET",
                url: "/data",
                dataType: "json"
            });

    fetchStatePromise.done(runD3)
      .fail(function() {
        alert("failed to load json")
      });

