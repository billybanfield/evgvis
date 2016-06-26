run = function() {
    var fetchDataPromise = $.ajax({
        dataType: "json",
        url: "/data",
    });

    fetchDataPromise.then(function(data) {
        alert(data);
    }).done(function() {
        alert("done");
    });
    console.log("hello");
};
run();
