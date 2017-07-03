$(document).ready(function () {
    //$('.croppieimg').croppie();
    var basic = $('#croppieimg').croppie({
        viewport: {
            width: 150,
            height: 200
        }
    });
    basic.croppie('bind', {
        url: '/img/stabi/Kupferstiche/PPN78358086X',
        points: [77, 469, 280, 739]
    });
    //on button click
    basic.croppie('result', 'html').then(function (html) {
        // html is div (overflow hidden)
        // with img positioned inside.
    });

});