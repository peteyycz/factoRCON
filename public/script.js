$(function () {
  $(".panel").on('keypress', ".in", function(e) {
    if (e.which == 13) {
      $(this).prop('readonly', true);
      var command = $(this).val();

      axios({
        method: 'post',
        url: '/command',
        data: {
          command
        }
      }).then(function (output) {
        $(".output").last().html(output.data.result)
        $(".panel").append(
          $("<div class='action'>")
          .html("<div class='action'><div class='command'><span class='symbol'>$</span><input class='in' type='text'></div><div class='output'></div></div>"));
          $(".in").last().focus();
      })
    }
  });
  $('.panel').stop().animate({
    scrollTop: $(".panel")[0].scrollHeight
  }, 800);
})
