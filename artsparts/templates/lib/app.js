$(document).ready(function () {
    $('.ui.dropdown').dropdown();
    $('.special.cards .image').dimmer({
        on: 'hover'
    });
    $('.openmodal').click(function(){
        $('.ui.modal').modal('show');
    });
    
    $('#image').cropper({

        viewMode: 1,
        preview: '#preview',
        guides: false,
        center: false,
        background: false,
        crop: function (e) {
            // Output the result data for cropping image.
            console.log(e.x);
            console.log(e.y);
            console.log(e.width);
            console.log(e.height);
            console.log(e.rotate);
            console.log(e.scaleX);
            console.log(e.scaleY);
        }
    });
});