$(document).ready(function () {

    

    var tweets = $(".tweet");

    $(tweets).each(function (t, tweet) {

        var id = $(this).attr('id');

        twttr.widgets.createTweet(
            id, tweet,
            {
                conversation: 'none',    // or all
                cards: 'visible',  // or visible 
                //linkColor: '#cc0000', // default is blue
                theme: 'light'    // or dark
            });
    });
});