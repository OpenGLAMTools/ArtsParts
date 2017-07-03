$(document).ready(function () {
    $('.ui.dropdown').dropdown();
    $('.special.cards .image').dimmer({
        on: 'hover'
    });
    $('.openmodal').click(function(){
        $('.ui.modal').modal('show');
    });
 
});